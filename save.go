package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

func runAndSaveCommandInContainer(filename string, c *client.Client, container *types.ContainerJSON, command string, args ...string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

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

	// read the output
	var outBuf, errBuf bytes.Buffer
	outputDone := make(chan error)

	go func() {
		// StdCopy demultiplexes the stream into two buffers
		_, err = stdcopy.StdCopy(&outBuf, &errBuf, resp.Reader)
		outputDone <- err
	}()

	select {
	case err := <-outputDone:
		if err != nil {
			return err
		}
		break

	case <-ctx.Done():
		return ctx.Err()
	}

	_, err = c.ContainerExecInspect(ctx, r.ID)
	if err != nil {
		fmt.Printf(errBuf.String())
		return err
	}

	_, err = io.Copy(f, &outBuf)
	if err != nil {
		return err
	}

	return nil
}
