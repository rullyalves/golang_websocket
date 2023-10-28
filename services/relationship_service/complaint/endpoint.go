package complaint

import "net/http"

func SendComplaintHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		//TODO:FAZER

		writer.WriteHeader(http.StatusCreated)
	}
}
