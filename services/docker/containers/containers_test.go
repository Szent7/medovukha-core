package containers

import (
	"context"
	"errors"
	"io"
	ts "medovukha/api/rest/v1/types"
	"os"
	"strings"
	"testing"

	dc "medovukha/services/docker"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/assert"
)

func TestCreateTestContainer(t *testing.T) {
	mockClient := new(dc.MockDockerClient)

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"80/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "9990",
				},
			},
		},
	}
	var nc *network.NetworkingConfig = nil
	var pl *v1.Platform = nil
	imageName := "docker/welcome-to-docker"

	// test: test container found and created
	mockClient.On("ImagePull", context.Background(), imageName, image.PullOptions{}).Return(io.NopCloser(strings.NewReader("creating web-test")), nil).Once()
	mockClient.On("ContainerCreate", context.Background(), &container.Config{
		Image: imageName,
		Tty:   false,
	}, hostConfig, nc, pl, "web-test").Return(container.CreateResponse{
		ID: "testID",
	}, nil).Once()
	mockClient.On("ContainerStart", context.Background(), "testID", container.StartOptions{}).Return(nil).Once()

	err := CreateTestContainer(mockClient)
	assert.Nil(t, err)
	mockClient.AssertExpectations(t)

	// test: PullImage throw error
	mockClient.On("ImagePull", context.Background(), imageName, image.PullOptions{}).Return(io.NopCloser(strings.NewReader("creating web-test")), errors.New("ImagePull error")).Once()

	err = CreateTestContainer(mockClient)
	assert.EqualError(t, err, "ImagePull error")
	mockClient.AssertExpectations(t)

	// test: ContainerCreate throw error
	mockClient.On("ImagePull", context.Background(), imageName, image.PullOptions{}).Return(io.NopCloser(strings.NewReader("creating web-test")), nil).Once()
	mockClient.On("ContainerCreate", context.Background(), &container.Config{
		Image: imageName,
		Tty:   false,
	}, hostConfig, nc, pl, "web-test").Return(container.CreateResponse{}, errors.New("ContainerCreate error")).Once()

	err = CreateTestContainer(mockClient)
	assert.EqualError(t, err, "ContainerCreate error")
	mockClient.AssertExpectations(t)

	// test: ContainerStart throw error
	mockClient.On("ImagePull", context.Background(), imageName, image.PullOptions{}).Return(io.NopCloser(strings.NewReader("creating web-test")), nil).Once()
	mockClient.On("ContainerCreate", context.Background(), &container.Config{
		Image: imageName,
		Tty:   false,
	}, hostConfig, nc, pl, "web-test").Return(container.CreateResponse{
		ID: "testID",
	}, nil).Once()
	mockClient.On("ContainerStart", context.Background(), "testID", container.StartOptions{}).Return(errors.New("ContainerStart error")).Once()

	err = CreateTestContainer(mockClient)
	assert.EqualError(t, err, "ContainerStart error")
	mockClient.AssertExpectations(t)
}

func TestGetContainerBaseInfoList(t *testing.T) {
	mockClient := new(dc.MockDockerClient)

	// test: container found
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{
		{ID: "1234567890ab"},
	}, nil).Once()

	result, err := GetContainerBaseInfoList(mockClient)
	assert.Equal(t, []ts.ContainerBaseInfo{{Id: "1234567890ab"}}, result)
	assert.Nil(t, err)
	mockClient.AssertExpectations(t)

	// test: container list empty
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{}, nil).Once()

	result, err = GetContainerBaseInfoList(mockClient)
	assert.Equal(t, []ts.ContainerBaseInfo{}, result)
	assert.Nil(t, err)
	mockClient.AssertExpectations(t)

	// test: ContainerList throw error
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{}, errors.New("ContainerList error")).Once()

	result, err = GetContainerBaseInfoList(mockClient)
	assert.Nil(t, result)
	assert.EqualError(t, err, "ContainerList error")
	mockClient.AssertExpectations(t)
}

func TestPauseContainerByID(t *testing.T) {
	mockClient := new(dc.MockDockerClient)

	// test: container found and paused
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{
		{ID: "1234567890ab"},
	}, nil).Once()
	mockClient.On("ContainerPause", context.Background(), "1234567890ab").Return(nil).Once()

	err := PauseContainerByID(mockClient, "1234567890ab")
	assert.Nil(t, err)
	mockClient.AssertExpectations(t)

	// test: container not found
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{}, nil).Once()

	err = PauseContainerByID(mockClient, "1234567890ab")
	assert.EqualError(t, err, ts.ErrContainerNotFound.Error())
	mockClient.AssertExpectations(t)

	// test: ContainerList throw error
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{}, errors.New("ContainerList error")).Once()

	err = PauseContainerByID(mockClient, "1234567890ab")
	assert.EqualError(t, err, "ContainerList error")
	mockClient.AssertExpectations(t)

	// test: ContainerPause throw error
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{
		{ID: "1234567890ab"},
	}, nil).Once()
	mockClient.On("ContainerPause", context.Background(), "1234567890ab").Return(errors.New("ContainerPause error")).Once()

	err = PauseContainerByID(mockClient, "1234567890ab")
	assert.EqualError(t, err, "ContainerPause error")
	mockClient.AssertExpectations(t)
}

func TestUnpauseContainerByID(t *testing.T) {
	mockClient := new(dc.MockDockerClient)

	// test: container found and paused
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{
		{ID: "1234567890ab"},
	}, nil).Once()
	mockClient.On("ContainerUnpause", context.Background(), "1234567890ab").Return(nil).Once()

	err := UnpauseContainerByID(mockClient, "1234567890ab")
	assert.Nil(t, err)
	mockClient.AssertExpectations(t)

	// test: container not found
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{}, nil).Once()

	err = UnpauseContainerByID(mockClient, "1234567890ab")
	assert.EqualError(t, err, ts.ErrContainerNotFound.Error())
	mockClient.AssertExpectations(t)

	// test: ContainerList throw error
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{}, errors.New("ContainerList error")).Once()

	err = UnpauseContainerByID(mockClient, "1234567890ab")
	assert.EqualError(t, err, "ContainerList error")
	mockClient.AssertExpectations(t)

	// test: ContainerUnpause throw error
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{
		{ID: "1234567890ab"},
	}, nil).Once()
	mockClient.On("ContainerUnpause", context.Background(), "1234567890ab").Return(errors.New("ContainerUnpause error")).Once()

	err = UnpauseContainerByID(mockClient, "1234567890ab")
	assert.EqualError(t, err, "ContainerUnpause error")
	mockClient.AssertExpectations(t)
}

func TestKillContainerByID(t *testing.T) {
	mockClient := new(dc.MockDockerClient)

	// test: container found and killed
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{
		{ID: "1234567890ab"},
	}, nil).Once()
	mockClient.On("ContainerKill", context.Background(), "1234567890ab", "").Return(nil).Once()

	err := KillContainerByID(mockClient, "1234567890ab")
	assert.Nil(t, err)
	mockClient.AssertExpectations(t)

	// test: container not found
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{}, nil).Once()

	err = KillContainerByID(mockClient, "1234567890ab")
	assert.EqualError(t, err, ts.ErrContainerNotFound.Error())
	mockClient.AssertExpectations(t)

	// test: ContainerList throw error
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{}, errors.New("ContainerList error")).Once()

	err = KillContainerByID(mockClient, "1234567890ab")
	assert.EqualError(t, err, "ContainerList error")
	mockClient.AssertExpectations(t)

	// test: ContainerKill throw error
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{
		{ID: "1234567890ab"},
	}, nil).Once()
	mockClient.On("ContainerKill", context.Background(), "1234567890ab", "").Return(errors.New("ContainerKill error")).Once()

	err = KillContainerByID(mockClient, "1234567890ab")
	assert.EqualError(t, err, "ContainerKill error")
	mockClient.AssertExpectations(t)
}

func TestStopContainerByID(t *testing.T) {
	mockClient := new(dc.MockDockerClient)

	// test: container found and stopped
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{
		{ID: "1234567890ab"},
	}, nil).Once()
	mockClient.On("ContainerStop", context.Background(), "1234567890ab", container.StopOptions{}).Return(nil).Once()

	err := StopContainerByID(mockClient, "1234567890ab")
	assert.Nil(t, err)
	mockClient.AssertExpectations(t)

	// test: container not found
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{}, nil).Once()

	err = StopContainerByID(mockClient, "1234567890ab")
	assert.EqualError(t, err, ts.ErrContainerNotFound.Error())
	mockClient.AssertExpectations(t)

	// test: ContainerList throw error
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{}, errors.New("ContainerList error")).Once()

	err = StopContainerByID(mockClient, "1234567890ab")
	assert.EqualError(t, err, "ContainerList error")
	mockClient.AssertExpectations(t)

	// test: ContainerStop throw error
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{
		{ID: "1234567890ab"},
	}, nil).Once()
	mockClient.On("ContainerStop", context.Background(), "1234567890ab", container.StopOptions{}).Return(errors.New("ContainerStop error")).Once()

	err = StopContainerByID(mockClient, "1234567890ab")
	assert.EqualError(t, err, "ContainerStop error")
	mockClient.AssertExpectations(t)
}

func TestStartContainerByID(t *testing.T) {
	mockClient := new(dc.MockDockerClient)

	// test: container found and started
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{
		{ID: "1234567890ab"},
	}, nil).Once()
	mockClient.On("ContainerStart", context.Background(), "1234567890ab", container.StartOptions{}).Return(nil).Once()

	err := StartContainerByID(mockClient, "1234567890ab")
	assert.Nil(t, err)
	mockClient.AssertExpectations(t)

	// test: container not found
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{}, nil).Once()

	err = StartContainerByID(mockClient, "1234567890ab")
	assert.EqualError(t, err, ts.ErrContainerNotFound.Error())
	mockClient.AssertExpectations(t)

	// test: ContainerList throw error
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{}, errors.New("ContainerList error")).Once()

	err = StartContainerByID(mockClient, "1234567890ab")
	assert.EqualError(t, err, "ContainerList error")
	mockClient.AssertExpectations(t)

	// test: ContainerKill throw error
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{
		{ID: "1234567890ab"},
	}, nil).Once()
	mockClient.On("ContainerStart", context.Background(), "1234567890ab", container.StartOptions{}).Return(errors.New("ContainerStart error")).Once()

	err = StartContainerByID(mockClient, "1234567890ab")
	assert.EqualError(t, err, "ContainerStart error")
	mockClient.AssertExpectations(t)
}

func TestRestartContainerByID(t *testing.T) {
	mockClient := new(dc.MockDockerClient)

	// test: container found and restarted
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{
		{ID: "1234567890ab"},
	}, nil).Once()
	mockClient.On("ContainerRestart", context.Background(), "1234567890ab", container.StopOptions{}).Return(nil).Once()

	err := RestartContainerByID(mockClient, "1234567890ab")
	assert.Nil(t, err)
	mockClient.AssertExpectations(t)

	// test: container not found
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{}, nil).Once()

	err = RestartContainerByID(mockClient, "1234567890ab")
	assert.EqualError(t, err, ts.ErrContainerNotFound.Error())
	mockClient.AssertExpectations(t)

	// test: ContainerList throw error
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{}, errors.New("ContainerList error")).Once()

	err = RestartContainerByID(mockClient, "1234567890ab")
	assert.EqualError(t, err, "ContainerList error")
	mockClient.AssertExpectations(t)

	// test: ContainerRestart throw error
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{
		{ID: "1234567890ab"},
	}, nil).Once()
	mockClient.On("ContainerRestart", context.Background(), "1234567890ab", container.StopOptions{}).Return(errors.New("ContainerRestart error")).Once()

	err = RestartContainerByID(mockClient, "1234567890ab")
	assert.EqualError(t, err, "ContainerRestart error")
	mockClient.AssertExpectations(t)
}

func TestRemoveContainerByID(t *testing.T) {
	mockClient := new(dc.MockDockerClient)

	// test: container found and started
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{
		{ID: "1234567890ab"},
	}, nil).Once()
	mockClient.On("ContainerRemove", context.Background(), "1234567890ab", container.RemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         false,
	}).Return(nil).Once()

	err := RemoveContainerByID(mockClient, "1234567890ab")
	assert.Nil(t, err)
	mockClient.AssertExpectations(t)

	// test: container not found
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{}, nil).Once()

	err = RemoveContainerByID(mockClient, "1234567890ab")
	assert.EqualError(t, err, ts.ErrContainerNotFound.Error())
	mockClient.AssertExpectations(t)

	// test: ContainerList throw error
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{}, errors.New("ContainerList error")).Once()

	err = RemoveContainerByID(mockClient, "1234567890ab")
	assert.EqualError(t, err, "ContainerList error")
	mockClient.AssertExpectations(t)

	// test: ContainerKill throw error
	mockClient.On("ContainerList", context.Background(), container.ListOptions{All: true}).Return([]types.Container{
		{ID: "1234567890ab"},
	}, nil).Once()
	mockClient.On("ContainerRemove", context.Background(), "1234567890ab", container.RemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         false,
	}).Return(errors.New("ContainerRemove error")).Once()

	err = RemoveContainerByID(mockClient, "1234567890ab")
	assert.EqualError(t, err, "ContainerRemove error")
	mockClient.AssertExpectations(t)
}

func TestCheckIsMedovukhaId(t *testing.T) {
	hn, _ := os.Hostname()

	result, err := CheckIsMedovukhaId(hn + "hgklj56j3o4i6oj5i7jo")
	assert.NoError(t, err)
	assert.Equal(t, true, result)

	result, err = CheckIsMedovukhaId("3gfj4u64u" + hn + "hfh55ey")
	assert.NoError(t, err)
	assert.Equal(t, false, result)

	result, err = CheckIsMedovukhaId("1234567890ab")
	assert.NoError(t, err)
	assert.Equal(t, false, result)
}
