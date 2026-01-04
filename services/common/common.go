package common

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/Szent7/medovukha-core/ipc/types"
	dockerpb "github.com/Szent7/medovukha-core/proto/docker/v1"
)

func SendLog(ch chan string, line string) {
	select {
	case ch <- line:
		{
			return
		}
	default:
		{
			<-ch
			ch <- line
		}
	}
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}

	fmt.Println("ошибка при проверке:", err)
	return false
}

func CreateFile(path string, content []byte) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(content); err != nil {
		return err
	}

	return nil
}

func RepoFromURL(raw string) (repo string, err error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid GitHub URL: %s", raw)
	}
	return parts[0] + "/" + parts[1], nil
}

func ConvertPorts(jsonPorts []types.Port) []*dockerpb.Port {
	if len(jsonPorts) == 0 {
		return nil
	}

	grpcPorts := make([]*dockerpb.Port, len(jsonPorts))
	for i := range jsonPorts {
		grpcPorts[i] = &dockerpb.Port{
			Ip:          jsonPorts[i].IP,
			PrivatePort: uint32(jsonPorts[i].PrivatePort),
			PublicPort:  uint32(jsonPorts[i].PublicPort),
			Type:        jsonPorts[i].Type,
		}
	}
	return grpcPorts
}
