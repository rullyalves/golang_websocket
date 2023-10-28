package complaint_category

import "net/http"

func FindComplaintCategoriesHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		//TODO:FAZER

		writer.WriteHeader(http.StatusCreated)
	}
}

func SaveComplaintCategoriesHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		//TODO:FAZER

		writer.WriteHeader(http.StatusCreated)
	}
}

func DeleteComplaintCategoriesHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		//TODO:FAZER

		writer.WriteHeader(http.StatusCreated)
	}
}
