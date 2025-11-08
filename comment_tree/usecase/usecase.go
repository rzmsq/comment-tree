package usecase

import (
	"CommentTree/comment_tree/core"
	"CommentTree/comment_tree/pkg/logger"
	"context"
	"strconv"

	"github.com/go-playground/validator/v10"
)

const defaultUsername = "Anonymous"

type UseCase struct {
	db       core.DB
	Log      logger.Interface
	Validate *validator.Validate
}

func NewUseCase(db core.DB, log logger.Interface, validate *validator.Validate) UseCase {
	return UseCase{
		db:       db,
		Log:      log,
		Validate: validate,
	}
}

func (u *UseCase) CreateComment(ctx context.Context, comment core.Comment) (int, error) {
	err := u.Validate.Struct(comment)
	if err != nil {
		return 0, core.ErrNotValidRequest
	}

	if comment.Username == "" {
		comment.Username = defaultUsername
	}

	return u.db.Add(ctx, comment)
}

func (u *UseCase) GetComments(ctx context.Context, commentIdStr string) (res []core.Comment, err error) {
	if commentIdStr != "" {
		var id int
		id, err = strconv.Atoi(commentIdStr)
		if id <= 0 {
			return nil, core.ErrNotValidRequest
		}
		if err != nil {
			return nil, err
		}
		res, err = u.db.Get(ctx, id)
		if err != nil {
			return nil, err
		}

		if len(res) == 0 {
			return nil, core.ErrNotFound
		}
	} else {
		res, err = u.db.GetAll(ctx)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (u *UseCase) Delete(ctx context.Context, commentIdStr string) error {
	id, err := validateAndReturnID(commentIdStr)
	if err != nil {
		return err
	}

	return u.db.Delete(ctx, id)
}

func validateAndReturnID(commentIdStr string) (int, error) {
	if commentIdStr == "" {
		return 0, core.ErrNotValidRequest
	}

	commentId, err := strconv.Atoi(commentIdStr)
	if commentId <= 0 {
		return 0, core.ErrNotValidRequest
	}
	if err != nil {
		return 0, err
	}

	return commentId, nil
}
