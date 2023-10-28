package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type GeocodingResult struct {
	City                     string  `json:"city"`
	PrincipalSubdivisionCode string  `json:"principalSubdivisionCode"`
	PrincipalSubdivision     string  `json:"principalSubdivision"`
	Locality                 string  `json:"locality"`
	Latitude                 float64 `json:"latitude"`
	Longitude                float64 `json:"longitude"`
}

func (g GeocodingResult) State() string {
	return strings.Split(g.PrincipalSubdivisionCode, "-")[1]
}

func (g GeocodingResult) Neighborhood() string {
	return strings.Split(g.Locality, " ")[0]
}

type FindAddressByLatLng func(latitude float64, longitude float64) (*GeocodingResult, error)

func FindAddress(latitude float64, longitude float64) (*GeocodingResult, error) {

	url := fmt.Sprintf(os.Getenv("GEOCODING_API_URL"), latitude, longitude)

	//TODO: abstrair http client no futuro

	result, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	var data GeocodingResult

	err = json.NewDecoder(result.Body).Decode(&data)

	if err != nil {
		return nil, err
	}

	return &data, nil
}
