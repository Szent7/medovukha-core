package volumes

import (
	"context"
	"errors"
	ts "medovukha/api/rest/v1/types"
	"testing"

	dc "medovukha/services/docker"

	"github.com/docker/docker/api/types/volume"
	"github.com/stretchr/testify/assert"
)

func TestGetVolumeList(t *testing.T) {
	mockClient := new(dc.MockDockerClient)

	// test: volume found
	mockClient.On("VolumeList", context.Background(), volume.ListOptions{}).Return(volume.ListResponse{
		Volumes: []*volume.Volume{{Driver: "test"}},
	}, nil).Once()

	result, err := GetVolumeList(mockClient)
	assert.Equal(t, []ts.VolumeBaseInfo{
		{Driver: "test"},
	}, result)
	assert.Nil(t, err)
	mockClient.AssertExpectations(t)

	// test: network list empty
	mockClient.On("VolumeList", context.Background(), volume.ListOptions{}).Return(volume.ListResponse{}, nil).Once()

	result, err = GetVolumeList(mockClient)
	assert.Equal(t, []ts.VolumeBaseInfo{}, result)
	assert.Nil(t, err)
	mockClient.AssertExpectations(t)

	// test: VolumeList throw error
	mockClient.On("VolumeList", context.Background(), volume.ListOptions{}).Return(volume.ListResponse{}, errors.New("VolumeList error")).Once()

	result, err = GetVolumeList(mockClient)
	assert.Nil(t, result)
	assert.EqualError(t, err, "VolumeList error")
	mockClient.AssertExpectations(t)
}
