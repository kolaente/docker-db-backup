package main

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"log"
)

type PostgresDumper struct {
	Container *types.ContainerJSON
}

func NewPostgresDumper(container *types.ContainerJSON) *PostgresDumper {
	return &PostgresDumper{
		Container: container,
	}
}

func (d *PostgresDumper) Dump(c *client.Client) error {
	log.Printf("Dumping postgres database from container %s...\n", d.Container.Name)

	env := parseEnv(d.Container.Config.Env)

	user := "root"
	if u, has := env["POSTGRES_USER"]; has {
		user = u
	}

	return runAndSaveCommandInContainer(c, d.Container, "pg_dumpall", "-U", user)
}
