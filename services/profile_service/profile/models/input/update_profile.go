package input

type UpdateProfileParamsInput struct {
	Name                         *string   `json:"name,omitempty"`
	Description                  *string   `json:"description,omitempty"`
	AnimalID                     *string   `json:"animalId,omitempty"`
	AstrologicalSignID           *string   `json:"astrologicalSignId,omitempty"`
	BiologicalSexID              *string   `json:"biologicalSexId,omitempty"`
	BodyTypeID                   *string   `json:"bodyTypeId,omitempty"`
	ChildID                      *string   `json:"childId,omitempty"`
	DrinkID                      *string   `json:"drinkId,omitempty"`
	GenderExpressionID           *string   `json:"genderExpressionId,omitempty"`
	GenderIdentityID             *string   `json:"genderIdentityId,omitempty"`
	HairColorID                  *string   `json:"hairColorId,omitempty"`
	HairLengthID                 *string   `json:"hairLengthId,omitempty"`
	HairTypeID                   *string   `json:"hairTypeId,omitempty"`
	LanguageID                   *string   `json:"languageId,omitempty"`
	MusicalStyleIds              []*string `json:"musicalStyleIds,omitempty"`
	NourishmentID                *string   `json:"nourishmentId,omitempty"`
	PersonalityID                *string   `json:"personalityId,omitempty"`
	PhysicalExerciseID           *string   `json:"physicalExerciseId,omitempty"`
	PoliticalPositionID          *string   `json:"politicalPositionId,omitempty"`
	RelationTypeID               *string   `json:"relationTypeId,omitempty"`
	ReligionID                   *string   `json:"religionId,omitempty"`
	SexualAffectiveOrientationID *string   `json:"sexualAffectiveOrientationId,omitempty"`
	SmokeID                      *string   `json:"smokeId,omitempty"`
	UnionTypeID                  *string   `json:"unionTypeId,omitempty"`
	PronounID                    *string   `json:"pronounId,omitempty"`
	ImageUrls                    []*string `json:"imageUrls,omitempty"`
}
