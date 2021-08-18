package main

import (
	"fmt"
	"github.com/docker/docker/api/types"
)

type MysqlDumper struct {
	Container *types.ContainerJSON
}

func NewMysqlDumper(container *types.ContainerJSON) *MysqlDumper {
	return &MysqlDumper{
		Container: container,
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

	args := []string{"--lock-tables=0", "--dump-date", "-u", user}

	if pw != "" {
		args = append(args, "-p"+pw)
	}

	return append(args, "--port", port, "-h", host, db)
}

func (m *MysqlDumper) Dump() error {
	fmt.Printf("Dumping mysql database from container %s...\n", m.Container.Name)

	args := m.buildDumpArgs()

	return runAndSaveCommand(getDumpFilename(m.Container.Name), "mysqldump", args...)
}
