package images

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	ts "medovukha/api/rest/v1/types"
	dc "medovukha/services/docker"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/image"
)

func PullImage(cli dc.IDockerClient, ctx context.Context, imageName string) error {
	out, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(os.Stdout, out)
	return nil
}

func GetImageList(cli dc.IDockerClient) ([]ts.ImageBaseInfo, error) {
	images, err := GetRawImageList(cli)
	if err != nil {
		return nil, err
	}

	imgList := make([]ts.ImageBaseInfo, len(images))
	for i, image := range images {
		imgList[i] = ts.ImageBaseInfo{
			Id:      image.ID,
			Tags:    image.RepoTags,
			Size:    image.Size,
			Created: image.Created,
		}
	}

	return imgList, nil
}

func GetRawImageList(cli dc.IDockerClient) ([]image.Summary, error) {
	ctx := context.Background()

	images, err := cli.ImageList(ctx, image.ListOptions{All: true, ContainerCount: true})
	if err != nil {
		return nil, err
	}

	return images, nil
}

func GetImageIDFromTag(cli dc.IDockerClient, tag string) (string, error) {
	images, err := GetRawImageList(cli)
	imageID := ""
	if err != nil {
		return imageID, err
	}
	for _, image := range images {
		if len(image.RepoTags) != 0 {
			for _, repoTag := range image.RepoTags {
				if tag == repoTag {
					imageID = image.ID
					return imageID, nil
				}
			}
		}
	}
	return imageID, nil
}

func GetImageListIDs(cli dc.IDockerClient) ([]string, error) {
	images, err := GetRawImageList(cli)
	if err != nil {
		return nil, err
	}
	imageListIds := make([]string, len(images))
	for i, image := range images {
		imageListIds[i] = image.ID
	}
	return imageListIds, nil
}

func RemoveImageByTag(ctx context.Context, cli dc.IDockerClient, imageTag string) error {
	images, err := cli.ImageList(ctx, image.ListOptions{All: true})
	if err != nil {
		return fmt.Errorf("cannot list images: %s", err.Error())
	}

	found := false
	for _, img := range images {
		for _, tag := range img.RepoTags {
			if tag == imageTag {
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		return nil
	}

	_, err = cli.ImageRemove(ctx, imageTag, image.RemoveOptions{
		Force:         true,
		PruneChildren: true,
	})
	if err != nil {
		return fmt.Errorf("cannot remove image %s: %s", imageTag, err.Error())
	}

	log.Printf("image %q removed\n", imageTag)

	return nil
}

func BuildImageNew(cli dc.IDockerClient, path string, tags []string) (string, error) {
	ctx := context.Background()

	if len(tags) < 1 {
		return "", ts.ErrEmptyTags
	}

	args := []string{
		"build",
		"--pull",
		"-t", tags[0],
		path,
	}

	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("build failed: %s", err.Error())
	}

	return GetImageIDFromTag(cli, tags[0])
}

// ? Perhaps archiving should be made into a separate function (ITarArchiver)
func BuildImage(cli dc.IDockerClient, tarArchiver ITarArchiver, path string, tags []string) (string, error) {
	ctx := context.Background()

	if len(tags) < 1 {
		return "", ts.ErrEmptyTags
	}

	buf := new(bytes.Buffer)
	if err := tarDirectory(path, buf); err != nil {
		return "", fmt.Errorf("tar directory: %w", err)
	}

	out, err := cli.ImageBuild(ctx, buf, build.ImageBuildOptions{
		Dockerfile:  "Dockerfile",
		Tags:        tags,
		Remove:      true,
		ForceRemove: true,
		NoCache:     true,
	})
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	defer out.Body.Close()
	io.Copy(os.Stdout, out.Body)

	return GetImageIDFromTag(cli, tags[0])
}

// RemoveIntermediateImages get a list of id images to be deleted
// returns a list of undeleted images
// if all images have been deleted, the list will be empty
func RemoveIntermediateImages(client dc.IDockerClient, IDs []string) []string {
	opt := image.RemoveOptions{Force: true, PruneChildren: true}
	var err error = nil
	undeletedIDs := make([]string, 0, len(IDs))
	for _, image := range IDs {
		_, err = client.ImageRemove(context.Background(), image, opt)
		if err != nil {
			undeletedIDs = append(undeletedIDs, image)
			log.Printf("failed to remove image %s: %v", image, err)
		} else {
			log.Printf("removed intermediate image %s", image)
		}
	}
	return undeletedIDs
}

func tarDirectory(srcDir string, out io.Writer) error {
	tw := tar.NewWriter(out)
	defer tw.Close()

	absSrc, err := filepath.Abs(srcDir)
	if err != nil {
		return err
	}

	return filepath.Walk(absSrc, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(absSrc, file)
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		if fi.IsDir() {
			header.Name = relPath + "/"
		} else {
			header.Name = relPath
		}

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !fi.IsDir() {
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
		}
		return nil
	})
}
