package main

import (
	"log"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

const containerLabelName = `de.kolaente.db-backup`

type Dumper interface {
	Dump(c *client.Client) error
}

func NewDumperFromContainer(ctr *container.InspectResponse) Dumper {

	// Containers contain the tags, therefore we need to check them one by one
	if strings.HasPrefix(ctr.Config.Image, "mysql") || strings.HasPrefix(ctr.Config.Image, "mariadb") || ctr.Config.Labels[containerLabelName] == "mysql" {
		return NewMysqlDumper(ctr)
	}

	if strings.HasPrefix(ctr.Config.Image, "postgres") || ctr.Config.Labels[containerLabelName] == "postgres" {
		return NewPostgresDumper(ctr)
	}

	return nil
}

func dumpAllDatabases(c *client.Client) {
	lock.Lock()
	defer lock.Unlock()

	for containerID, dumper := range store {
		err := dumper.Dump(c)
		if err != nil {
			log.Printf("Could not dump database from container %s: %v", containerID, err)
		}
	}
}

func getDumpFilename(containerName string) string {
	if strings.HasPrefix(containerName, "/") {
		containerName = strings.TrimPrefix(containerName, "/")
	}

	return config.fullCurrentBackupPath + containerName + ".sql"
}
