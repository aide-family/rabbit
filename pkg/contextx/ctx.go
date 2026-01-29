// Package contextx provides contextx for the application.
package contextx

import (
	"context"

	"github.com/bwmarrin/snowflake"
)

type (
	namespaceNameKey struct{}
	namespaceUIDKey  struct{}
	userUIDKey       struct{}
	usernameKey      struct{}
)

func WithNamespace(ctx context.Context, namespace string) context.Context {
	return context.WithValue(ctx, namespaceNameKey{}, namespace)
}

func GetNamespace(ctx context.Context) string {
	return ctx.Value(namespaceNameKey{}).(string)
}

func WithNamespaceUID(ctx context.Context, namespace snowflake.ID) context.Context {
	return context.WithValue(ctx, namespaceUIDKey{}, namespace)
}

func GetNamespaceUID(ctx context.Context) snowflake.ID {
	return ctx.Value(namespaceUIDKey{}).(snowflake.ID)
}

func WithUserUID(ctx context.Context, userUID snowflake.ID) context.Context {
	return context.WithValue(ctx, userUIDKey{}, userUID)
}

func GetUserUID(ctx context.Context) snowflake.ID {
	return ctx.Value(userUIDKey{}).(snowflake.ID)
}

func WithUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, usernameKey{}, username)
}

func GetUsername(ctx context.Context) string {
	return ctx.Value(usernameKey{}).(string)
}
