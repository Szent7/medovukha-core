package images

import (
	"io"

	"github.com/docker/docker/pkg/archive"
	"github.com/stretchr/testify/mock"
)

type MockTarArchiver struct {
	mock.Mock
}

func (m *MockTarArchiver) TarWithOptions(srcPath string, options *archive.TarOptions) (io.ReadCloser, error) {
	args := m.Called(srcPath, options)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}
