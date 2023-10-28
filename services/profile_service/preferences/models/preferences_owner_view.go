package models

type PreferencesOwnerView struct {
	OwnerID     string          `json:"owner_id"`
	Preferences PreferencesView `json:"data"`
}
