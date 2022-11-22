package bitbucket

type Files struct {
	Values        []string `json:"values"`
	Size          int      `json:"size"`
	IsLastPage    bool     `json:"isLastPage"`
	Start         int      `json:"start"`
	Limit         int      `json:"limit"`
	NextPageStart int      `json:"nextPageStart"`
}
