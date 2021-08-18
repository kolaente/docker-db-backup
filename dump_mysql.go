package main

import "github.com/docker/docker/api/types"

type MysqlDumper struct {
	Container *types.ContainerJSON
}

func NewMysqlDumper(container *types.ContainerJSON) *MysqlDumper {
	return &MysqlDumper{
		Container: container,
	}
}

func (m *MysqlDumper) Dump() error {
	panic("implement me")
}
