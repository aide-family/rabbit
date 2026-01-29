// Package authv1 is the repository for the auth service.
package authv1

import (
	"context"
)

type Repository interface {
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
}
