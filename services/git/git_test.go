package git

import (
	"errors"
	"os"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
)

func TestCloneRepo(t *testing.T) {
	mockClient := new(MockRepoCloner)
	URI := "https://github.com/Szent7/medovukha"
	path := "/tmp"

	// test: repo cloned
	mockClient.On("PlainClone", "/tmp", false, &git.CloneOptions{
		URL:      URI,
		Progress: os.Stdout,
	}).Return(new(git.Repository), nil).Once()

	err := CloneRepo(mockClient, URI, path)
	assert.Nil(t, err)
	mockClient.AssertExpectations(t)

	// test: PlainClone throw error
	mockClient.On("PlainClone", "/tmp", false, &git.CloneOptions{
		URL:      URI,
		Progress: os.Stdout,
	}).Return(new(git.Repository), errors.New("PlainClone error")).Once()

	err = CloneRepo(mockClient, URI, path)
	assert.EqualError(t, err, "PlainClone error")
	mockClient.AssertExpectations(t)
}
