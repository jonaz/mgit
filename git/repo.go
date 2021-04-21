package git

import (
	"fmt"
	"strings"

	"github.com/jonaz/mgit/utils"
)

type Repo struct {
	workdir string
}

func NewRepo(dir string) Repo {
	return Repo{workdir: dir}
}

func (repo Repo) WorkDir() string {
	return repo.workdir
}

func Clone(repo, dstFolder string) (Repo, error) {
	out, err := utils.Run("git", "clone", "--depth", "1", repo, dstFolder)
	if out != "" {
		fmt.Println(out)
	}
	return Repo{workdir: dstFolder}, err
}

func (repo Repo) RemoteURL() (string, error) {
	origin, err := utils.Run("git", "-C", repo.workdir, "config", "--get", "remote.origin.url")
	return strings.TrimSpace(origin), err
}

func (repo Repo) Add(file string) error {
	out, err := utils.Run("git", "-C", repo.workdir, "add", file)
	if out != "" {
		fmt.Println(out)
	}
	return err
}

func (repo Repo) Commit(msg, author string) error {
	var out string
	var err error

	if author != "" {
		out, err = utils.Run("git", "-C", repo.workdir, "commit", "--no-verify", "-a", "--author", author, "-m", msg)
	} else {
		out, err = utils.Run("git", "-C", repo.workdir, "commit", "--no-verify", "-a", "-m", msg)
	}
	if strings.TrimSpace(out) != "" {
		fmt.Println(out)
	}
	return err
}

func (repo Repo) Pull() error {
	out, err := utils.Run("git", "-C", repo.workdir, "pull", "--ff-only")
	if strings.TrimSpace(out) != "" {
		fmt.Println(out)
	}
	return err
}

func (repo Repo) Push(upstreamURL string, force bool) error {
	var err error
	if force {
		_, err = utils.Run("git", "-C", repo.workdir, "push", "--force", upstreamURL)
	} else {
		_, err = utils.Run("git", "-C", repo.workdir, "push", upstreamURL)
	}
	return err
}

func (repo Repo) Checkout(branch string) error {
	_, err := utils.Run("git", "-C", repo.workdir, "checkout", "-B", branch)
	return err
}

func (repo Repo) CurrentBranch() (string, error) {
	out, err := utils.Run("git", "-C", repo.workdir, "symbolic-ref", "--short", "HEAD")
	return strings.TrimSpace(out), err
}

/*

func (repo Repo) TrackPush(upstreamURL string) error {
	//git push --set-upstream origin `git symbolic-ref --short HEAD`

	branchName, err := repo.CurrentBranch()
	if err != nil {
		return err
	}

	_, err = utils.Run("git", "-C", repo.workdir, "push", "--set-upstream", "origin", branchName, upstreamURL)
	return err
}
*/
