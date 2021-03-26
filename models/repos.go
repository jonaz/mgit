package models

import (
	"fmt"
	"path/filepath"
)

type Repos struct {
	Size          int    `json:"size"`
	Limit         int    `json:"limit"`
	IsLastPage    bool   `json:"isLastPage"`
	Values        []Repo `json:"values"`
	Start         int    `json:"start"`
	NextPageStart int    `json:"nextPageStart"`
}

type Repo struct {
	Slug          string `json:"slug"`
	ID            int    `json:"id"`
	Name          string `json:"name"`
	HierarchyID   string `json:"hierarchyId"`
	ScmID         string `json:"scmId"`
	State         string `json:"state"`
	StatusMessage string `json:"statusMessage"`
	Forkable      bool   `json:"forkable"`
	Project       struct {
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
	} `json:"project"`
	Public bool `json:"public"`
	Links  struct {
		Clone CloneLinks `json:"clone"`
		Self  []struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"links"`
}

func (r Repo) RepoPath(dst string) string {
	return filepath.Join(dst, fmt.Sprintf("%s_%s", r.Project.Key, r.Slug))
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
