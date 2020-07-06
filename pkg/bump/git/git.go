package git

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/util"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/sirupsen/logrus"
)

// Git
type Git struct {
	Branch     string         `mapstructure:"branch"`
	Email      string         `mapstructure:"email"`
	Log        *logrus.Logger `mapstructure:"log"`
	RemoteName string         `mapstructure:"remote_name"`
	URL        string         `mapstructure:"url"`
	WorkDir    string         `mapstructure:"workdir"`
	auth       transport.AuthMethod
	filePaths  []string
	files      []billy.File
	fs         billy.Filesystem
	repo       *git.Repository
	worktree   *git.Worktree
}

var commitMsg = "[ci skip] ci: edit values with the new image tag\n\n\nskip-checks: true"
var name = "K8s Values Updater"

// SetAuth ...
func (g *Git) SetAuth(token string) (err error) {
	// Set Basic Auth
	g.auth = &http.BasicAuth{
		Username: "x-access-token",
		Password: token,
	}
	return nil
}

// Init ...
func (g *Git) Init(token string, fileNames string, dirPath string, separator string) (err error) {
	// Start In memory storage
	g.fs = memfs.New()

	// Set Basic Auth
	g.auth = &http.BasicAuth{
		Username: "x-access-token",
		Password: token,
	}

	g.Log.Info("Cloning")
	repo, err := git.Clone(memory.NewStorage(), g.fs, &git.CloneOptions{
		URL:           g.URL,
		Auth:          g.auth,
		SingleBranch:  true,
		ReferenceName: plumbing.NewBranchReferenceName(g.Branch),
	})
	if err != nil {
		return err
	}

	// Get WorkTree
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	// Store Repo and Worktree
	g.repo = repo
	g.worktree = worktree

	// Get The Files
	err = g.GetFiles(dirPath, fileNames, separator)
	if err != nil {
		return err
	}

	return nil
}

// Sync will force pull any changes
func (g *Git) Sync() error {
	// Pull
	g.Log.Info("pull")
	err := g.worktree.Pull(&git.PullOptions{
		RemoteName: g.RemoteName,
		Force:      true,
		Auth:       g.auth,
	})
	if err != nil {
		g.Log.Warn(err)
	}

	return nil
}

// Push will send any changes
func (g *Git) Push(fs billy.Filesystem) error {
	for _, fileName := range g.filePaths {
		file, err := fs.Open(fileName)
		if err != nil {
			return err
		}

		b, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}

		err = util.WriteFile(g.fs, fileName, []byte(b), os.ModePerm)
		if err != nil {
			return err
		}

		// Add
		_, err = g.worktree.Add(fileName)
		if err != nil {
			return err
		}
	}

	// Status
	status, err := g.worktree.Status()
	if err != nil {
		return err
	}
	g.Log.Info(status)

	// Create Commit
	commit, err := g.worktree.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  name,
			Email: g.Email,
			When:  time.Now(),
		},
	})

	// Commit
	_, err = g.repo.CommitObject(commit)
	if err != nil {
		return err
	}

	// Push
	g.Log.Info("push")
	err = g.repo.Push(&git.PushOptions{
		RemoteName: g.RemoteName,
		Auth:       g.auth,
	})
	if err != nil {
		return err
	}

	g.Log.Info("push success")
	return nil
}

// Files ...
func (g *Git) Files() []billy.File {
	return g.files
}

// GetFiles ...
func (g *Git) GetFiles(dir string, names string, separator string) error {
	filePaths := []string{}

	// Get the name of the files we will be working
	for _, fileName := range strings.Split(names, separator) {
		filePaths = append(g.filePaths, filepath.Join(dir, fileName))
	}

	// Loop Trough file names and get the actual files
	for _, filePath := range filePaths {
		matches, err := util.Glob(g.fs, filePath)
		if err != nil {
			return err
		}
		for _, match := range matches {
			g.filePaths = append(g.filePaths, match)
			f, err := g.fs.Open(match)
			if err != nil {
				return err
			}
			g.files = append(g.files, f)
		}
	}

	return nil
}
