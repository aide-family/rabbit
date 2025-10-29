package middler

import (
	"context"

	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/magicbox/strutil/cnst"
	"github.com/aide-family/rabbit/pkg/merr"
	"github.com/go-kratos/kratos/v2/metadata"
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

func MustNamespace() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			var namespace string
			if tr, ok := transport.FromServerContext(ctx); ok {
				namespace = tr.RequestHeader().Get(cnst.HTTPHeaderXNamespace)
				ctx = WithNamespace(ctx, namespace)
			}
			if strutil.IsNotEmpty(namespace) {
				return handler(ctx, req)
			}

			if metadata, ok := metadata.FromServerContext(ctx); ok {
				namespace = metadata.Get(cnst.MetadataGlobalKeyNamespace)
				ctx = WithNamespace(ctx, namespace)
			}
			if strutil.IsNotEmpty(namespace) {
				return handler(ctx, req)
			}

			return nil, merr.ErrorForbidden("namespace is required, please set the namespace in the request header or metadata, Example: %s: default", cnst.HTTPHeaderXNamespace)
		}
	}
}
