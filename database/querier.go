package database

import (
	"context"
)

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) error
}
