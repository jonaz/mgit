package providers

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/jonaz/mgit/config"
	"github.com/jonaz/mgit/git"
	"github.com/jonaz/mgit/service"
	"github.com/jonaz/mgit/utils"
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

func (b *Bitbucket) Clone(whitelist []string, hasFile string) error {
	username, password, err := utils.Credentials()
	if err != nil {
		return err
	}

	config := config.Config{
		BitbucketURL:      b.BitbucketURL,
		BitbucketUser:     username,
		BitbucketPassword: password,
	}
	bit := service.NewBitbucket(config)

	projects, err := bit.ListProjects()
	if err != nil {
		return err
	}

	for _, project := range projects.Values {
		repos, err := bit.ListRepos(project.Key)
		if err != nil {
			logrus.Error(err)
			continue
		}

		for _, repo := range repos.Values {
			if len(whitelist) > 0 {
				if !utils.InSlice(whitelist, repo.Name) {
					logrus.Debugf("skipping repo %s", repo.Slug)
					continue
				}

				logrus.Infof("cloning repo %s", repo.Name)
				_, err := git.Clone(repo.Links.Clone.GetSSH(), repo.RepoPath(b.Dir))
				if err != nil {
					logrus.Error(err)
				}

				continue
			}

			if hasFile != "" {
				files, err := bit.ListFiles(project.Key, repo.Slug)
				if err != nil {
					logrus.Error(err)
					continue
				}
				if utils.InSlice(files.Values, hasFile) {
					logrus.Infof("cloning repo %s", repo.Name)
					_, err := git.Clone(repo.Links.Clone.GetSSH(), repo.RepoPath(b.Dir))
					if err != nil {
						logrus.Error(err)
					}
				}
				continue
			}
		}
	}

	return nil
}

//PR opens PR for each repo in repos list. If list is zero it opens for each repo in the --dir path.
func (b *Bitbucket) PR(repos []string) error {
	if len(repos) == 0 {
		return utils.InEachRepo(b.Dir, func(path string) error {
			logrus.Info(path)
			return openBitbucketPR(b.BitbucketURL, "", path)
		})
	}

	var err error
	for _, repo := range repos {
		dir := utils.RepoDir(b.WorkDir(), repo)
		err = openBitbucketPR(b.BitbucketURL, "", dir)
		if err != nil {
			return err
		}
	}

	return nil
}

func openBitbucketPR(bitbucketURL, targetBranch, gitDir string) error {
	if targetBranch == "" {
		targetBranch = "master"
		_, err := utils.Run("git", "-C", gitDir, "rev-parse", "-q", "--verify", "master")
		if err != nil && err.Error() == "exit status 1" {
			targetBranch = "develop"
			_, err = utils.Run("git", "-C", gitDir, "rev-parse", "-q", "--verify", "develop")
			if err != nil && err.Error() == "exit status 1" {
				return err
			}
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
