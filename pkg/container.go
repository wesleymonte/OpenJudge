package pkg

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	dclient "github.com/docker/docker/client"
	"log"
	"os"
	"os/exec"
	"time"
)

const (
	DockerEngine string = "/usr/bin/docker"
	ExecCommand string = "exec"
	CopyCommand string = "cp"
)


type Spec struct {
	Name string
	Image string
	Mounts []ProblemMount
}

type ProblemMount struct {
	Source string
	Target string
	ReadOnly bool
}

func (m *ProblemMount) toDockerMount () mount.Mount {
	var dockerMount mount.Mount = mount.Mount{
		Type:     "bind",
		Source:   m.Source,
		Target:   m.Target,
		ReadOnly: m.ReadOnly,
	}
	return dockerMount
}

func NewProblemMount(problemId string) ProblemMount {
	source := fmt.Sprintf(os.Getenv(ServiceAbsolutePath) + "/problems/%s", problemId)
	target := fmt.Sprintf("/problems/%s", problemId)
	var problemMount = ProblemMount{
		Source:   source,
		Target:   target,
		ReadOnly: true,
	}
	return problemMount
}


func Start(cli *dclient.Client, spec Spec) (err error) {
	log.Println("Starting Container [" + spec.Name + "]")
	ctx := context.Background()
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

	if _, err = cli.ContainerCreate(ctx, &config, &hostConfig, nil, spec.Name); err != nil {
		log.Println("Error while creating container [" + spec.Name + "]")
	} else {
		log.Println("Successfully created container [" + spec.Name + "]")
		if err = cli.ContainerStart(ctx, spec.Name, types.ContainerStartOptions{}); err != nil {
			log.Println("Error while try start container [" + spec.Name + "]")
		} else {
			log.Println("Successfully started container [" + spec.Name + "]")
		}
	}
	return
}

func Stop(cli *dclient.Client, container string) (err error) {
	log.Println("Stopping Container [" + container + "]")
	ctx := context.Background()
	var timeout = 5 * time.Second
	if err = cli.ContainerStop(ctx, container, &timeout); err != nil {
		log.Println("Error while try stop container [" + container + "]")
	} else {
		log.Println("Successful container stop [" + container + "]")
		if err = cli.ContainerRemove(ctx, container, types.ContainerRemoveOptions{}); err != nil {
			log.Println("Error while try remove container [" + container + "]")
		} else {
			log.Println("Successful container remove")
		}
	}
	return
}

func Exec(container string, command string) (out []byte, err error) {
	log.Println("Executing command [" + command + "] inside container [" + container + "]")
	if out, err = exec.Command(DockerEngine, ExecCommand, "-t", container , command).Output(); err != nil {
		log.Println("Error while executing command [" + command + "] to container [" + container + "]")
	} else {
		log.Println("Successful command execution")
	}
	return
}

func Mkdir(container, dir string) (err error) {
	if err = exec.Command(DockerEngine, ExecCommand, "-t", container, "mkdir", dir).Run(); err != nil {
		log.Println("Error while creating folder")
	} else {
		log.Println("Successful folder creation")
	}
	return
}

func Send(container, src, des string) (err error) {
	log.Println("Sending [" + src + "] to [" + des + "] container [" + container + "].")
	if err = exec.Command(DockerEngine, CopyCommand, src, container + ":" + des).Run(); err != nil {
		log.Println("Error while sending file [" + src + "]")
	} else {
		log.Println("Successful send file")
	}
	return
}

func Run(container, problemId, submissionId string) (result []byte, err error) {
	testsPath := fmt.Sprintf("/problems/%s", problemId)
	scriptPath := "/" + SubmissionsDirName + "/" + "submission-" + submissionId + ".py"
	if result, err = exec.Command(DockerEngine, ExecCommand, "-t", container, "run.sh", testsPath, scriptPath).Output(); err != nil {
		log.Println("Error while executing run.sh to container [" + container + "]")
	} else {
		log.Println("Successful command execution")
	}
	return
}
