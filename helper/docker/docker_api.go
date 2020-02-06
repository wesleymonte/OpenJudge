package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

type ContainerConfig struct {
	Name string
	Image string
	Mounts []mount.Mount
}

func CreateContainer(cli *client.Client, config ContainerConfig) (string, error) {
	log.Printf("Creating Container [%s]", config.Name)
	ctx := context.Background()
	hostConfig := container.HostConfig{
		Mounts: config.Mounts,
	}

	dconfig := container.Config{
		Image:   config.Image,
		Tty:     true,
	}

	b, err := cli.ContainerCreate(ctx, &dconfig, &hostConfig, nil, config.Name)
	return b.ID, err
}

func StartContainer(cli *client.Client, id string) error {
	log.Printf("Starting Container [%s]", id)
	return cli.ContainerStart(context.Background(), id, types.ContainerStartOptions{})
}

func StopContainer(cli *client.Client, id string) error {
	log.Printf("Stopping Container [%s]", id)
	var timeout = 5 * time.Second
	return cli.ContainerStop(context.Background(), id, &timeout)
}

func RemoveContainer(cli *client.Client, id string) error {
	log.Printf("Removing Container [%s]", id)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return cli.ContainerRemove(ctx, id, types.ContainerRemoveOptions{})
}

func WriteToFile(cli *client.Client, containerId string, content []string, dest string) error {
	for _, c := range content {
		c = strings.ReplaceAll(c, "'", "'\"'\"'")
		cmd := fmt.Sprintf(`echo -E '%s' >> %s`, c, dest)
		log.Printf("Writing [%s] on [%s] from Container [%s]", c, dest, containerId)
		err := Exec(cli, containerId, cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

func Exec(cli *client.Client, containerId, cmd string)  error {
	log.Printf("Executing command [%s] on container [%s]", cmd, containerId)
	config := types.ExecConfig{
		Cmd: []string{"/bin/bash", "-c", cmd},
	}
	rid, _ := cli.ContainerExecCreate(context.Background(), containerId, config)
	return cli.ContainerExecStart(context.Background(), rid.ID, types.ExecStartCheck{})
}

func GetFileContent(cli *client.Client, id, path string) ([]byte, error) {
	log.Printf("Getting content of file [%s]", path)
	config := types.ExecConfig{
		Tty:true,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:[]string{"/bin/bash", "-c", "cat " + path},
	}
	rid, _ := cli.ContainerExecCreate(context.Background(), id, config)
	hijack, _ := cli.ContainerExecAttach(context.Background(), rid.ID, types.ExecConfig{Tty: true,})
	output := read(hijack.Conn)
	return output, nil
}

func read(conn net.Conn) []byte {
	result := make([]byte, 0)
	b := make([]byte, 10)
	for ; ; {
		n, _ := conn.Read(b)
		result = append(result, b...)
		if n < len(b) {break}
		b = make([]byte, 2 * len(b))
	}
	return result
}

func Pull(cli *client.Client, image string) (io.ReadCloser, error) {
	log.Printf("Pulling image [%s]", image)
	reader, err := cli.ImagePull(context.Background(), image, types.ImagePullOptions{})
	return reader, err
}

func CheckImage(cli *client.Client, image string) (exist bool, err error)  {
	log.Printf("Checking image [%s]", image)
	exist = false
	_, _, err = cli.ImageInspectWithRaw(context.Background(), image)
	if err == nil {
		exist = true
	}
	return
}