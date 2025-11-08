package core

import "context"

type DB interface {
	Add(context.Context, Comment) (int, error)
	Get(context.Context, int) ([]Comment, error)
	GetAll(context.Context) ([]Comment, error)
	Delete(context.Context, int) error
}
