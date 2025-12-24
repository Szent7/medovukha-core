package images

import (
	"io"

	"github.com/docker/docker/pkg/archive"
)

type ITarArchiver interface {
	TarWithOptions(srcPath string, options *archive.TarOptions) (io.ReadCloser, error)
}

type TarArchiver struct{}

func (r *TarArchiver) TarWithOptions(srcPath string, options *archive.TarOptions) (io.ReadCloser, error) {
	return archive.TarWithOptions(srcPath, options)
}
