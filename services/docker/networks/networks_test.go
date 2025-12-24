package networks

import (
	"context"
	"errors"
	ts "medovukha/api/rest/v1/types"
	"testing"

	dc "medovukha/services/docker"

	"github.com/docker/docker/api/types/network"
	"github.com/stretchr/testify/assert"
)

func TestGetNetworkList(t *testing.T) {
	mockClient := new(dc.MockDockerClient)

	// test: network found
	mockClient.On("NetworkList", context.Background(), network.ListOptions{}).Return([]network.Summary{
		{ID: "1234567890ab"},
	}, nil).Once()

	result, err := GetNetworkList(mockClient)
	assert.Equal(t, []ts.NetworkBaseInfo{{Id: "1234567890ab", Subnet: []string{}, Gateway: []string{}}}, result)
	assert.Nil(t, err)
	mockClient.AssertExpectations(t)

	// test: network list empty
	mockClient.On("NetworkList", context.Background(), network.ListOptions{}).Return([]network.Summary{}, nil).Once()

	result, err = GetNetworkList(mockClient)
	assert.Equal(t, []ts.NetworkBaseInfo{}, result)
	assert.Nil(t, err)
	mockClient.AssertExpectations(t)

	// test: NetworkList throw error
	mockClient.On("NetworkList", context.Background(), network.ListOptions{}).Return([]network.Summary{}, errors.New("NetworkList error")).Once()

	result, err = GetNetworkList(mockClient)
	assert.Nil(t, result)
	assert.EqualError(t, err, "NetworkList error")
	mockClient.AssertExpectations(t)
}
