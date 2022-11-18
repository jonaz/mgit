package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
)

type Client struct {
	URL      string
	Username string
	Password string
	Dry      bool
}

func NewClient(url, username, password string, dry bool) *Client {
	return &Client{
		URL:      url,
		Username: username,
		Password: password,
		Dry:      dry,
	}
}

func (b *Client) ListProjects() (Projects, error) {
	var projects Projects

	start := 0
	isLastPage := false
	for !isLastPage {
		resp := ProjectsResponse{}
		u := fmt.Sprintf("rest/api/1.0/projects/?limit=100&start=%d", start)
		err := b.do("GET", u, nil, &resp)
		if err != nil {
			return projects, err
		}
		start = resp.NextPageStart
		isLastPage = resp.IsLastPage
		projects = append(projects, resp.Values...)
	}
	return projects, nil
}

func (b *Client) ListRepos(projectKey string) (Repos, error) {
	var repos Repos
	start := 0
	isLastPage := false
	for !isLastPage {
		resp := ReposResponse{}
		u := fmt.Sprintf("rest/api/1.0/projects/%s/repos/?limit=100&start=%d", projectKey, start)
		err := b.do("GET", u, nil, &resp)
		if err != nil {
			return repos, err
		}
		start = resp.NextPageStart
		isLastPage = resp.IsLastPage
		repos = append(repos, resp.Values...)
	}
	return repos, nil
}

// GetRepo fetches a repo.
func (b *Client) GetRepo(projectKey, r string) (Repo, error) {
	u := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s", projectKey, r)
	repo := Repo{}
	err := b.do("GET", u, nil, &repo)
	return repo, err
}

// GetFileContent fetches the whole file from bitbucket API.
func (b *Client) GetFileContent(projectKey, repo, file, atRev string) ([]byte, error) {
	u := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/raw/%s", projectKey, repo, file)
	if atRev != "" {
		u += "?at=" + atRev
	}
	var content []byte
	err := b.do("GET", u, nil, &content)
	return content, err
}

func (b *Client) ListFiles(projectKey, repo, atRev string) ([]string, error) {
	var files []string
	start := 0
	isLastPage := false
	for !isLastPage {
		resp := Files{}
		u := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/files/?limit=1000&start=%d", projectKey, repo, start)
		if atRev != "" {
			u += "&at=" + atRev
		}
		err := b.do("GET", u, nil, &resp)
		if err != nil {
			return files, err
		}
		start = resp.NextPageStart
		isLastPage = resp.IsLastPage
		files = append(files, resp.Values...)
	}
	return files, nil
}

func (b *Client) ListChanges(projectKey, repo, commit string) ([]ChangeItem, error) {
	var files []ChangeItem
	start := 0
	isLastPage := false
	for !isLastPage {
		resp := Changes{}
		u := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/changes/?until=%s&limit=1000&until=%d", projectKey, repo, commit, start)
		err := b.do("GET", u, nil, &resp)
		if err != nil {
			return files, err
		}
		start = resp.NextPageStart
		isLastPage = resp.IsLastPage
		files = append(files, resp.Values...)
	}
	return files, nil
}

func (b *Client) ListWebhooks(projectSlug, repoSlug string) (Webhooks, error) {
	resp := WebhooksResponse{}
	u := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/webhooks?limit=100", projectSlug, repoSlug)
	err := b.do("GET", u, nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Values, nil
}
func (b *Client) CreateOrUpdateWebhook(projectSlug, repoSlug string, hook Webhook) error {
	webhooks, err := b.ListWebhooks(projectSlug, repoSlug)
	if err != nil {
		return err
	}

	var foundHook *Webhook
	for _, v := range webhooks {
		h := v
		if h.Name == hook.Name {
			foundHook = &h
			break
		}
	}

	if foundHook == nil {
		return b.CreateWebhook(projectSlug, repoSlug, hook)
	}

	if !foundHook.Equal(hook) {
		hook.ID = foundHook.ID
		return b.UpdateWebhook(projectSlug, repoSlug, hook)
	}
	return nil
}

func (b *Client) CreateWebhook(projectSlug, repoSlug string, hook Webhook) error {
	u := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/webhooks", projectSlug, repoSlug)
	jsonBody, err := json.Marshal(hook)
	if err != nil {
		return err
	}
	logrus.Infof("bitbucket: creating webhhook %s for repo %s", hook.Name, path.Join(projectSlug, repoSlug))

	return b.do("POST", u, bytes.NewBuffer(jsonBody), nil)
}

func (b *Client) UpdateWebhook(projectSlug, repoSlug string, hook Webhook) error {
	u := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/webhooks/%d", projectSlug, repoSlug, hook.ID)
	jsonBody, err := json.Marshal(hook)
	if err != nil {
		return err
	}

	logrus.Infof("bitbucket: updating webhhook %s for repo %s", hook.Name, path.Join(projectSlug, repoSlug))
	return b.do("PUT", u, bytes.NewBuffer(jsonBody), nil)
}

func (b *Client) CreatePullRequest(projectSlug, repoSlug string, pr PullRequest) error {
	u := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/pull-requests", projectSlug, repoSlug)
	jsonBody, err := json.Marshal(pr)
	if err != nil {
		return err
	}
	logrus.Infof("bitbucket: creating pull-request %s for repo %s", pr.Title, path.Join(projectSlug, repoSlug))

	return b.do("POST", u, bytes.NewBuffer(jsonBody), nil)
}

// GetDefaultReviwers fetches defaul reviwers when making a PR.
func (b *Client) GetDefaultReviwers(projectKey, slug string, repoId int, sourceRef, targetRef string) (DefaultReviwers, error) {
	u := fmt.Sprintf("rest/default-reviewers/1.0/projects/%s/repos/%s/reviewers?targetRepoId=%d&sourceRepoId=%d&targetRefId=%s&sourceRefId=%s",
		projectKey, slug, repoId, repoId, targetRef, sourceRef)
	dr := DefaultReviwers{}
	err := b.do("GET", u, nil, &dr)
	return dr, err
}

func (b *Client) GetDefaultBranch(projectKey, slug string) (Branch, error) {
	u := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/default-branch", projectKey, slug)
	dr := Branch{}
	err := b.do("GET", u, nil, &dr)
	return dr, err
}

type ResponseError struct {
	Status int
	Errors []struct {
		Context       interface{} `json:"context"`
		Message       string      `json:"message"`
		ExceptionName string      `json:"exceptionName"`
	} `json:"errors"`
}

func (e ResponseError) Error() string {
	errs := []string{}
	for _, v := range e.Errors {
		errs = append(errs, v.Message)
	}
	return strings.Join(errs, ",")
}

func (b *Client) do(method, uri string, body io.Reader, response interface{}) error {
	client := &http.Client{}
	u := fmt.Sprintf("%s/%s", b.URL, uri)
	// logrus.Debugf("bitbucket %s to: %s", method, u)

	if b.Dry {
		if method != "GET" {
			logrus.Warnf("bitbucket dryrun %s URL: %s", method, u)
			return nil
		}
		// Do GETs but log them if in dry mode
		logrus.Warnf("bitbucket dryrun %s URL: %s", method, u)
	} else {
		logrus.Debugf("bitbucket %s URL: %s", method, u)
	}

	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return err
	}
	req.SetBasicAuth(b.Username, b.Password)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("bitucket: error doing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		decoder := json.NewDecoder(resp.Body)
		errorResponse := &ResponseError{}
		err = decoder.Decode(errorResponse)
		if err != nil {
			return fmt.Errorf("got status code %d: %w", resp.StatusCode, err)
		}
		errorResponse.Status = resp.StatusCode
		return errorResponse
	}

	if response != nil {
		switch v := response.(type) {
		case *[]byte:
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("bitucket: error reading response body: %w", err)
			}
			*v = b
		default:
			decoder := json.NewDecoder(resp.Body)
			return decoder.Decode(response)
		}
	}

	return nil
}
