package utils

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// GetGitHash - Retrieve git hash from a local repo
func GetGitHash(repo string) string {
	var r *git.Repository
	var ref *plumbing.Reference
	var err error

	if r, err = git.PlainOpen(repo); err != nil {
		return ""
	}
	if ref, err = r.Head(); err != nil {
		return ""
	}

	return ref.Hash().String()
}

// GitClone - clone url to local directory
func GitClone(repoURL, targetDir string) error {
	_, err := git.PlainClone(targetDir, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: nil,
	})
	return err
}
