package containers

import (
	"context"
	"fmt"
	"log"
	"medovukha/api/rest/v1/types"
	dc "medovukha/services/docker"
	image "medovukha/services/docker/images"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/go-connections/nat"
)

func GetContainerBaseInfoList(cli dc.IDockerClient) ([]types.ContainerBaseInfo, error) {
	ctx := context.Background()

	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	conList := make([]types.ContainerBaseInfo, len(containers))
	for i, container := range containers {
		conList[i] = types.ContainerBaseInfo{
			Id:        container.ID,
			Names:     container.Names,
			ImageName: container.Image,
			Created:   container.Created,
			State:     container.State,
		}
		if len(container.Ports) == 0 {
			conList[i].Ports = nil
		} else {
			conList[i].Ports = make([]types.Port, len(container.Ports))
			for j, IPitem := range container.Ports {
				conList[i].Ports[j].IP = IPitem.IP
				conList[i].Ports[j].PrivatePort = IPitem.PrivatePort
				conList[i].Ports[j].PublicPort = IPitem.PublicPort
				conList[i].Ports[j].Type = IPitem.Type
			}
		}
		if check, err := CheckIsMedovukhaId(conList[i].Id); err != nil {
			return nil, err
		} else {
			conList[i].IsMedovukha = check
		}
	}

	return conList, nil
}

func ExecDockerRun(dockerRunCommand string) error {
	ctx := context.Background()

	args := strings.Split(dockerRunCommand, " ")
	if len(args) < 2 {
		return fmt.Errorf("wrong dockerRunCommand syntax: %v", args)
	}

	cmd := exec.CommandContext(ctx, "docker", args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run failed: %s", err.Error())
	}

	return nil
}

func ExecDockerComposeUp(composeFilepath string) error {
	ctx := context.Background()

	args := []string{
		"compose",
		"-f", composeFilepath,
		"up", "-d",
	}

	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("composeUp failed: %s", err.Error())
	}

	return nil
}

func CreateTestContainer(cli dc.IDockerClient) error {
	ctx := context.Background()

	imageName := "docker/welcome-to-docker"

	if err := image.PullImage(cli, ctx, imageName); err != nil {
		return err
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			// original port
			"80/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "9990", // new port
				},
			},
		},
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Tty:   false,
	}, hostConfig, nil, nil, "web-test")
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return err
	}

	fmt.Println(resp.ID)
	return nil
}

func PauseContainerByID(cli dc.IDockerClient, id string) error {
	ctx := context.Background()

	conList, err := GetContainerBaseInfoList(cli)
	if err != nil {
		return err
	}

	for _, container := range conList {
		if container.Id == id {
			if err := cli.ContainerPause(ctx, container.Id); err != nil {
				return err
			}
			fmt.Println("Paused: ", container.Id)
			return nil
		}
	}
	fmt.Println("Not found: ", id)
	return types.ErrContainerNotFound
}

func UnpauseContainerByID(cli dc.IDockerClient, id string) error {
	ctx := context.Background()

	conList, err := GetContainerBaseInfoList(cli)
	if err != nil {
		return err
	}

	for _, container := range conList {
		if container.Id == id {
			if err := cli.ContainerUnpause(ctx, container.Id); err != nil {
				return err
			}
			fmt.Println("Unpaused: ", container.Id)
			return nil
		}
	}
	fmt.Println("Not found: ", id)
	return types.ErrContainerNotFound
}

func KillContainerByID(cli dc.IDockerClient, id string) error {
	ctx := context.Background()

	conList, err := GetContainerBaseInfoList(cli)
	if err != nil {
		return err
	}

	for _, container := range conList {
		if container.Id == id {
			if err := cli.ContainerKill(ctx, container.Id, ""); err != nil {
				return err
			}
			fmt.Println("Killed: ", container.Id)
			return nil
		}
	}
	fmt.Println("Not found: ", id)
	return types.ErrContainerNotFound
}

func StartContainerByID(cli dc.IDockerClient, id string) error {
	ctx := context.Background()

	conList, err := GetContainerBaseInfoList(cli)
	if err != nil {
		return err
	}

	for _, con := range conList {
		if con.Id == id {
			if err := cli.ContainerStart(ctx, con.Id, container.StartOptions{}); err != nil {
				return err
			}
			fmt.Println("Started: ", con.Id)
			return nil
		}
	}
	fmt.Println("Not found: ", id)
	return types.ErrContainerNotFound
}

func RestartContainerByID(cli dc.IDockerClient, id string) error {
	ctx := context.Background()

	conList, err := GetContainerBaseInfoList(cli)
	if err != nil {
		return err
	}

	for _, con := range conList {
		if con.Id == id {
			if err := cli.ContainerRestart(ctx, con.Id, container.StopOptions{}); err != nil {
				return err
			}
			fmt.Println("Restarted: ", con.Id)
			return nil
		}
	}
	fmt.Println("Not found: ", id)
	return types.ErrContainerNotFound
}

func StopContainerByID(cli dc.IDockerClient, id string) error {
	ctx := context.Background()

	conList, err := GetContainerBaseInfoList(cli)
	if err != nil {
		return err
	}

	for _, con := range conList {
		if con.Id == id {
			if err := cli.ContainerStop(ctx, con.Id, container.StopOptions{}); err != nil {
				return err
			}
			fmt.Println("Stopped: ", con.Id)
			return nil
		}
	}
	fmt.Println("Not found: ", id)
	return types.ErrContainerNotFound
}

func RemoveContainerByID(cli dc.IDockerClient, id string) error {
	ctx := context.Background()

	conList, err := GetContainerBaseInfoList(cli)
	if err != nil {
		return err
	}

	for _, con := range conList {
		if con.Id == id {
			if err := cli.ContainerRemove(ctx, con.Id, container.RemoveOptions{
				RemoveVolumes: true,
				RemoveLinks:   false,
				Force:         false,
			}); err != nil {
				return err
			}
			fmt.Println("Started: ", con.Id)
			return nil
		}
	}
	fmt.Println("Not found: ", id)
	return types.ErrContainerNotFound
}

func RemoveContainerByImage(ctx context.Context, cli dc.IDockerClient, imageTag string) error {
	f := filters.NewArgs()
	f.Add("ancestor", imageTag)

	containers, err := cli.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: f,
	})
	if err != nil {
		return fmt.Errorf("cannot list containers for image %q: %s", imageTag, err.Error())
	}

	if len(containers) == 0 {
		return nil
	}

	for _, ctr := range containers {
		if ctr.State == container.StateRunning {
			if err := cli.ContainerStop(ctx, ctr.ID,
				container.StopOptions{Timeout: nil}); err != nil {
				return fmt.Errorf("cannot stop container %s: %s", ctr.ID, err.Error())
			}
		}

		if err := cli.ContainerRemove(ctx, ctr.ID, container.RemoveOptions{Force: true}); err != nil {
			return fmt.Errorf("cannot remove container %s: %s", ctr.ID, err.Error())
		}

		log.Printf("container %s removed\n", ctr.ID)
	}

	return nil
}

func CheckIsMedovukhaId(id string) (bool, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return false, err
	}
	//the hostname is part of the full container ID
	contains := strings.Contains(id, hostname)
	i := strings.Index(id, hostname)
	if contains && i == 0 {
		return true, nil
	}
	return false, nil
}
