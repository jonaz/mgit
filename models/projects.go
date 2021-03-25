package models

type Projects struct {
	Size       int  `json:"size"`
	Limit      int  `json:"limit"`
	IsLastPage bool `json:"isLastPage"`
	Values     []struct {
		Key    string `json:"key"`
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Public bool   `json:"public"`
		Type   string `json:"type"`
		Links  struct {
			Self []struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
		Description string `json:"description,omitempty"`
	} `json:"values"`
	Start int `json:"start"`
}
