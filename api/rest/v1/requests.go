package v1

import (
	"context"
	"errors"
	"fmt"
	"medovukha/api/rest/v1/types"
	"medovukha/services/docker"

	"net/http"
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
	"github.com/gin-gonic/gin"
)

// Containers
func CreateTestContainer(c *gin.Context) {
	cli, err := docker.CreateDockerClient()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Docker client error"})
		fmt.Printf("Docker client error: %s\n", err.Error())
		return
	}
	defer cli.Close()

	if err := containers.CreateTestContainer(cli); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "CreateTestContainer error"})
		fmt.Printf("CreateTestContainer error: %s\n", err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Created TestContainer"})
}

func GetContainerList(c *gin.Context) {
	cli, err := docker.CreateDockerClient()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Docker client error"})
		fmt.Printf("Docker client error: %s\n", err.Error())
		return
	}
	defer cli.Close()

	conList, err := containers.GetContainerBaseInfoList(cli)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "GetContainerList error"})
		fmt.Printf("GetContainerList error: %s\n", err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, conList)
}

func PauseContainerByID(c *gin.Context) {
	var containerId struct {
		Id string `json:"id"`
	}
	if err := c.BindJSON(&containerId); err != nil {
		return
	}

	if check, err := containers.CheckIsMedovukhaId(containerId.Id); err != nil {
		return
	} else if check {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": types.ErrContainerIsMedovukha.Error()})
		return
	}

	cli, err := docker.CreateDockerClient()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Docker client error"})
		fmt.Printf("Docker client error: %s\n", err.Error())
		return
	}
	defer cli.Close()

	if err := containers.PauseContainerByID(cli, containerId.Id); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "PauseContainerByID error"})
		fmt.Printf("PauseContainerByID error: %s\n", err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Paused: " + containerId.Id})
}

func UnpauseContainerByID(c *gin.Context) {
	var containerId struct {
		Id string `json:"id"`
	}
	if err := c.BindJSON(&containerId); err != nil {
		return
	}

	if check, err := containers.CheckIsMedovukhaId(containerId.Id); err != nil {
		return
	} else if check {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": types.ErrContainerIsMedovukha.Error()})
		return
	}

	cli, err := docker.CreateDockerClient()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Docker client error"})
		fmt.Printf("Docker client error: %s\n", err.Error())
		return
	}
	defer cli.Close()

	if err := containers.UnpauseContainerByID(cli, containerId.Id); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "UnpauseContainerByID error"})
		fmt.Printf("UnpauseContainerByID error: %s\n", err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Unpaused: " + containerId.Id})
}

func KillContainerByID(c *gin.Context) {
	var containerId struct {
		Id string `json:"id"`
	}
	if err := c.BindJSON(&containerId); err != nil {
		return
	}

	if check, err := containers.CheckIsMedovukhaId(containerId.Id); err != nil {
		return
	} else if check {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": types.ErrContainerIsMedovukha.Error()})
		return
	}

	cli, err := docker.CreateDockerClient()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Docker client error"})
		fmt.Printf("Docker client error: %s\n", err.Error())
		return
	}
	defer cli.Close()

	if err := containers.KillContainerByID(cli, containerId.Id); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "KillContainerByID error"})
		fmt.Printf("KillContainerByID error: %s\n", err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Killed: " + containerId.Id})
}

func StartContainerByID(c *gin.Context) {
	var containerId struct {
		Id string `json:"id"`
	}
	if err := c.BindJSON(&containerId); err != nil {
		return
	}

	if check, err := containers.CheckIsMedovukhaId(containerId.Id); err != nil {
		return
	} else if check {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": types.ErrContainerIsMedovukha.Error()})
		return
	}

	cli, err := docker.CreateDockerClient()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Docker client error"})
		fmt.Printf("Docker client error: %s\n", err.Error())
		return
	}
	defer cli.Close()

	if err := containers.StartContainerByID(cli, containerId.Id); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "StartContainerByID error"})
		fmt.Printf("StartContainerByID error: %s\n", err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Started: " + containerId.Id})
}

func StopContainerByID(c *gin.Context) {
	var containerId struct {
		Id string `json:"id"`
	}
	if err := c.BindJSON(&containerId); err != nil {
		return
	}

	if check, err := containers.CheckIsMedovukhaId(containerId.Id); err != nil {
		return
	} else if check {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": types.ErrContainerIsMedovukha.Error()})
		return
	}

	cli, err := docker.CreateDockerClient()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Docker client error"})
		fmt.Printf("Docker client error: %s\n", err.Error())
		return
	}
	defer cli.Close()

	if err := containers.StopContainerByID(cli, containerId.Id); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "StopContainerByID error"})
		fmt.Printf("StopContainerByID error: %s\n", err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Stopped: " + containerId.Id})
}

func RestartContainerByID(c *gin.Context) {
	var containerId struct {
		Id string `json:"id"`
	}
	if err := c.BindJSON(&containerId); err != nil {
		return
	}

	if check, err := containers.CheckIsMedovukhaId(containerId.Id); err != nil {
		return
	} else if check {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": types.ErrContainerIsMedovukha.Error()})
		return
	}

	cli, err := docker.CreateDockerClient()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Docker client error"})
		fmt.Printf("Docker client error: %s\n", err.Error())
		return
	}
	defer cli.Close()

	if err := containers.RestartContainerByID(cli, containerId.Id); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "RestartContainerByID error"})
		fmt.Printf("RestartContainerByID error: %s\n", err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Restarted: " + containerId.Id})
}

func RemoveContainerByID(c *gin.Context) {
	var containerId struct {
		Id string `json:"id"`
	}
	if err := c.BindJSON(&containerId); err != nil {
		return
	}

	if check, err := containers.CheckIsMedovukhaId(containerId.Id); err != nil {
		return
	} else if check {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": types.ErrContainerIsMedovukha.Error()})
		return
	}

	cli, err := docker.CreateDockerClient()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Docker client error"})
		fmt.Printf("Docker client error: %s\n", err.Error())
		return
	}
	defer cli.Close()

	if err := containers.RemoveContainerByID(cli, containerId.Id); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "RemoveContainerByID error"})
		fmt.Printf("RemoveContainerByID error: %s\n", err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Removed: " + containerId.Id})
}

//Images

func GetImageList(c *gin.Context) {
	cli, err := docker.CreateDockerClient()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Docker client error"})
		fmt.Printf("Docker client error: %s\n", err.Error())
		return
	}
	defer cli.Close()

	imgList, err := images.GetImageList(cli)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "GetImageList error"})
		fmt.Printf("GetImageList error: %s\n", err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, imgList)
}

// Networks
func GetNetworkList(c *gin.Context) {
	cli, err := docker.CreateDockerClient()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Docker client error"})
		fmt.Printf("Docker client error: %s\n", err.Error())
		return
	}
	defer cli.Close()

	imgList, err := networks.GetNetworkList(cli)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "GetNetworkList error"})
		fmt.Printf("GetNetworkList error: %s\n", err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, imgList)
}

// Volumes
func GetVolumeList(c *gin.Context) {
	cli, err := docker.CreateDockerClient()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Docker client error"})
		fmt.Printf("Docker client error: %s\n", err.Error())
		return
	}
	defer cli.Close()

	imgList, err := volumes.GetVolumeList(cli)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "GetVolumeList error"})
		fmt.Printf("GetVolumeList error: %s\n", err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, imgList)
}

// Deploy
func CreateFromGit(c *gin.Context) {
	ctx := context.Background()
	//
	//
	// Parsing Request
	//
	//
	var NewDeploy types.DeployFromGit
	if err := c.BindJSON(&NewDeploy); err != nil {
		return
	}
	//
	//
	// Create Docker Client
	//
	//
	cli, err := docker.CreateDockerClient()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Docker client error"})
		fmt.Printf("Docker client error: %s\n", err.Error())
		return
	}
	defer cli.Close()
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

	if err := git.CloneRepo(&git.RepoCloner{}, NewDeploy.URL, tempDir); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "CloneRepo error"})
		fmt.Printf("CloneRepo error: %s\n", err.Error())
		return
	}
	//
	//
	// Check Dockerfile in request
	//
	//
	dockerfileDir := filepath.Join(tempDir, "Dockerfile")
	dockercomposeDir := filepath.Join(tempDir, "docker-compose.yml")
	if NewDeploy.Dockerfile == "" {
		// Dockerfile is not specified in the request, check in the directory
		// If Dockerfile doesn`t exist, throw error
		if !fileExists(dockerfileDir) {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Dockerfile empty error"})
			fmt.Println("Dockerfile empty error")
			return
		}
	} else {
		// If Dockerfile specified in the request, rewrite/create new Dockerfile
		if err := createFile(dockerfileDir, []byte(NewDeploy.Dockerfile)); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Dockerfile write error"})
			fmt.Printf("Dockerfile write error: %s\n", err.Error())
			return
		} else {
			fmt.Println("Dockerfile redefined")
		}
	}
	//
	//
	// Parse tags from repo (for image name)
	//
	//
	repo, err := repoFromURL(NewDeploy.URL)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Repo error"})
		fmt.Printf("Repo error: %s\n", err.Error())
		return
	}
	tags := []string{repo + ":latest"} //NewDeploy.URL
	//
	//
	// Delete old containers
	//
	//
	if err := containers.RemoveContainerByImage(ctx, cli, tags[0]); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error remove container"})
		fmt.Printf("Error remove container: %s\n", err.Error())
		return
	}
	//
	//
	// Delete old images
	//
	//
	if err := images.RemoveImageByTag(ctx, cli, tags[0]); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error remove image"})
		fmt.Printf("Error remove image: %s\n", err.Error())
		return
	}
	//
	//
	// Build image
	//
	//
	newImageId, err := images.BuildImageNew(cli, tempDir, tags)
	if err != nil {
		fmt.Println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Docker client error"})
		fmt.Printf("Docker client error: %s\n", err.Error())
		return
	}
	//
	//
	// Clear build cache
	//
	//
	pruneFilters := filters.NewArgs()
	pruneFilters.Add("dangling", "true")

	_, err = cli.ImagesPrune(ctx, pruneFilters)
	if err != nil {
		fmt.Println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Docker client error"})
		fmt.Printf("Docker client error: %s\n", err.Error())
		return
	}
	_, err = cli.BuildCachePrune(ctx, build.CachePruneOptions{All: true})
	if err != nil {
		fmt.Println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Docker client error"})
		fmt.Printf("Docker client error: %s\n", err.Error())
		return
	}
	//
	//
	// Check DockerCompose in request
	//
	//
	if NewDeploy.DockerCompose == "" {
		// DockerCompose is not specified in the request, check DockerRun in the request
		if NewDeploy.DockerRun == "" {
			// DockerRun is not specified in the request, check DockerCompose in the directory
			if !fileExists(dockercomposeDir) {
				c.IndentedJSON(http.StatusOK, gin.H{"message": "created but not launched: " + newImageId})
				fmt.Println("Created image, but not command to launch container")
				return
			} else {
				if err := containers.ExecDockerComposeUp(dockercomposeDir); err != nil {
					c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "the image was created, but an error occurred when starting the container (dockerCompose)"})
					return
				}
			}
		} else {
			if err := containers.ExecDockerRun(NewDeploy.DockerRun); err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "the image was created, but an error occurred when starting the container (dockerRun)"})
				return
			}
		}
	} else {
		// If Dockerfile specified in the request, rewrite/create new Dockerfile
		if err := createFile(dockercomposeDir, []byte(NewDeploy.DockerCompose)); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "DockerCompose write error"})
			fmt.Printf("DockerCompose write error: %s\n", err.Error())
			return
		} else {
			fmt.Println("DockerCompose redefined")
			if err := containers.ExecDockerComposeUp(dockercomposeDir); err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "the image was created, but an error occurred when starting the container (dockerCompose)"})
				return
			}
		}
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "created and launched: " + newImageId})
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
