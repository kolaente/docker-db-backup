package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func runAndSaveCommand(filename, command string, args ...string) error {
	c := exec.Command(command, args...)

	//fmt.Printf("Running %s\n\n", c.String())

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}

	var stderr bytes.Buffer
	c.Stderr = &stderr

	err = c.Start()
	if err != nil {
		return err
	}

	_, err = io.Copy(f, stdout)
	if err != nil {
		return err
	}

	err = c.Wait()
	if err != nil {
		fmt.Printf(stderr.String())
		return err
	}

	return nil
}
