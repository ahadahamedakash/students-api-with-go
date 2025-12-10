package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/ahadahamedakash/students-api-with-go/internal/types"
	"github.com/ahadahamedakash/students-api-with-go/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("creating a student")

		var student types.Student

		err := json.NewDecoder((r.Body)).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))

			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		}

		// request validation

		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors) // type cast

			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))

			return
		}

		// w.Write([]byte("Welcome to students api with go"))

		response.WriteJson(w, http.StatusCreated, map[string]string{"success": "ok"})
	}
}
