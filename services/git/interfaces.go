package git

import (
	git "github.com/go-git/go-git/v5"
)

type IRepoCloner interface {
	PlainClone(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error)
}

type RepoCloner struct{}

func (r *RepoCloner) PlainClone(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
	return git.PlainClone(path, isBare, o)
}
