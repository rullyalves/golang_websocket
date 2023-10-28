package utils

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"time"
)

func ToLocalDateTime(date *time.Time) *neo4j.LocalDateTime {
	if date == nil {
		return nil
	}

	result := neo4j.LocalDateTimeOf(*date)

	return &result
}
