package providers

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/jonaz/mgit/pkg/bitbucket"
	"github.com/jonaz/mgit/pkg/git"
	"github.com/jonaz/mgit/pkg/utils"
	"github.com/sirupsen/logrus"
)

type Bitbucket struct {
	DefaultProvider
	BitbucketURL string
}

func NewBitbucket(dir, url string) *Bitbucket {
	return &Bitbucket{
		DefaultProvider: DefaultProvider{Dir: dir},
		BitbucketURL:    url,
	}
}

func (b *Bitbucket) Clone(whitelist []string, hasFile, contentRegexp string) error {
	username, password, err := utils.Credentials()
	if err != nil {
		return err
	}

	bit := bitbucket.NewClient(b.BitbucketURL, username, password, false)

	contentReg, err := regexp.Compile(contentRegexp)
	if err != nil {
		return err
	}

	projects, err := bit.ListProjects()
	if err != nil {
		return err
	}

	for _, project := range projects {
		repos, err := bit.ListRepos(project.Key)
		if err != nil {
			logrus.Error(err)
			continue
		}

		for _, repo := range repos {
			if len(whitelist) > 0 {
				if !utils.InSlice(whitelist, repo.Name) {
					logrus.Debugf("skipping repo %s", repo.Slug)
					continue
				}

				logrus.Infof("cloning repo %s", repo.Name)
				_, err := git.Clone(repo.Links.Clone.GetSSH(), repo.RepoPath(b.Dir))
				if err != nil {
					logrus.Errorf("%s: %s", repo.Name, err.Error())
				}

				continue
			}

			if hasFile == "" {
				continue
			}

			content, err := bit.GetFileContent(project.Key, repo.Slug, hasFile, "")
			if err != nil {
				if strings.Contains(err.Error(), "does not exist at revision") {
					continue // file not found in repo
				}
				logrus.Errorf("%s: %s", repo.Name, err.Error())
				continue
			}

			if contentRegexp != "" && !contentReg.Match(content) {
				logrus.Infof("%s: file %s does not match regexp %s", repo.Name, hasFile, contentRegexp)
				continue
			}

			logrus.Infof("cloning repo %s", repo.Name)
			_, err = git.Clone(repo.Links.Clone.GetSSH(), repo.RepoPath(b.Dir))
			if err != nil {
				logrus.Errorf("%s: %s", repo.Name, err.Error())
			}
			// }
		}
	}

	return nil
}

// PR opens PR for each repo in repos list. If list is zero it opens for each repo in the --dir path.
func (b *Bitbucket) PR(repos []string, prMode string) error {
	var username, password string
	var err error

	open := func(bitbucketURL, gitDir string) error {
		return openBitbucketPRBrowser(bitbucketURL, gitDir)
	}

	if prMode == "api" {
		username, password, err = utils.Credentials()
		if err != nil {
			return err
		}
		bit := bitbucket.NewClient(b.BitbucketURL, username, password, false)
		open = func(bitbucketURL, gitDir string) error {
			return openBitbucketPRAPI(bit, gitDir)
		}
	}

	if len(repos) == 0 {
		return utils.InEachRepo(b.Dir, func(path string) error {
			logrus.Info(path)
			return open(b.BitbucketURL, path)
		})
	}

	for _, repo := range repos {
		dir := utils.RepoDir(b.WorkDir(), repo)
		err = open(b.BitbucketURL, dir)
		if err != nil {
			logrus.Errorf("error opening PR with api: %s", err)
		}
	}

	return nil
}

func openBitbucketPRAPI(bit *bitbucket.Client, gitDir string) error {
	origin, err := utils.Run("git", "-C", gitDir, "config", "--get", "remote.origin.url")
	if err != nil {
		return err
	}

	u, err := url.Parse(strings.TrimSpace(origin))
	if err != nil {
		return err
	}
	paths := strings.Split(u.Path, "/")
	projectId := strings.TrimSpace(paths[1])
	slug := strings.TrimSuffix(strings.TrimSpace(paths[2]), ".git")

	repo, err := bit.GetRepo(projectId, slug)
	if err != nil {
		return err
	}

	defaultBranch, err := bit.GetDefaultBranch(repo.Project.Key, repo.Slug)
	if err != nil {
		return err
	}
	sourceBranch, err := utils.Run("git", "-C", gitDir, "symbolic-ref", "HEAD")
	if err != nil {
		return err
	}

	defaultReviewers, err := bit.GetDefaultReviwers(repo.Project.Key, repo.Slug, repo.ID, strings.TrimSpace(sourceBranch), defaultBranch.ID)
	if err != nil {
		return err
	}

	var reviewers bitbucket.Reviewers

	for _, r := range defaultReviewers {
		reviewers = append(reviewers, bitbucket.Reviewer{
			User: bitbucket.User{
				Name: r.Name,
			},
		})
	}

	subject, err := utils.Run("git", "-C", gitDir, "show", "-s", "--format=%s")
	if err != nil {
		return err
	}
	body, err := utils.Run("git", "-C", gitDir, "show", "-s", "--format=%b")
	if err != nil {
		return err
	}

	pullRequest := bitbucket.PullRequest{
		Title:       subject,
		Description: body,
		FromRef: bitbucket.Ref{
			ID:         sourceBranch,
			Repository: repo,
		},
		ToRef: bitbucket.Ref{
			ID:         defaultBranch.ID,
			Repository: repo,
		},
		Reviewers: reviewers,
	}

	return bit.CreatePullRequest(projectId, slug, pullRequest)
}

func openBitbucketPRBrowser(bitbucketURL, gitDir string) error {
	targetBranch := "master"
	_, err := utils.Run("git", "-C", gitDir, "rev-parse", "-q", "--verify", "master")

	if err != nil && strings.Contains(err.Error(), "exit status 1") {
		targetBranch = "develop"
		_, err = utils.Run("git", "-C", gitDir, "rev-parse", "-q", "--verify", "develop")
		if err != nil && strings.Contains(err.Error(), "exit status 1") {
			return err
		}
	}

	origin, err := utils.Run("git", "-C", gitDir, "config", "--get", "remote.origin.url")
	if err != nil {
		return err
	}

	u, err := url.Parse(strings.TrimSpace(origin))
	if err != nil {
		return err
	}

	sourceBranch, err := utils.Run("git", "-C", gitDir, "symbolic-ref", "HEAD")
	if err != nil {
		return err
	}

	paths := strings.Split(u.Path, "/")
	burl := path.Join(bitbucketURL, "projects/%s/repos/%s/pull-requests?create&targetBranch=%s&sourceBranch=%s")
	prURL := fmt.Sprintf(burl, strings.TrimSpace(paths[1]), strings.TrimSuffix(strings.TrimSpace(paths[2]), ".git"), "refs/heads/"+targetBranch, strings.TrimSpace(sourceBranch))
	fmt.Println("Opening in browser:", prURL)
	return utils.OpenBrowser(prURL)
}
