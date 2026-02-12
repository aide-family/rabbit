// Package state is the implementation package for the message task state.
package state

import (
	"context"
	"fmt"

	"github.com/aide-family/magicbox/enum"
	"github.com/bwmarrin/snowflake"
)

// 状态转换规则矩阵（设计文档参考）
// +------------+----------+------------+-------------+----------+---------+
// | 当前状态    | Start    | SendSuccess| SendFailure | Cancel   | Retry   |
// +------------+----------+------------+-------------+----------+---------+
// | PENDING    | SENDING  | -          | -           | CANCELLED| -       |
// | SENDING    | -        | SENT       | FAILED      | CANCELLED| -       |
// | FAILED     | -        | -          | -           | CANCELLED| PENDING |
// | SENT       | 终止状态（拒绝所有事件）                            		   |
// | CANCELLED  | 终止状态（拒绝所有事件）                            		   |
// | UNKNOWN    | 仅允许 Start → PENDING（初始化）                    		  |
// +------------+----------+------------+-------------+----------+---------+

type MessageTask struct {
	NamespaceUID snowflake.ID
	MessageUID   snowflake.ID
	retryCount   int
}

func (m *MessageTask) RetryIncrement() {
	m.retryCount++
}

func (m *MessageTask) IsMaxRetry() bool {
	return m.retryCount >= 2
}

type ProcessFunc func(task *MessageTask) (enum.MessageStatus, bool)

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
	status, isNext := m.processFunc(task)
	if !isNext {
		return
	}
	if nextState, ok := GetMessageTaskState(status); ok {
		nextState.Process(ctx, task)
	}
}
