package rest

import (
	"CommentTree/comment_tree/core"
	"CommentTree/comment_tree/usecase"
	"encoding/json"
	"errors"
	"net/http"
)

func NewGetHandler(u usecase.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("parent")
		u.Log.Debug("Get Request", "idStr", idStr)

		comments, err := u.GetComments(r.Context(), idStr)
		if errors.Is(err, core.ErrNotValidRequest) || errors.Is(err, core.ErrNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err != nil {
			u.Log.Debug("Error getting comments", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		u.Log.Debug("Comments", "comments", comments)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(comments); err != nil {
			u.Log.Error("Error encoding response", "err", err)
			return
		}
	}
}
