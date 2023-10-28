package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go_ws/services/profile_service/profile/models"
	"go_ws/shared/collections"
	neo4jdb "go_ws/shared/database/neo4j"
	"time"
)

type CreateProfileDataParams struct {
	ID          string     `json:"id"`
	Name        *string    `json:"name"`
	Age         *int       `json:"age"`
	Height      *float32   `json:"height"`
	Description *string    `json:"description"`
	CreatedAt   *time.Time `json:"createdAt"`
}

type UpdateProfileDataParams struct {
	Name        string  `json:"name"`
	Age         int     `json:"age"`
	Height      float32 `json:"height"`
	Description string  `json:"description"`
}

type FindSenderProfileByBlockIdIn func(context.Context, []string) (*[]models.ProfileOwnerView, error)

type FindReceiverProfileByBlockIdIn func(context.Context, []string) (*[]models.ProfileOwnerView, error)

type FindSenderProfileByComplaintIdIn func(context.Context, []string) (*[]models.ProfileOwnerView, error)

type FindReceiverProfileByComplaintIdIn func(context.Context, []string) (*[]models.ProfileOwnerView, error)

type FindReceiverProfileByDeliveryIdIn func(context.Context, []string) (*[]models.ProfileOwnerView, error)

type FindSenderProfileByDislikeIdIn func(context.Context, []string) (*[]models.ProfileOwnerView, error)

type FindReceiverProfileByDislikeIdIn func(context.Context, []string) (*[]models.ProfileOwnerView, error)

type FindSenderProfileByLikeIdIn func(context.Context, []string) (*[]models.ProfileOwnerView, error)

type FindReceiverProfileByLikeIdIn func(context.Context, []string) (*[]models.ProfileOwnerView, error)

type FindSenderProfileByMessageIdIn func(context.Context, []string) (*[]models.ProfileOwnerView, error)

type FindProfilesByChatIdIn func(context.Context, []string) (*[]models.ProfileOwnerView, error)

type FindInterestProfilesByProfileId func(context.Context, string, int) (*[]models.ProfileView, error)

type FindInterestedProfilesByProfileId func(context.Context, string, int) (*[]models.ProfileView, error)

type FindParticipantIdsByChatIdIn func(context.Context, []string) (*[]string, error)

type FindProfileByIdIn func(context.Context, []string) (*[]models.ProfileView, error)

type FindProfileNameById func(context.Context, string) (*string, error)

type SaveProfile func(context.Context, CreateProfileDataParams) error

type DeleteProfileById func(context.Context, string) error

func FindSenderByBlockIdIn(driver *neo4j.DriverWithContext) FindSenderProfileByBlockIdIn {
	return func(ctx context.Context, ids []string) (*[]models.ProfileOwnerView, error) {

		params := map[string]any{
			"ids": ids,
		}

		query :=
			"MATCH (b :Block) WHERE b.id IN ($ids) " +
				"MATCH (b)-->(p: Profile) " +
				"RETURN {ownerId: b.id, data: p}"

		results, err := neo4jdb.ExecuteWithMapping[models.ProfileOwnerView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return &results, nil
	}
}

func FindReceiverByBlockIdIn(driver *neo4j.DriverWithContext) FindReceiverProfileByBlockIdIn {
	return func(ctx context.Context, ids []string) (*[]models.ProfileOwnerView, error) {

		params := map[string]any{
			"ids": ids,
		}

		query :=
			"MATCH (b :Block) WHERE b.id IN ($ids) " +
				"MATCH (b)<--(p: Profile) " +
				"RETURN {ownerId: b.id, data: p}"

		results, err := neo4jdb.ExecuteWithMapping[models.ProfileOwnerView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return &results, nil
	}
}

func FindReceiverByComplaintIdIn(driver *neo4j.DriverWithContext) FindReceiverProfileByComplaintIdIn {
	return func(ctx context.Context, ids []string) (*[]models.ProfileOwnerView, error) {

		params := map[string]any{
			"ids": ids,
		}

		query :=
			"MATCH (c :Complaint) WHERE c.id IN ($ids) " +
				"MATCH (c)<--(p: Profile) " +
				"RETURN {ownerId: c.id, data: p}"

		results, err := neo4jdb.ExecuteWithMapping[models.ProfileOwnerView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return &results, nil
	}
}

func FindSenderByComplaintIdIn(driver *neo4j.DriverWithContext) FindSenderProfileByComplaintIdIn {
	return func(ctx context.Context, ids []string) (*[]models.ProfileOwnerView, error) {

		params := map[string]any{
			"ids": ids,
		}

		query :=
			"MATCH (c :Complaint) WHERE c.id IN ($ids) " +
				"MATCH (c)-->(p: Profile) " +
				"RETURN {ownerId: c.id, data: p}"

		results, err := neo4jdb.ExecuteWithMapping[models.ProfileOwnerView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return &results, nil
	}
}

func FindReceiverByDeliveryIdIn(driver *neo4j.DriverWithContext) FindReceiverProfileByDeliveryIdIn {
	return func(ctx context.Context, ids []string) (*[]models.ProfileOwnerView, error) {

		params := map[string]any{
			"ids": ids,
		}

		query :=
			"MATCH (d:Delivery) WHERE d.id IN ($ids) " +
				"MATCH (d)<--(p :Profile) " +
				"RETURN {ownerId: d.id, data: p}"

		results, err := neo4jdb.ExecuteWithMapping[models.ProfileOwnerView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return &results, nil
	}
}

func FindSenderByDislikeIdIn(driver *neo4j.DriverWithContext) FindSenderProfileByDislikeIdIn {
	return func(ctx context.Context, ids []string) (*[]models.ProfileOwnerView, error) {

		params := map[string]any{
			"ids": ids,
		}

		query :=
			"MATCH (d:Dislike) WHERE d.id IN ($ids) " +
				"MATCH (d)-->(p: Profile) " +
				"RETURN {ownerId: d.id, data: p}"

		results, err := neo4jdb.ExecuteWithMapping[models.ProfileOwnerView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return &results, nil
	}
}

func FindReceiverByDislikeIdIn(driver *neo4j.DriverWithContext) FindReceiverProfileByDislikeIdIn {
	return func(ctx context.Context, ids []string) (*[]models.ProfileOwnerView, error) {

		params := map[string]any{
			"ids": ids,
		}

		query :=
			"MATCH (d:Dislike) WHERE d.id IN ($ids) " +
				"MATCH (d)<--(p: Profile) " +
				"RETURN {ownerId: d.id, data: p}"

		results, err := neo4jdb.ExecuteWithMapping[models.ProfileOwnerView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return &results, nil
	}
}

func FindSenderByLikeIdIn(driver *neo4j.DriverWithContext) FindSenderProfileByLikeIdIn {
	return func(ctx context.Context, ids []string) (*[]models.ProfileOwnerView, error) {

		params := map[string]any{
			"ids": ids,
		}

		query :=
			"MATCH (l:Like) WHERE l.id IN ($ids) " +
				"MATCH (l)<-[:SEND]-(p: Profile) " +
				"RETURN {ownerId: l.id, data: p}"

		results, err := neo4jdb.ExecuteWithMapping[models.ProfileOwnerView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return &results, nil
	}
}

func FindReceiverByLikeIdIn(driver *neo4j.DriverWithContext) FindReceiverProfileByLikeIdIn {
	return func(ctx context.Context, ids []string) (*[]models.ProfileOwnerView, error) {

		params := map[string]any{
			"ids": ids,
		}

		query :=
			"MATCH (l:Like) WHERE l.id IN ($ids) " +
				"MATCH (l)-[:RECEIVE]->(p: Profile) " +
				"RETURN {ownerId: l.id, data: p}"

		results, err := neo4jdb.ExecuteWithMapping[models.ProfileOwnerView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return &results, nil
	}
}

func FindSenderByMessageIdIn(driver *neo4j.DriverWithContext) FindSenderProfileByMessageIdIn {
	return func(ctx context.Context, ids []string) (*[]models.ProfileOwnerView, error) {

		params := map[string]any{
			"ids": ids,
		}

		query :=
			"MATCH (m:Message) WHERE m.id IN ($ids) " +
				"MATCH (p: Profile)-[:SEND]->(m) RETURN m.id as ownerId, p"

		results, err := neo4jdb.ExecuteWithMapping[models.ProfileOwnerView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return &results, nil
	}
}

func FindIdsByChatIdIn(driver *neo4j.DriverWithContext) FindParticipantIdsByChatIdIn {
	return func(ctx context.Context, ids []string) (*[]string, error) {

		params := map[string]any{
			"ids": ids,
		}

		query := "MATCH (c:Chat) WHERE c.id IN ($ids) " +
			"MATCH (c)<--(p: Profile) " +
			"RETURN {ids: collect(p.id)}"

		results, err := neo4jdb.ExecuteWithMapping[map[string]any](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		data := collections.FromList[string](results[0]["ids"].([]any))

		return &data, nil
	}
}

func FindByChatIdIn(driver *neo4j.DriverWithContext) FindProfilesByChatIdIn {
	return func(ctx context.Context, ids []string) (*[]models.ProfileOwnerView, error) {

		params := map[string]any{
			"ids": ids,
		}

		query := "MATCH (c:Chat)<--(p: Profile) WHERE c.id IN ($ids) RETURN c.id as ownerId, p"

		results, err := neo4jdb.ExecuteWithMapping[models.ProfileOwnerView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return &results, nil
	}
}

func FindInterestProfilesIn(driver *neo4j.DriverWithContext) FindInterestProfilesByProfileId {
	return func(ctx context.Context, profileId string, limit int) (*[]models.ProfileView, error) {

		params := map[string]any{
			"id":    profileId,
			"limit": limit,
		}

		//TODO: adicionei um "OR" option is NULL para caso exista um filtro e ele não possua relações, dai ele é ignorado
		query := `
		MATCH 
		(p: Profile{id: $id}),
		(p)-->(address :Location),
		(p)-->(pref :Preferences)

		MATCH (p2 :Profile)-->(p2address :Location) 

		WITH *, 
		point({latitude: address.latitude, longitude: address.longitude}) as pLocation,
		point({latitude: p2address.latitude, longitude: p2address.longitude}) as p2Location

		WHERE point.distance(pLocation, p2Location) <= pref.distance

		AND (p2.age >= pref.minAge AND p2.age <= pref.maxAge)
		AND NOT EXISTS((p)--(:Dislike|Match|Block)--(p2))
		AND NOT EXISTS((p)-->(:Like)-->(p2))
		WITH p, p2, pref
			OPTIONAL MATCH (pref)--(filter: Filter)--(option: Option)
			OPTIONAL MATCH (option)--(filtered: Profile) WHERE filtered <> p
				WITH
					CASE
						WHEN ((filter IS NULL OR filter.isRequired = FALSE OR option IS NULL) AND (filtered IS NULL))
					THEN p2
					ELSE filtered
				END as p2, pref
		WHERE p <> p2
		RETURN DISTINCT p2 LIMIT $limit`

		results, err := neo4jdb.ExecuteWithMapping[models.ProfileView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return &results, nil
	}
}

func FindInterestedProfilesIn(driver *neo4j.DriverWithContext) FindInterestedProfilesByProfileId {
	return func(ctx context.Context, profileId string, limit int) (*[]models.ProfileView, error) {

		params := map[string]any{
			"id":    profileId,
			"limit": limit,
		}

		query :=
			"MATCH (p)-->(:Like)-->(p2:Profile{id: $id}) " +
				"WHERE NOT (p)--(:Dislike|Match|Block)--(p2) " +
				"RETURN p LIMIT $limit"

		results, err := neo4jdb.ExecuteWithMapping[models.ProfileView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return &results, nil
	}
}

func FindNameById(driver *neo4j.DriverWithContext) FindProfileNameById {
	return func(ctx context.Context, id string) (*string, error) {

		params := map[string]any{
			"id": id,
		}

		query := "MATCH (p :Profile{id: $id}) RETURN {name: p.name}"

		results, err := neo4jdb.ExecuteWithMapping[map[string]any](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		if len(results) == 0 {
			return nil, nil
		}

		data := results[0]["name"].(string)

		return &data, nil
	}
}

func FindByIdIn(driver *neo4j.DriverWithContext) FindProfileByIdIn {
	return func(ctx context.Context, ids []string) (*[]models.ProfileView, error) {

		params := map[string]any{
			"ids": ids,
		}

		query := "MATCH (p :Profile) WHERE p.id IN ($ids) RETURN p"

		results, err := neo4jdb.ExecuteWithMapping[models.ProfileView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		if len(results) == 0 {
			return nil, nil
		}

		return &results, nil
	}
}

func Save(driver *neo4j.DriverWithContext) SaveProfile {
	return func(ctx context.Context, data CreateProfileDataParams) error {

		params := map[string]any{
			"id":          data.ID,
			"name":        data.Name,
			"age":         data.Age,
			"height":      data.Height,
			"description": data.Description,
			"createdAt":   data.CreatedAt,
		}

		query := `
		MERGE (p :Profile{id: $id}) 
		ON CREATE SET
			p.createdAt = COALESCE($createdAt, p.createdAt) 
		SET 
			p.name = COALESCE($name, p.name),
			p.age = COALESCE($age, p.age),
			p.height = COALESCE($height, p.height),
			p.description = COALESCE($description, p.description)`

		_, err := neo4jdb.ExecuteQuery(ctx, driver, query, params)

		return err
	}
}

func Delete(driver *neo4j.DriverWithContext) DeleteProfileById {
	return func(ctx context.Context, id string) error {
		params := map[string]any{"id": id}
		_, err := neo4jdb.ExecuteQuery(ctx, driver, "MATCH (p :Profile{id: $id}) DETACH DELETE p", params)
		return err
	}
}
