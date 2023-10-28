package models

type SettingsOwnerView struct {
	OwnerID  string       `json:"ownerId"`
	Settings SettingsView `json:"data"`
}
