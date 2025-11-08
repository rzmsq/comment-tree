package rest

import (
	"CommentTree/comment_tree/usecase"
	"net/http"
)

func NewDeleteHandler(u usecase.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		err := u.Delete(r.Context(), id)
		if err != nil {
			u.Log.Debug("Error deleting comment", "id", id, "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
