package main

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func runAndSaveCommandInContainer(c *client.Client, ctr *container.InspectResponse, command string, args ...string) error {
	ctx := context.Background()

	filename := getDumpFilename(ctr.Name)
	if config.CompressBackups {
		filename += ".gz"
	}

	execConfig := container.ExecOptions{
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          append([]string{command}, args...),
	}

	r, err := c.ContainerExecCreate(ctx, ctr.ID, execConfig)
	if err != nil {
		return err
	}

	resp, err := c.ContainerExecAttach(ctx, r.ID, container.ExecStartOptions{})
	if err != nil {
		return err
	}
	defer resp.Close()

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	var target io.Writer = f

	if config.CompressBackups {
		gw, err := gzip.NewWriterLevel(f, gzip.BestCompression)
		if err != nil {
			return err
		}
		defer gw.Close()
		target = gw
	}

	_, err = io.Copy(target, resp.Reader)
	if err != nil {
		return err
	}

	execInspect, err := c.ContainerExecInspect(ctx, r.ID)
	if err != nil {
		return err
	}
	if execInspect.ExitCode != 0 {
		return fmt.Errorf("backup from container %s failed with exit code %d", ctr.Name, execInspect.ExitCode)
	}
	return err
}
