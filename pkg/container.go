package pkg

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	dclient "github.com/docker/docker/client"
	"log"
	"time"
)

type Spec struct {
	Name string
	Image string
	Mounts []Mount
}

type Mount struct {
	Source string
	Target string
	ReadOnly bool
}

func (m *Mount) toDockerMount () mount.Mount {
	var mount mount.Mount = mount.Mount{
		Type:     "bind",
		Source:   m.Source,
		Target:   m.Target,
		ReadOnly: false,
	}
	return mount
}


func Start(spec Spec) error {
	log.Println("Starting Container [" + spec.Name + "]")
	ctx := context.Background()
	cli, err := dclient.NewEnvClient()
	if err != nil {
		return err
	}

	m := make([]mount.Mount, 0, len(spec.Mounts))
	for _, _m := range spec.Mounts {
		m = append(m, _m.toDockerMount())
	}

	hostConfig := container.HostConfig{
		Mounts: m,
	}

	config := container.Config{
		Image:   spec.Image,
		Tty:     true,
	}

	if resp, err := cli.ContainerCreate(ctx, &config, &hostConfig, nil, spec.Name); err != nil {
		log.Println("Error while creating container [" + spec.Name + "]")
	} else {
		log.Println("Successfully created container [" + spec.Name + "]")
		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			return err
		} else {
			log.Println("Started container [" + spec.Name + "]")
		}
	}
	return err
}

func Stop(name string) error {
	log.Println("Stopping Container [" + name + "]")
	ctx := context.Background()
	cli, err := dclient.NewEnvClient()
	if err != nil {
		return err
	}
	var timeout time.Duration = time.Duration(5 * time.Second)
	if err = cli.ContainerStop(ctx, name, &timeout); err != nil {
		log.Println("Error while try stop container [" + name + "]")
	}
	if err = cli.ContainerRemove(ctx, name, types.ContainerRemoveOptions{}); err != nil {
		log.Println("Error while try remove container [" + name + "]")
		return err
	}
	return err
}

func Exec(container string, command string) error {
	log.Println("Executing command [" + command + "] inside container [" + container + "]")
	ctx := context.Background()
	cli, err := dclient.NewEnvClient()
	if err != nil {
		return err
	}
	config := types.ExecConfig{
		AttachStdout:true,
		Detach:true,
		Cmd:          []string{
			"/bin/bash",
			"-c",
			command,
		},
	}

	if response, err := cli.ContainerExecCreate(ctx, container, config); err != nil {
		log.Println("Error while create exec instante to container [" + container + "]")
		return err
	} else {
		log.Println("Successful creation of exec instante to container [" + container + "]")
		execId := response.ID
		if hijacker, err := cli.ContainerExecAttach(ctx, execId, config); err != nil {
			log.Println("Error while start exec instance [" + execId + "]")
		} else {
			if err := cli.ContainerExecStart(ctx, execId, types.ExecStartCheck{}); err != nil {
				log.Println("Error while start exec [" + execId + "]")
			} else {
				var p []byte = make([]byte, 0, 5000)
				_, err := hijacker.Conn.Read(p)
				if err != nil {
					log.Println(err.Error())
				}
				//size := hijacker.Reader.Size()
				//_, err = hijacker.Reader.Read(p)
				//if err != nil {
				//	log.Println(err.Error())
				//}
				s := string(p)
				log.Println("Output:\n" + s)
			}
		}
		return err
	}
}
