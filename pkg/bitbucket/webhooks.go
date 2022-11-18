package bitbucket

type WebhooksResponse struct {
	Size          int      `json:"size"`
	Limit         int      `json:"limit"`
	IsLastPage    bool     `json:"isLastPage"`
	Values        Webhooks `json:"values"`
	Start         int      `json:"start"`
	NextPageStart int      `json:"nextPageStart"`
}
type Webhooks []Webhook

type Webhook struct {
	ID          int      `json:"id,omitempty"`
	Name        string   `json:"name"`
	CreatedDate int64    `json:"createdDate,omitempty"`
	UpdatedDate int64    `json:"updatedDate,omitempty"`
	Events      []string `json:"events"`
	// Configuration struct {
	// CreatedBy string `json:"createdBy,omitempty"`
	// } `json:"configuration,omitempty"`
	URL    string `json:"url"`
	Active bool   `json:"active"`
}

func (w Webhook) Equal(hook Webhook) bool {
	if w.URL != hook.URL {
		return false
	}
	if w.Active != hook.Active {
		return false
	}
	return AssertSameStringSlice(w.Events, hook.Events)
}

func AssertSameStringSlice(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}

	itemAppearsTimes := make(map[string]int, len(x))
	for _, i := range x {
		itemAppearsTimes[i]++
	}

	for _, i := range y {
		if _, ok := itemAppearsTimes[i]; !ok {
			return false
		}

		itemAppearsTimes[i]--

		if itemAppearsTimes[i] == 0 {
			delete(itemAppearsTimes, i)
		}
	}

	if len(itemAppearsTimes) == 0 {
		return true
	}

	return false
}
