package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func runAndSaveCommandInContainer(filename string, c *client.Client, container *types.ContainerJSON, command string, args ...string) error {
	ctx := context.Background()

	config := types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          append([]string{command}, args...),
	}

	r, err := c.ContainerExecCreate(ctx, container.ID, config)
	if err != nil {
		return err
	}

	resp, err := c.ContainerExecAttach(ctx, r.ID, types.ExecStartCheck{})
	if err != nil {
		return err
	}
	defer resp.Close()

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Reader)
	if err != nil {
		return err
	}

	execInspect, err := c.ContainerExecInspect(ctx, r.ID)
	if execInspect.ExitCode != 0 {
		return fmt.Errorf("backup from container %s failed with exit code %d", container.Name, execInspect.ExitCode)
	}
	return err
}
