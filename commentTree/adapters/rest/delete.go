package rest

import (
	"CommentTree/commentTree/pkg/logger"
	"net/http"
)

func NewDeleteHandler(log logger.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
