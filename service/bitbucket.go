package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jonaz/mgit/config"
	"github.com/jonaz/mgit/models"
)

type Bitbucket struct {
	Config config.Config
}

func NewBitbucket(config config.Config) *Bitbucket {
	return &Bitbucket{
		Config: config,
	}
}

func (b *Bitbucket) ListProjects() (models.Projects, error) {
	u := "rest/api/1.0/projects/"
	projects := models.Projects{}
	err := b.do("GET", u, nil, &projects)
	return projects, err
}

func (b *Bitbucket) ListRepos(projectKey string) (models.Repos, error) {
	u := fmt.Sprintf("rest/api/1.0/projects/%s/repos/", projectKey)
	repos := models.Repos{}
	err := b.do("GET", u, nil, &repos)
	return repos, err
}

func (b *Bitbucket) ListFiles(projectKey, repo string) (models.Files, error) {
	u := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/files/?limit=1000", projectKey, repo)
	files := models.Files{}
	err := b.do("GET", u, nil, &files)
	return files, err
}

type ErrorResponse struct {
	Errors []struct {
		Context       interface{} `json:"context"`
		Message       string      `json:"message"`
		ExceptionName string      `json:"exceptionName"`
	} `json:"errors"`
}

func (e ErrorResponse) Error() string {
	errs := []string{}
	for _, v := range e.Errors {
		errs = append(errs, v.Message)
	}
	return strings.Join(errs, ",")
}

func (b *Bitbucket) do(method, uri string, body io.Reader, response interface{}) error {
	client := &http.Client{}
	u := fmt.Sprintf("%s/%s", b.Config.BitbucketURL, uri)
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return err
	}
	req.SetBasicAuth(b.Config.BitbucketUser, b.Config.BitbucketPassword)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		decoder := json.NewDecoder(resp.Body)
		errorResponse := &ErrorResponse{}
		err := decoder.Decode(errorResponse)
		if err != nil {
			return err
		}
		return errorResponse
	}

	if response != nil {
		decoder := json.NewDecoder(resp.Body)
		return decoder.Decode(response)
	}

	return nil
}
