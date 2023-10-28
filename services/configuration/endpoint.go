package configuration

import (
	"net/http"
)

func FindConfigurationHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		//TODO:FAZER

		writer.WriteHeader(http.StatusCreated)
	}
}

func UpdateConfigurationHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		//TODO:FAZER

		writer.WriteHeader(http.StatusCreated)
	}
}
