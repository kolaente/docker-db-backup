package main

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"reflect"
	"testing"
)

func TestNewDumperFromContainer(t *testing.T) {
	tests := []struct {
		name      string
		container *types.ContainerJSON
		want      string
	}{
		{
			name:      "mysql",
			container: &types.ContainerJSON{Config: &container.Config{Image: "mysql"}},
			want:      "MysqlDumper",
		},
		{
			name:      "mariadb",
			container: &types.ContainerJSON{Config: &container.Config{Image: "mysql"}},
			want:      "MysqlDumper",
		},
		{
			name:      "postgres",
			container: &types.ContainerJSON{Config: &container.Config{Image: "postgres"}},
			want:      "PostgresDumper",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDumperFromContainer(tt.container)
			if reflect.TypeOf(got).Elem().Name() != tt.want {
				t.Errorf("NewDumperFromContainer() is %s, want %s", got, tt.want)
			}
		})
	}
}
