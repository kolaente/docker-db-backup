package main

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"strings"
)

type Dumper interface {
	Dump(c *client.Client) error
}

func NewDumperFromContainer(container *types.ContainerJSON) Dumper {
	switch container.Config.Image {
	case "mysql":
		fallthrough
	case "mariadb":
		return NewMysqlDumper(container)
	case "postgres":
		return NewPostgresDumper(container)
	}

	return nil
}

func dumpAllDatabases(c *client.Client) error {
	lock.Lock()
	defer lock.Unlock()

	for _, dumper := range store {
		err := dumper.Dump(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func getDumpFilename(containerName string) string {
	if strings.HasPrefix(containerName, "/") {
		containerName = strings.TrimPrefix(containerName, "/")
	}

	return config.fullCurrentBackupPath + containerName + ".sql"
}
