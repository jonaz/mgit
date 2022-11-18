package bitbucket

type ProjectsResponse struct {
	Size          int      `json:"size"`
	Limit         int      `json:"limit"`
	IsLastPage    bool     `json:"isLastPage"`
	Values        Projects `json:"values"`
	Start         int      `json:"start"`
	NextPageStart int      `json:"nextPageStart"`
}

type Projects []Project

type Project struct {
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
}
