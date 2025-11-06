package core

import "context"

type DB interface {
	Add(context.Context, Comment) error
	Get(context.Context) ([]Comment, error)
	Drop(context.Context) error
}
