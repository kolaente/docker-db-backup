package main

import (
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type PostgresDumper struct {
	Container *container.InspectResponse
}

func NewPostgresDumper(ctr *container.InspectResponse) *PostgresDumper {
	return &PostgresDumper{
		Container: ctr,
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
