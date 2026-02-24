// Package state is the implementation package for the message task state.
package state

import (
	"context"
	"fmt"
	"sync"

	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/magicbox/enum"
)

type MessageTask struct {
	NamespaceUID snowflake.ID
	MessageUID   snowflake.ID
	retryCount   int
	mu           sync.Mutex
}

func (m *MessageTask) RetryIncrement() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.retryCount++
}

func (m *MessageTask) IsMaxRetry() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.retryCount >= 2
}

type ProcessFunc func(task *MessageTask) (nextStatus enum.MessageStatus, isNext bool)

/*
MessageTaskState state transition rule matrix (reference design document)

	+-------------------+----------+------------+-------------+----------+---------+
	| current status    | Start    | SendSuccess| SendFailure | Cancel   | Retry   |
	+-------------------+----------+------------+-------------+----------+---------+
	| PENDING           | SENDING  | -          | -           | CANCELLED| -       |
	| SENDING           | -        | SENT       | FAILED      | CANCELLED| -       |
	| FAILED            | -        | -          | -           | CANCELLED| PENDING |
	| SENT              | terminated state (reject all events)                     |
	| CANCELLED         | terminated state (reject all events)                     |
	| UNKNOWN           | only allowed Start → PENDING (initialization)            |
	+-------------------+----------+------------+-------------+----------+---------+
*/
type MessageTaskState struct {
	nextState   map[enum.MessageStatus]MessageTaskState
	processFunc ProcessFunc
}

func NewMessageTaskState(status enum.MessageStatus) MessageTaskState {
	processFunc, ok := GetMessageTaskProcess(status)
	if !ok {
		panic(fmt.Sprintf("message status %s state process func not found", status))
	}
	return MessageTaskState{
		processFunc: processFunc,
	}
}

func (m *MessageTaskState) Process(ctx context.Context, task *MessageTask) {
	nextStatus, isNext := m.processFunc(task)
	if !isNext {
		return
	}
	if nextState, ok := GetMessageTaskState(nextStatus); ok {
		nextState.Process(ctx, task)
	}
}
