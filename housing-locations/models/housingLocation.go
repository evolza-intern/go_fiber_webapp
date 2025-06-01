package models

type HousingLocation struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	City           string `json:"city"`
	State          string `json:"state"`
	Photo          string `json:"photo"`
	AvailableUnits int    `json:"availableUnits"`
	Wifi           bool   `json:"wifi"`
	Laundry        bool   `json:"laundry"`
}
