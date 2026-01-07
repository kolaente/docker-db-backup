package main

import (
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type MysqlDumper struct {
	Container *container.InspectResponse
}

func NewMysqlDumper(ctr *container.InspectResponse) *MysqlDumper {
	return &MysqlDumper{
		Container: ctr,
	}
}

func (m *MysqlDumper) buildDumpArgs() []string {
	env := parseEnv(m.Container.Config.Env)

	user := "root"
	if u, has := env["MYSQL_USER"]; has {
		user = u
	}

	db := "--all-databases"
	if d, has := env["MYSQL_DATABASE"]; has {
		db = d
	}

	pw := env["MYSQL_ROOT_PASSWORD"]
	if p, has := env["MYSQL_PASSWORD"]; has {
		pw = p
	}

	port := "3306"
	if p, has := env["MYSQL_PORT"]; has {
		port = p
	}

	host := m.Container.NetworkSettings.DefaultNetworkSettings.IPAddress

	args := []string{
		"--lock-tables=0",
		"--dump-date",
		"--single-transaction",
		"--routines",
		"--triggers",
		"--events",
		"-u", user,
	}

	if pw != "" {
		args = append(args, "-p"+pw)
	}

	return append(args, "--port", port, "-h", host, db)
}

func (m *MysqlDumper) Dump(c *client.Client) error {
	log.Printf("Dumping mysql database from container %s...\n", m.Container.Name)

	args := m.buildDumpArgs()

	return runAndSaveCommandInContainer(c, m.Container, "mysqldump", args...)
}
