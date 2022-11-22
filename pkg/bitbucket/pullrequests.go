package bitbucket

type PullRequest struct {
	HTMLDescription string `json:"htmlDescription,omitempty"`
	CreatedDate     int64  `json:"createdDate,omitempty"`
	ClosedDate      int64  `json:"closedDate,omitempty"`
	FromRef         Ref    `json:"fromRef,omitempty"`
	Participants    []struct {
		LastReviewedCommit string `json:"lastReviewedCommit,omitempty"`
		Approved           bool   `json:"approved,omitempty"`
		Status             string `json:"status,omitempty"`
		Role               string `json:"role,omitempty"`
		User               User   `json:"user,omitempty"`
	} `json:"participants,omitempty"`
	Reviewers   Reviewers `json:"reviewers,omitempty"`
	Description string    `json:"description,omitempty"`
	UpdatedDate int64     `json:"updatedDate,omitempty"`
	Closed      bool      `json:"closed,omitempty"`
	Title       string    `json:"title,omitempty"`
	ToRef       Ref       `json:"toRef,omitempty"`
	Version     int       `json:"version,omitempty"`
	Locked      bool      `json:"locked,omitempty"`
	ID          int       `json:"id,omitempty"`
	State       string    `json:"state,omitempty"`
	Open        bool      `json:"open,omitempty"`
	Links       struct {
	} `json:"links,omitempty"`
}
type Reviewer struct {
	LastReviewedCommit string `json:"lastReviewedCommit,omitempty"`
	Approved           bool   `json:"approved,omitempty"`
	Status             string `json:"status,omitempty"`
	Role               string `json:"role,omitempty"`
	User               User   `json:"user,omitempty"`
}

type Reviewers []Reviewer

type Ref struct {
	Repository   Repo   `json:"keyName,omitempty"`
	DisplayID    string `json:"displayId,omitempty"`
	LatestCommit string `json:"latestCommit,omitempty"`
	ID           string `json:"id,omitempty"`
	Type         string `json:"type,omitempty"`
}

type DefaultReviwers []User

type Branch struct {
	ID              string `json:"id"`
	DisplayID       string `json:"displayId"`
	Type            string `json:"type"`
	LatestCommit    string `json:"latestCommit"`
	LatestChangeset string `json:"latestChangeset"`
	IsDefault       bool   `json:"isDefault"`
}
