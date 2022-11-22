package bitbucket

type User struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	EmailAddress string `json:"emailAddress,omitempty"`
	Slug         string `json:"slug,omitempty"`
	Active       bool   `json:"active,omitempty"`
	Links        struct {
	} `json:"links,omitempty"`
	Type        string `json:"type,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	AvatarURL   string `json:"avatarUrl,omitempty"`
}
