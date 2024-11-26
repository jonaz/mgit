package bitbucket

type WebhookEvent struct {
	EventKey    string   `json:"eventKey"`
	Date        string   `json:"date"`
	Repository  *Repo    `json:"repository"`
	Changes     []Change `json:"changes"`
	PullRequest struct {
		ID          int    `json:"id"`
		Version     int    `json:"version"`
		Title       string `json:"title"`
		Description string `json:"description"`
		State       string `json:"state"`
		Open        bool   `json:"open"`
		Closed      bool   `json:"closed"`
		CreatedDate int64  `json:"createdDate"`
		UpdatedDate int64  `json:"updatedDate"`
		FromRef     struct {
			ID           string `json:"id"`
			DisplayID    string `json:"displayId"`
			LatestCommit string `json:"latestCommit"`
			Type         string `json:"type"`
			Repository   Repo   `json:"repository"`
		} `json:"fromRef"`
		ToRef struct {
			ID           string `json:"id"`
			DisplayID    string `json:"displayId"`
			LatestCommit string `json:"latestCommit"`
			Type         string `json:"type"`
			Repository   Repo   `json:"repository"`
		} `json:"toRef"`
		Locked bool `json:"locked"`
		Author struct {
			User     User   `json:"user"`
			Role     string `json:"role"`
			Approved bool   `json:"approved"`
			Status   string `json:"status"`
		} `json:"author"`
		Reviewers []struct {
			User               User   `json:"user"`
			LastReviewedCommit string `json:"lastReviewedCommit,omitempty"`
			Role               string `json:"role"`
			Approved           bool   `json:"approved"`
			Status             string `json:"status"`
		} `json:"reviewers"`
		Participants []interface{} `json:"participants"`
		Links        struct {
			Self []struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
	} `json:"pullRequest"`
	PreviousTitle       string `json:"previousTitle"`
	PreviousDescription string `json:"previousDescription"`
	PreviousDraft       bool   `json:"previousDraft"`
	PreviousTarget      struct {
		ID              string `json:"id"`
		DisplayID       string `json:"displayId"`
		Type            string `json:"type"`
		LatestCommit    string `json:"latestCommit"`
		LatestChangeset string `json:"latestChangeset"`
	} `json:"previousTarget"`
}
