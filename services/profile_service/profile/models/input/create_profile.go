package input

type CreateUserProfileParamsInput struct {
	Username string                    `json:"username"`
	Password string                    `json:"password"`
	Profile  *CreateProfileParamsInput `json:"profile"`
}

type CreateProfileParamsInput struct {
	Name                         string   `json:"name"`
	Age                          int      `json:"age"`
	Height                       float64  `json:"height"`
	Description                  string   `json:"description"`
	ChildID                      *string  `json:"childId,omitempty"`
	SmokeID                      *string  `json:"smokeId,omitempty"`
	DrinkID                      *string  `json:"drinkId,omitempty"`
	PronounID                    *string  `json:"pronounId,omitempty"`
	UnionTypeID                  *string  `json:"unionTypeId,omitempty"`
	ReligionID                   *string  `json:"religionId,omitempty"`
	MusicalStyleIds              []string `json:"musicalStyleIds,omitempty"`
	PersonalityID                *string  `json:"personalityId,omitempty"`
	NourishmentID                *string  `json:"nourishmentId,omitempty"`
	PoliticalPositionID          *string  `json:"politicalPositionId,omitempty"`
	RelationTypeID               *string  `json:"relationTypeId,omitempty"`
	GenderIdentityID             *string  `json:"genderIdentityId,omitempty"`
	GenderExpressionID           *string  `json:"genderExpressionId,omitempty"`
	BiologicalSexID              *string  `json:"biologicalSexId,omitempty"`
	SexualAffectiveOrientationID *string  `json:"sexualAffectiveOrientationId,omitempty"`
	AstrologicalSignID           *string  `json:"astrologicalSignId,omitempty"`
	ImageUrls                    []string `json:"imageUrls,omitempty"`
}
