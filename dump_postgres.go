package main

import "github.com/docker/docker/api/types"

type PostgresDumper struct {
	Container *types.ContainerJSON
}

func NewPostgresDumper(container *types.ContainerJSON) *PostgresDumper {
	return &PostgresDumper{
		Container: container,
	}
}

func (p *PostgresDumper) Dump() error {
	panic("implement me")
}
