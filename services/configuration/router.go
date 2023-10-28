package configuration

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go_ws/shared/http_router"
)

func Handle(router http_router.Router, neo4jDriver *neo4j.DriverWithContext) {

	// configuration
	router.Get("/configurations", FindConfigurationHandler())
	router.Post("/configurations", UpdateConfigurationHandler())
}
