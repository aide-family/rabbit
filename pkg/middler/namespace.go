package middler

import (
	"context"

	"github.com/aide-family/magicbox/strutil/cnst"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

type namespaceKey struct{}

func WithNamespace(ctx context.Context, namespace string) context.Context {
	return context.WithValue(ctx, namespaceKey{}, namespace)
}

func GetNamespace(ctx context.Context) string {
	return ctx.Value(namespaceKey{}).(string)
}

func BindNamespace() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			var namespace string
			if tr, ok := transport.FromServerContext(ctx); ok {
				namespace = tr.RequestHeader().Get(cnst.HTTPHeaderXNamespace)
				ctx = WithNamespace(ctx, namespace)
			}

			return handler(ctx, req)
		}
	}
}
