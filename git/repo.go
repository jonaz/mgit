package git

import (
	"fmt"
	"strings"

	"github.com/jonaz/mgit/utils"
	"github.com/sirupsen/logrus"
)

type Repo struct {
	workdir string
}

func (repo Repo) WorkDir() string {
	return repo.workdir
}

func Clone(repo, dstFolder string) (Repo, error) {
	// GIT_SSH_COMMAND='ssh -i /Users/UR_USERNAME/.ssh/UR_PRIVATE_KEY'
	logrus.Infof("clone git repo %s into %s ", repo, dstFolder)
	out, err := utils.Run("git", "clone", "--depth", "1", repo, dstFolder)
	if out != "" {
		fmt.Println(out)
	}
	return Repo{workdir: dstFolder}, err
}

func (repo Repo) Add(file string) error {
	out, err := utils.Run("git", "-C", repo.workdir, "add", file)
	if out != "" {
		fmt.Println(out)
	}
	return err
}

func (repo Repo) HasLocalChanges() (bool, error) {
	out, err := utils.Run("git", "-C", repo.workdir, "status", "--porcelain")
	return out != "", err
}

func (repo Repo) CommitAndPushIfLocalChanges(msg, author, upstreamURL string) error {
	localChanges, err := repo.HasLocalChanges()
	if err != nil {
		return err
	}
	if !localChanges {
		logrus.Debug("skipping commit and push since repo is already up to date")
		return nil
	}
	return repo.CommitAndPush(msg, author, upstreamURL)
}

func (repo Repo) CommitAndPush(msg, author, upstreamURL string) error {
	logrus.Infof("commit and push in %s to %s", repo.workdir, upstreamURL)
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
	if err != nil {
		return err
	}

	err = repo.Push(upstreamURL)

	if err != nil && strings.Contains(err.Error(), "Updates were rejected because the remote contains work that you do") {
		logrus.Error(err)
		logrus.Info(`push failed because "Updates were rejected because the remote contains work that you do". We will do a pull and then try the push again. `)
		err = repo.Pull()
		if err != nil {
			return err
		}
		return repo.Push(upstreamURL)
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

func (repo Repo) Push(upstreamURL string) error {
	out, err := utils.Run("git", "-C", repo.workdir, "push", upstreamURL)
	if strings.TrimSpace(out) != "" {
		fmt.Println(out)
	}
	return err
}
