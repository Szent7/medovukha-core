package ipc

import (
	"context"

	"github.com/docker/docker/client"

	"errors"
	"fmt"
	"medovukha/ipc/types"

	"net/url"
	"os"
	"path/filepath"
	"strings"

	containers "medovukha/services/docker/containers"
	images "medovukha/services/docker/images"
	networks "medovukha/services/docker/networks"
	volumes "medovukha/services/docker/volumes"
	git "medovukha/services/git"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/filters"
)

type DockerCore struct {
	cli *client.Client
}

func NewDockerCore() (*DockerCore, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}
	return &DockerCore{cli: cli}, nil
}

func (d *DockerCore) Close() error {
	return d.cli.Close()
}

// Containers
func (d *DockerCore) GetContainerList(args *types.Empty, reply *[]types.ContainerBaseInfo) error {
	conList, err := containers.GetContainerBaseInfoList(d.cli)
	if err != nil {
		fmt.Printf("GetContainerList error: %s\n", err.Error())
		return err
	}

	*reply = conList
	return nil
}

func (d *DockerCore) PauseContainerByID(args *types.BaseID, reply *types.BaseMessage) error {
	if err := containers.PauseContainerByID(d.cli, args.ID); err != nil {
		fmt.Printf("PauseContainerByID error: %s\n", err.Error())
		return fmt.Errorf("PauseContainerByID error: %s", err.Error())
	}

	*reply = types.BaseMessage{Message: "Paused: " + args.ID}
	return nil
}

func (d *DockerCore) UnpauseContainerByID(args *types.BaseID, reply *types.BaseMessage) error {
	if err := containers.UnpauseContainerByID(d.cli, args.ID); err != nil {
		fmt.Printf("UnpauseContainerByID error: %s\n", err.Error())
		return fmt.Errorf("UnpauseContainerByID error: %s", err.Error())
	}

	*reply = types.BaseMessage{Message: "Unpaused: " + args.ID}
	return nil
}

func (d *DockerCore) KillContainerByID(args *types.BaseID, reply *types.BaseMessage) error {
	if err := containers.KillContainerByID(d.cli, args.ID); err != nil {
		fmt.Printf("KillContainerByID error: %s\n", err.Error())
		return fmt.Errorf("KillContainerByID error: %s", err.Error())
	}

	*reply = types.BaseMessage{Message: "Killed: " + args.ID}
	return nil
}

func (d *DockerCore) StartContainerByID(args *types.BaseID, reply *types.BaseMessage) error {
	if err := containers.StartContainerByID(d.cli, args.ID); err != nil {
		fmt.Printf("StartContainerByID error: %s\n", err.Error())
		return fmt.Errorf("StartContainerByID error: %s", err.Error())
	}

	*reply = types.BaseMessage{Message: "Started: " + args.ID}
	return nil
}

func (d *DockerCore) StopContainerByID(args *types.BaseID, reply *types.BaseMessage) error {
	if err := containers.StopContainerByID(d.cli, args.ID); err != nil {
		fmt.Printf("StopContainerByID error: %s\n", err.Error())
		return fmt.Errorf("StopContainerByID error: %s", err.Error())
	}

	*reply = types.BaseMessage{Message: "Stopped: " + args.ID}
	return nil
}

func (d *DockerCore) RestartContainerByID(args *types.BaseID, reply *types.BaseMessage) error {
	if err := containers.RestartContainerByID(d.cli, args.ID); err != nil {
		fmt.Printf("RestartContainerByID error: %s\n", err.Error())
		return fmt.Errorf("RestartContainerByID error: %s", err.Error())
	}

	*reply = types.BaseMessage{Message: "Restarted: " + args.ID}
	return nil
}

func (d *DockerCore) RemoveContainerByID(args *types.BaseID, reply *types.BaseMessage) error {
	if err := containers.RemoveContainerByID(d.cli, args.ID); err != nil {
		fmt.Printf("RemoveContainerByID error: %s\n", err.Error())
		return fmt.Errorf("RemoveContainerByID error: %s", err.Error())
	}

	*reply = types.BaseMessage{Message: "Removed: " + args.ID}
	return nil
}

// Images
func (d *DockerCore) GetImageList(args *types.Empty, reply *[]types.ImageBaseInfo) error {
	imgList, err := images.GetImageList(d.cli)
	if err != nil {
		fmt.Printf("GetImageList error: %s\n", err.Error())
		return fmt.Errorf("GetImageList error: %s", err.Error())
	}

	*reply = imgList
	return nil
}

// Networks
func (d *DockerCore) GetNetworkList(args *types.Empty, reply *[]types.NetworkBaseInfo) error {
	networkList, err := networks.GetNetworkList(d.cli)
	if err != nil {
		fmt.Printf("GetNetworkList error: %s\n", err.Error())
		return fmt.Errorf("GetNetworkList error: %s", err.Error())
	}

	*reply = networkList
	return nil
}

// Volumes
func (d *DockerCore) GetVolumeList(args *types.Empty, reply *[]types.VolumeBaseInfo) error {
	volumeList, err := volumes.GetVolumeList(d.cli)
	if err != nil {
		fmt.Printf("GetVolumeList error: %s\n", err.Error())
		return fmt.Errorf("GetVolumeList error: %s", err.Error())
	}

	*reply = volumeList
	return nil
}

// Deploy
func (d *DockerCore) CreateFromGit(args *types.DeployFromGit, reply *types.BaseMessage) error {
	ctx := context.Background()
	//
	//
	// Clonning repo from git into temp dir
	//
	//
	tempDir, err := os.MkdirTemp("", "docker-repo")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Created tempDir: " + tempDir)
	defer os.RemoveAll(tempDir)
	defer fmt.Println("Deleted tempDir: " + tempDir)

	if err := git.CloneRepo(&git.RepoCloner{}, args.URL, tempDir); err != nil {
		fmt.Printf("CloneRepo error: %s\n", err.Error())
		return fmt.Errorf("CreateFromGit error: %s", err.Error())
	}
	//
	//
	// Check Dockerfile in request
	//
	//
	dockerfileDir := filepath.Join(tempDir, "Dockerfile")
	dockercomposeDir := filepath.Join(tempDir, "docker-compose.yml")
	if args.Dockerfile == "" {
		// Dockerfile is not specified in the request, check in the directory
		// If Dockerfile doesn`t exist, throw error
		if !fileExists(dockerfileDir) {
			fmt.Println("Dockerfile empty error")
			return fmt.Errorf("CreateFromGit error: Dockerfile empty")
		}
	} else {
		// If Dockerfile specified in the request, rewrite/create new Dockerfile
		if err := createFile(dockerfileDir, []byte(args.Dockerfile)); err != nil {
			fmt.Printf("Dockerfile write error: %s\n", err.Error())
			return fmt.Errorf("CreateFromGit error: %s", err.Error())
		} else {
			fmt.Println("Dockerfile redefined")
		}
	}
	//
	//
	// Parse tags from repo (for image name)
	//
	//
	repo, err := repoFromURL(args.URL)
	if err != nil {
		fmt.Printf("Repo error: %s\n", err.Error())
		return fmt.Errorf("CreateFromGit error: %s", err.Error())
	}
	tags := []string{repo + ":latest"} //args.URL
	//
	//
	// Delete old containers
	//
	//
	if err := containers.RemoveContainerByImage(ctx, d.cli, tags[0]); err != nil {
		fmt.Printf("Error remove container: %s\n", err.Error())
		return fmt.Errorf("CreateFromGit error: %s", err.Error())
	}
	//
	//
	// Delete old images
	//
	//
	if err := images.RemoveImageByTag(ctx, d.cli, tags[0]); err != nil {
		fmt.Printf("Error remove image: %s\n", err.Error())
		return fmt.Errorf("CreateFromGit error: %s", err.Error())
	}
	//
	//
	// Build image
	//
	//
	newImageId, err := images.BuildImageNew(d.cli, tempDir, tags)
	if err != nil {
		fmt.Printf("Docker client error: %s\n", err.Error())
		return fmt.Errorf("CreateFromGit error: %s", err.Error())
	}
	//
	//
	// Clear build cache
	//
	//
	pruneFilters := filters.NewArgs()
	pruneFilters.Add("dangling", "true")

	_, err = d.cli.ImagesPrune(ctx, pruneFilters)
	if err != nil {
		fmt.Printf("Docker client error: %s\n", err.Error())
		return fmt.Errorf("CreateFromGit error: %s", err.Error())
	}
	_, err = d.cli.BuildCachePrune(ctx, build.CachePruneOptions{All: true})
	if err != nil {
		fmt.Printf("Docker client error: %s\n", err.Error())
		return fmt.Errorf("CreateFromGit error: %s", err.Error())
	}
	//
	//
	// Check DockerCompose in request
	//
	//
	if args.DockerCompose == "" {
		// DockerCompose is not specified in the request, check DockerRun in the request
		if args.DockerRun == "" {
			// DockerRun is not specified in the request, check DockerCompose in the directory
			if !fileExists(dockercomposeDir) {
				fmt.Println("Created image, but not command to launch container")
				*reply = types.BaseMessage{Message: "created but not launched: " + newImageId}
				return nil
			} else {
				if err := containers.ExecDockerComposeUp(dockercomposeDir); err != nil {
					*reply = types.BaseMessage{Message: "the image was created, but an error occurred when starting the container (dockerCompose): " + newImageId}
					return nil
				}
			}
		} else {
			if err := containers.ExecDockerRun(args.DockerRun); err != nil {
				*reply = types.BaseMessage{Message: "the image was created, but an error occurred when starting the container (dockerRun): " + newImageId}
				return nil
			}
		}
	} else {
		// If Dockerfile specified in the request, rewrite/create new Dockerfile
		if err := createFile(dockercomposeDir, []byte(args.DockerCompose)); err != nil {
			fmt.Printf("DockerCompose write error: %s\n", err.Error())
			return fmt.Errorf("CreateFromGit error: %s", err.Error())
		} else {
			fmt.Println("DockerCompose redefined")
			if err := containers.ExecDockerComposeUp(dockercomposeDir); err != nil {
				*reply = types.BaseMessage{Message: "the image was created, but an error occurred when starting the container (dockerCompose): " + newImageId}
				return nil
			}
		}
	}

	*reply = types.BaseMessage{Message: "created and launched: " + newImageId}
	return nil
}

func fileExists(path string) bool {
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

func createFile(path string, content []byte) error {
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

func repoFromURL(raw string) (repo string, err error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	// u.Path выглядит как "/owner/repo" (или "/owner/repo/…")
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid GitHub URL: %s", raw)
	}
	return parts[0] + "/" + parts[1], nil
}
