package middler

import (
	"context"

	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/magicbox/strutil/cnst"
	"github.com/bwmarrin/snowflake"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/aide-family/rabbit/pkg/contextx"
	"github.com/aide-family/rabbit/pkg/merr"
)

func MustNamespace() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			var namespace string
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			namespace = tr.RequestHeader().Get(cnst.HTTPHeaderXNamespace)
			ctx = contextx.WithNamespace(ctx, namespace)
			tr.RequestHeader().Set(cnst.MetadataGlobalKeyNamespace, namespace)

			if strutil.IsNotEmpty(namespace) {
				return handler(ctx, req)
			}

			if md, ok := metadata.FromServerContext(ctx); ok {
				namespace = md.Get(cnst.MetadataGlobalKeyNamespace)
				ctx = contextx.WithNamespace(ctx, namespace)
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
func MustNamespaceExist(hasNamespace func(ctx context.Context) (snowflake.ID, error)) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			namespace, err := hasNamespace(ctx)
			if err != nil {
				return nil, err
			}
			ctx = contextx.WithNamespaceUID(ctx, namespace)
			return handler(ctx, req)
		}
	}
}
