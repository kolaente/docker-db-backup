package main

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"strings"
)

const containerLabelName = `de.kolaente.db-backup`

type Dumper interface {
	Dump(c *client.Client) error
}

func NewDumperFromContainer(container *types.ContainerJSON) Dumper {

	// Containers contain the tags, therefore we need to check them one by one
	if strings.HasPrefix(container.Config.Image, "mysql") || strings.HasPrefix(container.Config.Image, "mariadb") || container.Config.Labels[containerLabelName] == "mysql" {
		return NewMysqlDumper(container)
	}

	if strings.HasPrefix(container.Config.Image, "postgres") || container.Config.Labels[containerLabelName] == "postgres" {
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
