package http_router

import (
	"encoding/json"
	"errors"
	"go_ws/services/shared/validations"
	"net/http"
)

type MutationAction[T any, S any] func(body S) (*T, error)

type VoidMutationAction[T any] func() (*T, error)

type ReadAction[T any] func() (*T, error)

func handleError[T any](writer http.ResponseWriter, err error) error {

	var invalidPathErr *InvalidPathParameterErr

	if errors.As(err, &invalidPathErr) {
		HandleError(err, writer, http.StatusBadRequest)
		return err
	}

	if err != nil {
		HandleError(err, writer, http.StatusInternalServerError)
		return err
	}

	return nil
}

func validateValues(writer http.ResponseWriter, request *http.Request, value any) (proceed bool) {
	invalidList := validations.ValidateErrors(request.Context(), "pt_BR", value)

	validationErr := struct {
		Errors []validations.ErrorItem `json:"errors"`
	}{Errors: invalidList}

	if len(invalidList) != 0 {
		data, _ := json.Marshal(validationErr)
		http.Error(writer, string(data), http.StatusBadRequest)
		return false
	}

	return true
}

func ExecuteReadRequest[T any](writer http.ResponseWriter, action ReadAction[T]) {
	result, err := action()
	err = handleError[T](writer, err)

	if err != nil {
		return
	}

	err = json.NewEncoder(writer).Encode(&result)

	if err != nil {
		HandleError(err, writer, http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func ExecuteVoidMutationRequest[T any](writer http.ResponseWriter, action VoidMutationAction[T]) {
	_, err := action()
	_ = handleError[T](writer, err)
}

func ExecuteMutationRequest[T any, S any](writer http.ResponseWriter, request *http.Request, action MutationAction[T, S]) {

	var body S

	err := json.NewDecoder(request.Body).Decode(&body)

	if err != nil {
		HandleError(err, writer, http.StatusBadRequest)
		return
	}

	if proceed := validateValues(writer, request, body); !proceed {
		return
	}

	result, err := action(body)

	err = handleError[T](writer, err)

	if result == nil {
		return
	}

	err = json.NewEncoder(writer).Encode(&result)

	if err != nil {
		HandleError(err, writer, http.StatusInternalServerError)
		return
	}
}
