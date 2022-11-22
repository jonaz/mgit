package bitbucket

type Changes struct {
	FromHash   interface{} `json:"fromHash"`
	ToHash     string      `json:"toHash"`
	Properties struct {
	} `json:"properties"`
	Values        []ChangeItem `json:"values"`
	Size          int          `json:"size"`
	IsLastPage    bool         `json:"isLastPage"`
	Start         int          `json:"start"`
	Limit         int          `json:"limit"`
	NextPageStart int          `json:"nextPageStart"`
}

type ChangeItem struct {
	ContentID     string `json:"contentId"`
	FromContentID string `json:"fromContentId"`
	Path          struct {
		Components []string `json:"components"`
		Parent     string   `json:"parent"`
		Name       string   `json:"name"`
		Extension  string   `json:"extension"`
		ToString   string   `json:"toString"`
	} `json:"path"`
	Executable       bool   `json:"executable"`
	PercentUnchanged int    `json:"percentUnchanged"`
	Type             string `json:"type"`
	NodeType         string `json:"nodeType"`
	SrcPath          struct {
		Components []string `json:"components"`
		Parent     string   `json:"parent"`
		Name       string   `json:"name"`
		Extension  string   `json:"extension"`
		ToString   string   `json:"toString"`
	} `json:"srcPath"`
	SrcExecutable bool `json:"srcExecutable"`
	Links         struct {
		Self []struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"links"`
	Properties struct {
		GitChangeType string `json:"gitChangeType"`
	} `json:"properties"`
}
