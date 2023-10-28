package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go_ws/services/profile_service/options/models"
	neo4jdb "go_ws/shared/database/neo4j"
	"time"
)

type CreateOptionDataParams struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Type      models.OptionType `json:"type"`
	CreatedAt time.Time         `json:"createdAt"`
}

type FindOptions func(ctx context.Context) ([]models.OptionView, error)

type SaveOption func(context.Context, CreateOptionDataParams) error

type LinkOptionsWithProfile func(context.Context, string, []string) error

type DeleteOptionById func(ctx context.Context, id string) error

func FindAll(driver *neo4j.DriverWithContext) FindOptions {
	return func(ctx context.Context) ([]models.OptionView, error) {

		query := "MATCH (n: Option) RETURN n"

		results, err := neo4jdb.ExecuteWithMapping[models.OptionView](ctx, driver, query, nil)
		if err != nil {
			return nil, err
		}

		return results, nil
	}
}

func LinkWithProfile(driver *neo4j.DriverWithContext) LinkOptionsWithProfile {
	return func(ctx context.Context, profileId string, data []string) error {

		params := map[string]any{
			"profileId": profileId,
			"optionIds": data,
		}

		//TODO: melhorar update/delete das options, talvez deletar os que não tiverem nenhuma relação

		query := `
		MATCH (p :Profile{id: $profileId})
		WITH p
		OPTIONAL MATCH (p)-[relation :HAS]->(option) WHERE NOT option.id IN ($optionIds) DELETE relation
		WITH p
		MATCH (o :Option) WHERE o.id IN ($optionIds)
		MERGE (p)-[relation :HAS]->(o)
`
		_, err := neo4jdb.ExecuteQuery(ctx, driver, query, params)
		return err
	}
}

func Save(driver *neo4j.DriverWithContext) SaveOption {
	return func(ctx context.Context, data CreateOptionDataParams) error {
		params := map[string]any{
			"id":        data.ID,
			"name":      data.Name,
			"type":      data.Type,
			"createdAt": data.CreatedAt,
		}

		query := `
		MERGE (n: Option{id: $id}) 
		ON CREATE SET
			i.createdAt = COALESCE($createdAt, n.createdAt) 
		SET 
			n.name = COALESCE($name, n.name),
			n.name = COALESCE($type, n.type)`

		_, err := neo4jdb.ExecuteQuery(ctx, driver, query, params)
		return err
	}
}

func Delete(driver *neo4j.DriverWithContext) DeleteOptionById {
	return func(ctx context.Context, id string) error {
		params := map[string]any{"id": id}

		query := "MATCH (n: Option{id: $id}) DETACH DELETE n"

		_, err := neo4jdb.ExecuteQuery(ctx, driver, query, params)
		return err
	}
}
