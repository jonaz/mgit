package bitbucket

import (
	"fmt"
	"path/filepath"
	"strings"
)

type ReposResponse struct {
	Size          int   `json:"size"`
	Limit         int   `json:"limit"`
	IsLastPage    bool  `json:"isLastPage"`
	Values        Repos `json:"values"`
	Start         int   `json:"start"`
	NextPageStart int   `json:"nextPageStart"`
}
type Repos []Repo

type Repo struct {
	Slug          string  `json:"slug"`
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	HierarchyID   string  `json:"hierarchyId"`
	ScmID         string  `json:"scmId"`
	State         string  `json:"state"`
	StatusMessage string  `json:"statusMessage"`
	Forkable      bool    `json:"forkable"`
	Project       Project `json:"project"`
	Public        bool    `json:"public"`
	Links         struct {
		Clone CloneLinks `json:"clone"`
		Self  []struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"links"`
}

func (r Repo) RepoPath(dst string) string {
	return filepath.Join(dst, fmt.Sprintf("%s_%s", strings.ToLower(r.Project.Key), r.Slug))
}

type CloneLinks []struct {
	Href string `json:"href"`
	Name string `json:"name"`
}

func (cl CloneLinks) GetSSH() string {
	for _, v := range cl {
		if v.Name == "ssh" {
			return v.Href
		}
	}
	return ""
}
