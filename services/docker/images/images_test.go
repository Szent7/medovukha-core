package images

import (
	"context"
	"errors"
	"io"
	ts "medovukha/api/rest/v1/types"
	"strings"
	"testing"

	dc "medovukha/services/docker"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/pkg/archive"
	"github.com/stretchr/testify/assert"
)

func TestPullImage(t *testing.T) {
	mockClient := new(dc.MockDockerClient)

	//test: image found and pulled
	mockClient.On("ImagePull", context.Background(), "docker/welcome-to-docker", image.PullOptions{}).Return(io.NopCloser(strings.NewReader("creating web-test")), nil).Once()

	err := PullImage(mockClient, context.Background(), "docker/welcome-to-docker")
	assert.Nil(t, err)
	mockClient.AssertExpectations(t)

	//test: image not pulled
	mockClient.On("ImagePull", context.Background(), "docker/welcome-to-docker", image.PullOptions{}).Return(io.NopCloser(strings.NewReader("creating web-test")), errors.New("ImagePull error")).Once()

	err = PullImage(mockClient, context.Background(), "docker/welcome-to-docker")
	assert.EqualError(t, err, "ImagePull error")
	mockClient.AssertExpectations(t)
}

func TestGetImageList(t *testing.T) {
	mockClient := new(dc.MockDockerClient)

	// test: image found
	mockClient.On("ImageList", context.Background(), image.ListOptions{All: true, ContainerCount: true}).Return([]image.Summary{
		{ID: "1234567890ab"},
	}, nil).Once()

	result, err := GetImageList(mockClient)
	assert.Equal(t, []ts.ImageBaseInfo{{Id: "1234567890ab"}}, result)
	assert.Nil(t, err)
	mockClient.AssertExpectations(t)

	// test: image list empty
	mockClient.On("ImageList", context.Background(), image.ListOptions{All: true, ContainerCount: true}).Return([]image.Summary{}, nil).Once()

	result, err = GetImageList(mockClient)
	assert.Equal(t, []ts.ImageBaseInfo{}, result)
	assert.Nil(t, err)
	mockClient.AssertExpectations(t)

	// test: ImageList throw error
	mockClient.On("ImageList", context.Background(), image.ListOptions{All: true, ContainerCount: true}).Return([]image.Summary{}, errors.New("ImageList error")).Once()

	result, err = GetImageList(mockClient)
	assert.Nil(t, result)
	assert.EqualError(t, err, "ImageList error")
	mockClient.AssertExpectations(t)
}

func TestGetRawImageList(t *testing.T) {
	mockClient := new(dc.MockDockerClient)

	// test: image found
	mockClient.On("ImageList", context.Background(), image.ListOptions{All: true, ContainerCount: true}).Return([]image.Summary{
		{ID: "1234567890ab"},
	}, nil).Once()

	result, err := GetRawImageList(mockClient)
	assert.Equal(t, []image.Summary{
		{ID: "1234567890ab"}}, result)
	assert.Nil(t, err)
	mockClient.AssertExpectations(t)

	// test: image list empty
	mockClient.On("ImageList", context.Background(), image.ListOptions{All: true, ContainerCount: true}).Return([]image.Summary{}, nil).Once()

	result, err = GetRawImageList(mockClient)
	assert.Equal(t, []image.Summary{}, result)
	assert.Nil(t, err)
	mockClient.AssertExpectations(t)

	// test: ImageList throw error
	mockClient.On("ImageList", context.Background(), image.ListOptions{All: true, ContainerCount: true}).Return([]image.Summary{}, errors.New("ImageList error")).Once()

	result, err = GetRawImageList(mockClient)
	assert.Nil(t, result)
	assert.EqualError(t, err, "ImageList error")
	mockClient.AssertExpectations(t)
}

func TestGetImageListIDs(t *testing.T) {
	mockClient := new(dc.MockDockerClient)
	IDs := []string{"1234567890ab", "1234567890ac", "1234567890ad"}
	imageListRes := []image.Summary{
		{ID: "1234567890ab"},
		{ID: "1234567890ac"},
		{ID: "1234567890ad"},
	}

	// test: image found
	mockClient.On("ImageList", context.Background(), image.ListOptions{All: true, ContainerCount: true}).Return(imageListRes, nil).Once()

	result, err := GetImageListIDs(mockClient)
	assert.Equal(t, result, IDs)
	assert.Nil(t, err)
	mockClient.AssertExpectations(t)

	// test: image list empty
	mockClient.On("ImageList", context.Background(), image.ListOptions{All: true, ContainerCount: true}).Return([]image.Summary{}, nil).Once()

	result, err = GetImageListIDs(mockClient)
	assert.Equal(t, result, []string{})
	assert.Nil(t, err)
	mockClient.AssertExpectations(t)

	// test: ImageList throw error
	mockClient.On("ImageList", context.Background(), image.ListOptions{All: true, ContainerCount: true}).Return([]image.Summary{}, errors.New("ImageList error")).Once()

	result, err = GetImageListIDs(mockClient)
	assert.Nil(t, result)
	assert.EqualError(t, err, "ImageList error")
	mockClient.AssertExpectations(t)
}

func TestBuildImage(t *testing.T) {
	mockClientDocker := new(dc.MockDockerClient)
	mockClientTar := new(MockTarArchiver)
	tags := []string{"test/test:latest"}
	imageListRes := []image.Summary{
		{ID: "1234567890ab", RepoTags: []string{""}},
		{ID: "1234567890ac", RepoTags: []string{tags[0]}},
		{ID: "1234567890ad", RepoTags: []string{"test/test:v1"}},
	}

	// test: tags are empty
	result, err := BuildImage(mockClientDocker, mockClientTar, "/tmp", []string{})
	assert.Equal(t, result, "")
	assert.EqualError(t, err, ts.ErrEmptyTags.Error())

	// test: image created successfully
	mockClientTar.On("TarWithOptions", "/tmp", &archive.TarOptions{}).Return(io.NopCloser(strings.NewReader("creating tar archive")), nil).Once()
	mockClientDocker.On("ImageBuild", context.Background(), io.NopCloser(strings.NewReader("creating tar archive")), types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       tags,
		Remove:     true,
	}).Return(types.ImageBuildResponse{
		Body: io.NopCloser(strings.NewReader("building image")),
	}, nil).Once()
	mockClientDocker.On("ImageList", context.Background(), image.ListOptions{All: true, ContainerCount: true}).Return(imageListRes, nil).Once()

	result, err = BuildImage(mockClientDocker, mockClientTar, "/tmp", tags)
	assert.Equal(t, result, imageListRes[1].ID)
	assert.Nil(t, err)
	mockClientTar.AssertExpectations(t)
	mockClientDocker.AssertExpectations(t)

	// test: TarWithOptions throw error
	mockClientTar.On("TarWithOptions", "/tmp", &archive.TarOptions{}).Return(io.NopCloser(strings.NewReader("tar error")), errors.New("TarWithOptions error")).Once()

	result, err = BuildImage(mockClientDocker, mockClientTar, "/tmp", tags)
	assert.Equal(t, result, "")
	assert.EqualError(t, err, "TarWithOptions error")
	mockClientTar.AssertExpectations(t)

	// test: ImageBuild throw error
	mockClientTar.On("TarWithOptions", "/tmp", &archive.TarOptions{}).Return(io.NopCloser(strings.NewReader("creating tar archive")), nil).Once()
	mockClientDocker.On("ImageBuild", context.Background(), io.NopCloser(strings.NewReader("creating tar archive")), types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       tags,
		Remove:     true,
	}).Return(types.ImageBuildResponse{
		Body: io.NopCloser(strings.NewReader("building image error")),
	}, errors.New("ImageBuild error")).Once()

	result, err = BuildImage(mockClientDocker, mockClientTar, "/tmp", tags)
	assert.Equal(t, result, "")
	assert.EqualError(t, err, "ImageBuild error")
	mockClientTar.AssertExpectations(t)
	mockClientDocker.AssertExpectations(t)

	// test: ImageList throw error
	mockClientTar.On("TarWithOptions", "/tmp", &archive.TarOptions{}).Return(io.NopCloser(strings.NewReader("creating tar archive")), nil).Once()
	mockClientDocker.On("ImageBuild", context.Background(), io.NopCloser(strings.NewReader("creating tar archive")), types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       tags,
		Remove:     true,
	}).Return(types.ImageBuildResponse{
		Body: io.NopCloser(strings.NewReader("building image")),
	}, nil).Once()
	mockClientDocker.On("ImageList", context.Background(), image.ListOptions{All: true, ContainerCount: true}).Return([]image.Summary{}, errors.New("ImageList error")).Once()

	result, err = BuildImage(mockClientDocker, mockClientTar, "/tmp", tags)
	assert.Equal(t, result, "")
	assert.EqualError(t, err, "ImageList error")
	mockClientTar.AssertExpectations(t)
	mockClientDocker.AssertExpectations(t)
}

func TestGetImageIDFromTag(t *testing.T) {
	mockClient := new(dc.MockDockerClient)
	tag := "test/test:latest"
	imageListRes := []image.Summary{
		{ID: "1234567890ab", RepoTags: []string{""}},
		{ID: "1234567890ac", RepoTags: []string{tag}},
		{ID: "1234567890ad", RepoTags: []string{"test/test:v1"}},
	}

	// test: image found
	mockClient.On("ImageList", context.Background(), image.ListOptions{All: true, ContainerCount: true}).Return(imageListRes, nil).Once()

	result, err := GetImageIDFromTag(mockClient, tag)
	assert.Equal(t, result, imageListRes[1].ID)
	assert.Nil(t, err)
	mockClient.AssertExpectations(t)

	// test: image list empty
	mockClient.On("ImageList", context.Background(), image.ListOptions{All: true, ContainerCount: true}).Return([]image.Summary{}, nil).Once()

	result, err = GetImageIDFromTag(mockClient, tag)
	assert.Equal(t, result, "")
	assert.Nil(t, err)
	mockClient.AssertExpectations(t)

	// test: ImageList throw error
	mockClient.On("ImageList", context.Background(), image.ListOptions{All: true, ContainerCount: true}).Return([]image.Summary{}, errors.New("ImageList error")).Once()

	result, err = GetImageIDFromTag(mockClient, tag)
	assert.Equal(t, result, "")
	assert.EqualError(t, err, "ImageList error")
	mockClient.AssertExpectations(t)
}

func TestRemoveIntermediateImages(t *testing.T) {
	mockClient := new(dc.MockDockerClient)
	imagesIDs := []string{"1234567890ab", "1234567890ac", "1234567890ad"}
	// test: all images removed
	for _, id := range imagesIDs {
		mockClient.On("ImageRemove", context.Background(), id, image.RemoveOptions{Force: true}).Return([]image.DeleteResponse{
			{Deleted: id},
		}, nil).Once()
	}

	result := RemoveIntermediateImages(mockClient, imagesIDs)
	assert.Equal(t, make([]string, 0, 3), result)
	mockClient.AssertExpectations(t)

	// test: some images are not removed
	for i, id := range imagesIDs {
		if i >= 1 {
			mockClient.On("ImageRemove", context.Background(), id, image.RemoveOptions{Force: true}).Return([]image.DeleteResponse{},
				errors.New("ImageRemove error")).Once()
		} else {
			mockClient.On("ImageRemove", context.Background(), id, image.RemoveOptions{Force: true}).Return([]image.DeleteResponse{
				{Deleted: id},
			}, nil).Once()
		}
	}

	result = RemoveIntermediateImages(mockClient, imagesIDs)
	assert.Equal(t, imagesIDs[1:], result)
	mockClient.AssertExpectations(t)

	// test: all images are not removed
	for _, id := range imagesIDs {
		mockClient.On("ImageRemove", context.Background(), id, image.RemoveOptions{Force: true}).Return([]image.DeleteResponse{},
			errors.New("ImageRemove error")).Once()
	}

	result = RemoveIntermediateImages(mockClient, imagesIDs)
	assert.Equal(t, imagesIDs, result)
	mockClient.AssertExpectations(t)

	// test: image list empty
	result = RemoveIntermediateImages(mockClient, make([]string, 0))
	assert.Equal(t, make([]string, 0), result)
	mockClient.AssertExpectations(t)
}
