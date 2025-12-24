package git

import (
	git "github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/mock"
)

type MockRepoCloner struct {
	mock.Mock
}

func (m *MockRepoCloner) PlainClone(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
	args := m.Called(path, isBare, o)
	return args.Get(0).(*git.Repository), args.Error(1)
}
