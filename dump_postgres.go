package main

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type PostgresDumper struct {
	Container *types.ContainerJSON
}

func NewPostgresDumper(container *types.ContainerJSON) *PostgresDumper {
	return &PostgresDumper{
		Container: container,
	}
}

func (d *PostgresDumper) buildConnStr() string {
	env := parseEnv(d.Container.Config.Env)

	user := "root"
	if u, has := env["POSTGRES_USER"]; has {
		user = u
	}

	db := "postgres"
	if d, has := env["POSTGRES_DB"]; has {
		db = d
	}

	pw := env["POSTGRES_ROOT_PASSWORD"]
	if p, has := env["POSTGRES_PASSWORD"]; has {
		pw = p
	}

	port := "5432"
	if p, has := env["POSTGRES_PORT"]; has {
		port = p
	}

	host := d.Container.NetworkSettings.DefaultNetworkSettings.IPAddress

	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, pw, host, port, db)
}

func (d *PostgresDumper) Dump(c *client.Client) error {
	fmt.Printf("Dumping postgres database from container %s...\n", d.Container.Name)

	connStr := d.buildConnStr()

	return runAndSaveCommandInContainer(getDumpFilename(d.Container.Name), c, d.Container, "pg_dump", "--dbname", connStr)
}
