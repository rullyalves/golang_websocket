package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go_ws/services/profile_service/filter/models"
	optionModels "go_ws/services/profile_service/options/models"
	neo4jdb "go_ws/shared/database/neo4j"
	"time"
)

type CreateFilterDataParams struct {
	ID         string                  `json:"id"`
	OptionIds  []string                `json:"data"`
	IsRequired bool                    `json:"isRequired"`
	Type       optionModels.OptionType `json:"type"`
	CreatedAt  time.Time               `json:"createdAt"`
}

type FindFilterByPreferencesId func(context.Context, string) (*[]models.FilterView, error)

type SaveFilter func(context.Context, string, CreateFilterDataParams) error

type SaveAllFilters func(context.Context, string, []CreateFilterDataParams) error

type DeleteFilterById func(ctx context.Context, id string) error

func FindByPreferencesId(driver *neo4j.DriverWithContext) FindFilterByPreferencesId {
	return func(ctx context.Context, id string) (*[]models.FilterView, error) {
		params := map[string]any{
			"preferencesId": id,
		}

		query := `
		MATCH (pref: Preferences{id: $preferencesId})
		MATCH (pref)-[:HAS]->(f :Filter) 
		RETURN f`

		results, err := neo4jdb.ExecuteWithMapping[models.FilterView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return &results, nil
	}
}

func SaveAll(driver *neo4j.DriverWithContext) SaveAllFilters {
	return func(ctx context.Context, preferencesId string, data []CreateFilterDataParams) error {

		var items []map[string]any

		for _, item := range data {

			items = append(items, map[string]any{
				"id":         item.ID,
				"optionIds":  item.OptionIds,
				"type":       item.Type,
				"isRequired": item.IsRequired,
				"createdAt":  item.CreatedAt,
			})
		}

		params := map[string]any{
			"preferencesId": preferencesId,
			"items":         items,
		}

		//TODO: melhorar update/delete dos filtros, talvez deletar os que não tiverem nenhuma relação
		query := `
		MATCH (preferences: Preferences{id: $preferencesId})
		WITH preferences
		UNWIND $items as item
		OPTIONAL MATCH (preferences)-->(filter :Filter{type: item.type})
		OPTIONAL MATCH (filter)-[relation :HAS]->(option:Option) WHERE NOT option.id IN (item.optionIds) DELETE relation
		WITH preferences, item, size(item.optionIds) as count
		WHERE count > 0
		MATCH (option: Option) WHERE option.id IN (item.optionIds)
		MERGE (preferences)-[:HAS]->(f: Filter{type: item.type})
		MERGE (f)-[:HAS]->(option)
		ON CREATE SET
			f.createdAt = COALESCE(item.createdAt, f.createdAt) 
		SET 
			f.optionIds = COALESCE(item.optionIds, f.optionIds),
			f.type = COALESCE(item.type, f.type),
			f.isRequired = COALESCE(item.isRequired, f.isRequired)
`

		_, err := neo4jdb.ExecuteQuery(ctx, driver, query, params)
		return err
	}
}

func Delete(driver *neo4j.DriverWithContext) DeleteFilterById {
	return func(ctx context.Context, id string) error {
		params := map[string]any{"id": id}

		query := "MATCH (n: Filter{id: $id}) DETACH DELETE n"

		_, err := neo4jdb.ExecuteQuery(ctx, driver, query, params)
		return err
	}
}
