package git

import (
	"os"

	"github.com/go-git/go-git/v5"
)

func CloneRepo(repocloner IRepoCloner, URI string, path string) error {

	_, err := repocloner.PlainClone(path, false, &git.CloneOptions{
		URL:      URI,
		Progress: os.Stdout,
	})

	if err != nil {
		return err
	}

	return nil
}
