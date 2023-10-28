package models

import "time"

type OptionType string

const (
	OptionTypeReligion                   OptionType = "religion"
	OptionTypeAnimal                     OptionType = "animal"
	OptionTypeBodyType                   OptionType = "bodyType"
	OptionTypeChild                      OptionType = "child"
	OptionTypeDrink                      OptionType = "drink"
	OptionTypeHairColor                  OptionType = "hairColor"
	OptionTypeHairLength                 OptionType = "hairLength"
	OptionTypeHairType                   OptionType = "hairType"
	OptionTypeLanguage                   OptionType = "language"
	OptionTypeMusicalStyle               OptionType = "musicalStyle"
	OptionTypeNourishment                OptionType = "nourishment"
	OptionTypeOccupation                 OptionType = "occupation"
	OptionTypePersonality                OptionType = "personality"
	OptionTypePhysicalExercise           OptionType = "physicalExercise"
	OptionTypePoliticalPosition          OptionType = "politicalPosition"
	OptionTypePronoun                    OptionType = "pronoun"
	OptionTypeRelationType               OptionType = "relationType"
	OptionTypeSmoke                      OptionType = "smoke"
	OptionTypeUnionType                  OptionType = "unionType"
	OptionTypeAstrologicalSign           OptionType = "astrologicalSign"
	OptionTypeSexualAffectiveOrientation OptionType = "sexualAffectiveOrientation"
	OptionTypeBiologicalSex              OptionType = "biologicalSex"
	OptionTypeGenderExpression           OptionType = "genderExpression"
	OptionTypeGenderIdentity             OptionType = "genderIdentity"
)

type OptionView struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Type      OptionType `json:"type"`
	CreatedAt time.Time  `json:"createdAt"`
}
