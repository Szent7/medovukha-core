package containers

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/Szent7/medovukha-core/ipc/types"
	"github.com/Szent7/medovukha-core/services/common"
	dc "github.com/Szent7/medovukha-core/services/docker"
	image "github.com/Szent7/medovukha-core/services/docker/images"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
)

func GetContainer(ctx context.Context, cli dc.IDockerClient, containerID string) (types.ContainerBaseInfo, error) {
	containerInspect, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return types.ContainerBaseInfo{}, err
	}

	formattedTime, _ := common.StrToUNIXTime(containerInspect.Created)

	conSummary := types.ContainerBaseInfo{
		Id:        containerInspect.ID,
		Names:     []string{containerInspect.Name},
		ImageName: containerInspect.Image,
		Created:   formattedTime,
		State:     containerInspect.State.Status,
	}
	if containerInspect.HostConfig.PortBindings != nil {
		conSummary.Ports = make([]types.Port, len(containerInspect.HostConfig.PortBindings))
		for i := 0; i < len(containerInspect.HostConfig.PortBindings); i++ {
			for natPort, portBindings := range containerInspect.HostConfig.PortBindings {
				for _, pb := range portBindings {
					parsedHostPort, _ := common.StrToUint16(pb.HostPort)
					parsedContainerPort, _ := common.StrToUint16(natPort.Port())
					conSummary.Ports[i].IP = pb.HostIP
					conSummary.Ports[i].PublicPort = parsedHostPort
					conSummary.Ports[i].PrivatePort = parsedContainerPort
					conSummary.Ports[i].Type = natPort.Proto()
				}
			}
		}
	}

	return conSummary, nil
}

func GetContainerBaseInfoList(ctx context.Context, cli dc.IDockerClient) ([]types.ContainerBaseInfo, error) {
	containers, err := GetContainerRawList(ctx, cli)
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
	}

	return conList, nil
}

func IsImageUsed(ctx context.Context, cli dc.IDockerClient, imageID string) (bool, error) {
	containers, err := GetContainerRawListByImage(ctx, cli, imageID)
	if err != nil {
		return false, err
	}

	return len(containers) != 0, nil
}

func IsVolumeUsed(ctx context.Context, cli dc.IDockerClient, volumeName string) (bool, error) {
	containers, err := GetContainerRawList(ctx, cli)
	if err != nil {
		return false, err
	}

	if len(containers) != 0 {
		for _, con := range containers {
			for _, m := range con.Mounts {
				if m.Type == mount.TypeVolume && m.Name == volumeName {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func GetContainerRawList(ctx context.Context, cli dc.IDockerClient) ([]container.Summary, error) {
	return cli.ContainerList(ctx, container.ListOptions{All: true})
}

func GetContainerRawListByImage(ctx context.Context, cli dc.IDockerClient, imageID string) ([]container.Summary, error) {
	f := filters.NewArgs()
	f.Add("ancestor", imageID)

	return cli.ContainerList(ctx, container.ListOptions{All: true, Filters: f})
}

func ExecDockerRun(ctx context.Context, dockerRunCommand string, logCh chan string) error {
	args := strings.Split(dockerRunCommand, " ")
	if len(args) < 2 {
		return fmt.Errorf("wrong dockerRunCommand syntax: %v", args)
	}

	cmd := exec.CommandContext(ctx, "docker", args[1:]...)
	pr, pw := io.Pipe()
	cmd.Stdout = pw
	cmd.Stderr = pw

	go func() {
		scanner := bufio.NewScanner(pr)
		for scanner.Scan() {
			common.SendLog(logCh, scanner.Text())
		}
	}()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run failed: %s", err.Error())
	}

	return nil
}

func ExecDockerComposeUp(ctx context.Context, composeFilepath string, logCh chan string) error {
	args := []string{
		"compose",
		"-f", composeFilepath,
		"up", "-d",
	}

	cmd := exec.CommandContext(ctx, "docker", args...)
	pr, pw := io.Pipe()
	cmd.Stdout = pw
	cmd.Stderr = pw

	go func() {
		scanner := bufio.NewScanner(pr)
		for scanner.Scan() {
			common.SendLog(logCh, scanner.Text())
		}
	}()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("composeUp failed: %s", err.Error())
	}

	return nil
}

func CreateTestContainer(ctx context.Context, cli dc.IDockerClient) error {
	imageName := "docker/welcome-to-docker"

	if err := image.PullImage(ctx, cli, imageName); err != nil {
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

func PauseContainerByID(ctx context.Context, cli dc.IDockerClient, id string) error {
	conList, err := GetContainerBaseInfoList(ctx, cli)
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

func UnpauseContainerByID(ctx context.Context, cli dc.IDockerClient, id string) error {
	conList, err := GetContainerBaseInfoList(ctx, cli)
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

func KillContainerByID(ctx context.Context, cli dc.IDockerClient, id string) error {
	conList, err := GetContainerBaseInfoList(ctx, cli)
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

func StartContainerByID(ctx context.Context, cli dc.IDockerClient, id string) error {
	conList, err := GetContainerBaseInfoList(ctx, cli)
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

func RestartContainerByID(ctx context.Context, cli dc.IDockerClient, id string) error {
	conList, err := GetContainerBaseInfoList(ctx, cli)
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

func StopContainerByID(ctx context.Context, cli dc.IDockerClient, id string) error {
	conList, err := GetContainerBaseInfoList(ctx, cli)
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

func RemoveContainerByID(ctx context.Context, cli dc.IDockerClient, id string) error {
	conList, err := GetContainerBaseInfoList(ctx, cli)
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
			fmt.Println("Removed: ", con.Id)
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
