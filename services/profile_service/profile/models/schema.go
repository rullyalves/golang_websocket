package models

import (
	"go_ws/services/profile_service/address/models"
	image "go_ws/services/profile_service/image/models"
	options "go_ws/services/profile_service/options/models"
	"time"
)

type ProfileSchema struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Age         int                `json:"age"`
	Height      *float64           `json:"height,omitempty"`
	Description string             `json:"description"`
	CreatedAt   time.Time          `json:"createdAt"`
	Photos      []*image.ImageView `json:"photos"`

	//Preferences                *PreferencesSchema                `json:"preferences"`
	Child                      *options.OptionView   `json:"child,omitempty"`
	Drink                      *options.OptionView   `json:"drink,omitempty"`
	Animal                     *options.OptionView   `json:"animal,omitempty"`
	AstrologicalSign           *options.OptionView   `json:"astrologicalSign,omitempty"`
	BiologicalSex              *options.OptionView   `json:"biologicalSex,omitempty"`
	GenderExpression           *options.OptionView   `json:"genderExpression,omitempty"`
	GenderIdentity             *options.OptionView   `json:"genderIdentity,omitempty"`
	SexualAffectiveOrientation *options.OptionView   `json:"sexualAffectiveOrientation,omitempty"`
	BodyType                   *options.OptionView   `json:"bodyType,omitempty"`
	HairLength                 *options.OptionView   `json:"hairLength,omitempty"`
	HairColor                  *options.OptionView   `json:"hairColor,omitempty"`
	HairType                   *options.OptionView   `json:"hairType,omitempty"`
	Language                   *options.OptionView   `json:"language,omitempty"`
	MusicalStyles              []*options.OptionView `json:"musicalStyles"`
	Nourishment                *options.OptionView   `json:"nourishment,omitempty"`
	Personality                *options.OptionView   `json:"personality,omitempty"`
	PhysicalExercise           *options.OptionView   `json:"physicalExercise,omitempty"`
	PoliticalPosition          *options.OptionView   `json:"politicalPosition,omitempty"`
	Pronoun                    *options.OptionView   `json:"pronoun,omitempty"`
	RelationType               *options.OptionView   `json:"relationType,omitempty"`
	Religion                   *options.OptionView   `json:"religion,omitempty"`
	Smoke                      *options.OptionView   `json:"smoke,omitempty"`
	UnionType                  *options.OptionView   `json:"unionType,omitempty"`
	Address                    *models.AddressView   `json:"address,omitempty"`
}
