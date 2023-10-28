package input

type CreateAddressByLocationParamsInput struct {
	UserID    string  `json:"userId"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
