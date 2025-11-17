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
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			namespace = tr.RequestHeader().Get(cnst.HTTPHeaderXNamespace)
			ctx = WithNamespace(ctx, namespace)
			tr.RequestHeader().Set(cnst.MetadataGlobalKeyNamespace, namespace)

			if strutil.IsNotEmpty(namespace) {
				return handler(ctx, req)
			}

			if md, ok := metadata.FromServerContext(ctx); ok {
				namespace = md.Get(cnst.MetadataGlobalKeyNamespace)
				ctx = WithNamespace(ctx, namespace)
				tr.RequestHeader().Set(cnst.MetadataGlobalKeyNamespace, namespace)
			}

			if strutil.IsNotEmpty(namespace) {
				return handler(ctx, req)
			}

			return nil, merr.ErrorForbidden("namespace is required, please set the namespace in the request header or metadata, Example: %s: default", cnst.HTTPHeaderXNamespace)
		}
	}
}

// MustNamespaceExist 检查namespace必须存在且有效
func MustNamespaceExist(hasNamespace func(ctx context.Context) error) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			if err := hasNamespace(ctx); err != nil {
				return nil, err
			}
			return handler(ctx, req)
		}
	}
}
