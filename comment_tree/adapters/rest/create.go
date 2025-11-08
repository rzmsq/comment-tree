package rest

import (
	"CommentTree/comment_tree/core"
	"CommentTree/comment_tree/usecase"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/goccy/go-json"
)

func NewCreateHandler(u usecase.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := r.Body.Close()
			if err != nil {
				u.Log.Error("Error closing body", "err", err)
				return
			}
		}()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			u.Log.Error("Error reading body", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		requestData := core.Comment{}
		err = json.Unmarshal(body, &requestData)
		if err != nil {
			u.Log.Error("Error unmarshalling body", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var id int
		id, err = u.CreateComment(r.Context(), requestData)
		if errors.Is(err, core.ErrNotValidRequest) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err != nil {
			u.Log.Error("Error creating comment", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		write, err := w.Write([]byte(fmt.Sprintf("%d", id)))
		if err != nil {
			u.Log.Error("Error writing response", "err", err, "write", write)
			return
		}
	}
}
