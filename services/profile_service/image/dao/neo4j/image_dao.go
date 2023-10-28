package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go_ws/services/profile_service/image/models"

	neo4jdb "go_ws/shared/database/neo4j"
	"time"
)

type CreateImageDataParams struct {
	ID        string    `json:"id"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"createdAt"`
}

type FindImageByProfileIdIn func(ctx context.Context, ids []string) ([]models.ImageOwnerView, error)

type SaveImages func(context.Context, []CreateImageDataParams) error

type LinkImageWithProfile func(context.Context, string, []CreateImageDataParams) error

type DeleteImageById func(ctx context.Context, id string) error

func FindProfileIdIn(driver *neo4j.DriverWithContext) FindImageByProfileIdIn {
	return func(ctx context.Context, ids []string) ([]models.ImageOwnerView, error) {
		params := map[string]any{
			"ids": ids,
		}

		query := "MATCH (p:Profile)-[h]->(a:Image) WHERE p.id IN ($ids) RETURN {ownerId: p.id, data: a}"

		results, err := neo4jdb.ExecuteWithMapping[models.ImageOwnerView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return results, nil
	}
}

func LinkWithProfile(driver *neo4j.DriverWithContext) LinkImageWithProfile {
	return func(ctx context.Context, profileId string, data []CreateImageDataParams) error {

		var items []map[string]any

		var urls []string

		for _, item := range data {
			items = append(items, map[string]any{
				"id":        item.ID,
				"url":       item.Url,
				"createdAt": item.CreatedAt,
			})

			urls = append(urls, item.Url)
		}

		params := map[string]any{
			"items":     items,
			"profileId": profileId,
			"urls":      urls,
		}

		// TODO: melhorar update/delete das imagens
		query := `
		MATCH (p :Profile{id: $profileId})
		WITH p
		OPTIONAL MATCH (p)-[relation :HAS]->(image :Image) WHERE NOT image.url IN ($urls) DELETE relation, image
		WITH p
		UNWIND $items as item 
		MERGE (p)-[:HAS]->(i:Image{id: item.id}) 
		SET i.url = item.url, i.createdAt = item.createdAt`

		_, err := neo4jdb.ExecuteQuery(ctx, driver, query, params)
		return err
	}
}

func SaveAll(driver *neo4j.DriverWithContext) SaveImages {
	return func(ctx context.Context, data []CreateImageDataParams) error {

		var items []map[string]any

		for _, item := range data {
			items = append(items, map[string]any{
				"id":        item.ID,
				"url":       item.Url,
				"createdAt": item.CreatedAt,
			})
		}

		params := map[string]any{
			"items": items,
		}

		query :=
			"UNWIND $items as item " +
				"MERGE (i:Image{id: item.id}) " +
				"SET i.url = item.url, i.createdAt = item.createdAt "

		_, err := neo4jdb.ExecuteQuery(ctx, driver, query, params)
		return err
	}
}

func Delete(driver *neo4j.DriverWithContext) DeleteImageById {
	return func(ctx context.Context, id string) error {
		params := map[string]any{"id": id}
		_, err := neo4jdb.ExecuteQuery(ctx, driver, "MATCH (i: Image{id: $id}) DETACH DELETE i", params)
		return err
	}
}
