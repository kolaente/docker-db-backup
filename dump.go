package main

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/google/martian/log"
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

func dumpAllDatabases(c *client.Client) {
	lock.Lock()
	defer lock.Unlock()

	for containerID, dumper := range store {
		err := dumper.Dump(c)
		if err != nil {
			log.Errorf("Could not dump database from container %s: %v", containerID, err)
		}
	}
}

func getDumpFilename(containerName string) string {
	if strings.HasPrefix(containerName, "/") {
		containerName = strings.TrimPrefix(containerName, "/")
	}

	return config.fullCurrentBackupPath + containerName + ".sql"
}
