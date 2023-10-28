package models

import (
	filter "go_ws/services/profile_service/filter/models"
	"time"
)

type PreferencesSchema struct {
	ID                          string             `json:"id"`
	MinAge                      int                `json:"minAge"`
	MaxAge                      int                `json:"maxAge"`
	Distance                    int                `json:"distance"`
	CreatedAt                   time.Time          `json:"createdAt"`
	Children                    *filter.FilterView `json:"children,omitempty"`
	Drink                       *filter.FilterView `json:"drink,omitempty"`
	Animals                     *filter.FilterView `json:"animals,omitempty"`
	AstrologicalSigns           *filter.FilterView `json:"astrologicalSigns,omitempty"`
	BiologicalSexes             *filter.FilterView `json:"biologicalSexes,omitempty"`
	GenderExpressions           *filter.FilterView `json:"genderExpressions,omitempty"`
	GenderIdentities            *filter.FilterView `json:"genderIdentities,omitempty"`
	SexualAffectiveOrientations *filter.FilterView `json:"sexualAffectiveOrientations,omitempty"`
	BodyType                    *filter.FilterView `json:"bodyType,omitempty"`
	HairLengths                 *filter.FilterView `json:"hairLengths,omitempty"`
	HairColors                  *filter.FilterView `json:"hairColors,omitempty"`
	HairTypes                   *filter.FilterView `json:"hairTypes,omitempty"`
	Languages                   *filter.FilterView `json:"languages,omitempty"`
	MusicalStyles               *filter.FilterView `json:"musicalStyles,omitempty"`
	Nourishments                *filter.FilterView `json:"nourishments,omitempty"`
	Personalities               *filter.FilterView `json:"personalities,omitempty"`
	PhysicalExercises           *filter.FilterView `json:"physicalExercises,omitempty"`
	PoliticalPositions          *filter.FilterView `json:"politicalPositions,omitempty"`
	RelationTypes               *filter.FilterView `json:"relationTypes,omitempty"`
	Religions                   *filter.FilterView `json:"religions,omitempty"`
	Smokes                      *filter.FilterView `json:"smokes,omitempty"`
	UnionTypes                  *filter.FilterView `json:"unionTypes,omitempty"`
}
