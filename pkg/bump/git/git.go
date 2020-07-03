package git

import (
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/sirupsen/logrus"
)

// Git
type Git struct {
	WorkDir    string
	RemoteName string
	Email      string
	auth       transport.AuthMethod
	Log        *logrus.Logger
}

// SetBasicAuth Sets the chosen auth
func (g *Git) SetBasicAuth(pass string) error {
	g.auth = &http.TokenAuth{
		Token: pass,
	}

	return nil
}

// Sync will force pull any changes
func (g *Git) Sync() error {
	// Opens an already existing repository.
	r, err := git.PlainOpen(g.WorkDir)
	if err != nil {
		return err
	}

	// Get WorkTree
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	// Pull
	g.Log.Info("pull")
	err = w.Pull(&git.PullOptions{
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
func (g *Git) Push(files []string) error {
	//
	commitMsg := "[ci skip] ci: edit values with the new image tag\n\n\nskip-checks: true"
	name := "K8s Values Updater"

	// Opens an already existing repository.
	r, err := git.PlainOpen(g.WorkDir)
	if err != nil {
		return err
	}

	// Get WorkTree
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	// Add
	for _, f := range files {
		_, err = w.Add(f)
		if err != nil {
			return err
		}
	}

	// Status
	status, err := w.Status()
	if err != nil {
		return err
	}
	g.Log.Debug(status)

	// Create Commit
	commit, err := w.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  name,
			Email: g.Email,
			When:  time.Now(),
		},
	})

	// Commit
	obj, err := r.CommitObject(commit)
	if err != nil {
		return err
	}
	g.Log.Debug(obj)

	// Push
	g.Log.Info("push")
	err = r.Push(&git.PushOptions{
		RemoteName: g.RemoteName,
		Auth:       g.auth,
	})
	if err != nil {
		return err
	}

	//
	g.Log.Info("ok")
	return nil
}
