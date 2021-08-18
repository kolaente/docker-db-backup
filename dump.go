package main

import "github.com/docker/docker/api/types"

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
