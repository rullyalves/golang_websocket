package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go_ws/services/relationship_service/complaint/models"

	neo4jdb "go_ws/shared/database/neo4j"
	"time"
)

type CreateComplaintCategoryDataParams struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	ParentID  *string   `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
}

type FindAllComplaintCategories func(ctx context.Context) ([]models.ComplaintView, error)

type FindComplaintCategoryByComplaintIdIn func(ctx context.Context, ids []string) ([]models.ComplaintView, error)

type FindComplaintCategoryByParentIdIn func(ctx context.Context, ids []string) ([]models.ComplaintView, error)

type SaveComplaintCategory func(context.Context, CreateComplaintCategoryDataParams) error

type DeleteComplaintCategoryById func(ctx context.Context, id string) error

func FindByComplaintIdIn(driver *neo4j.DriverWithContext) FindComplaintCategoryByComplaintIdIn {
	return func(ctx context.Context, ids []string) ([]models.ComplaintView, error) {

		params := map[string]any{
			"ids": ids,
		}

		query :=
			"MATCH (c: Complaint) WHERE c.id IN ($ids) " +
				"MATCH (c)-[:BELONGS_TO]->(cc: ComplaintCategory) " +
				"RETURN c.id as ownerId, cc"

		results, err := neo4jdb.ExecuteWithMapping[models.ComplaintView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return results, nil
	}
}

func FindByParentIdIn(driver *neo4j.DriverWithContext) FindComplaintCategoryByParentIdIn {
	return func(ctx context.Context, ids []string) ([]models.ComplaintView, error) {
		params := map[string]any{
			"ids": ids,
		}

		query :=
			"MATCH (ccParent: ComplaintCategory) WHERE ccParent.id IN ($ids) " +
				"MATCH (ccParent)-[:PARENT_OF]->(cc: ComplaintCategory) " +
				"RETURN ccParent.id as ownerId, cc"

		results, err := neo4jdb.ExecuteWithMapping[models.ComplaintView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return results, nil
	}
}

func FindAll(driver *neo4j.DriverWithContext) FindAllComplaintCategories {
	return func(ctx context.Context) ([]models.ComplaintView, error) {

		query := "MATCH (cc: ComplaintCategory) WHERE (cc)-[:PARENT_OF]->(:ComplaintCategory) RETURN cc"

		results, err := neo4jdb.ExecuteWithMapping[models.ComplaintView](ctx, driver, query, nil)
		if err != nil {
			return nil, err
		}

		return results, nil
	}
}

func Save(driver *neo4j.DriverWithContext) SaveComplaintCategory {
	return func(ctx context.Context, data CreateComplaintCategoryDataParams) error {

		params := map[string]any{
			"id":        data.ID,
			"name":      data.Name,
			"parentId":  data.ParentID,
			"createdAt": data.CreatedAt,
		}

		query := `
		MATCH (ccParent: ComplaintCategory{id: $parentId})
		WITH ccParent 
		MERGE (ccParent)-[:PARENT_OF]->(cc: ComplaintCategory{id: $id})
		ON CREATE SET
			cc.createdAt = COALESCE($createdAt, cc.createdAt)
		SET 
			cc.name = COALESCE($name, cc.name)`

		_, err := neo4jdb.ExecuteQuery(ctx, driver, query, params)
		return err
	}

}

func Delete(driver *neo4j.DriverWithContext) DeleteComplaintCategoryById {
	return func(ctx context.Context, id string) error {
		params := map[string]any{"id": id}
		_, err := neo4jdb.ExecuteQuery(ctx, driver, "MATCH (i: ComplaintCategory{id: $id}) DETACH DELETE i", params)
		return err
	}
}
