package main

import (
	"github.com/docker/docker/api/types"
	"strings"
)

type Dumper interface {
	Dump() error
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

func dumpAllDatabases() error {
	lock.Lock()
	defer lock.Unlock()

	for _, dumper := range store {
		err := dumper.Dump()
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
