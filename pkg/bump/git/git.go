package git

import (
	"fmt"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// Git
type Git struct {
	WorkDir    string
	RemoteName string
	Email      string
	ForceSSH   bool
	auth       transport.AuthMethod
}

// SetBasicAuth Sets the chosen auth
func (g *Git) SetBasicAuth(user string, pass string) error {
	if pass != "" && !g.ForceSSH {
		fmt.Println("Will authenticate with basic auth")
		g.auth = &http.BasicAuth{
			Username: user,
			Password: pass,
		}

		return nil
	}
	fmt.Println("Will authenticate with SSH")
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
	fmt.Println("Pull")
	err = w.Pull(&git.PullOptions{
		RemoteName: g.RemoteName,
		Force:      true,
		Auth:       g.auth,
	})
	if err != nil {
		fmt.Println(err)
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
	fmt.Println(status)

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
	fmt.Println(obj)

	// Push
	fmt.Println("Push")
	err = r.Push(&git.PushOptions{
		RemoteName: g.RemoteName,
		Auth:       g.auth,
	})
	if err != nil {
		return err
	}

	//
	fmt.Println("Ok")
	return nil
}
