package main

import (
	"reflect"
	"testing"

	"github.com/docker/docker/api/types/container"
)

func TestNewDumperFromContainer(t *testing.T) {
	tests := []struct {
		name string
		ctr  *container.InspectResponse
		want string
	}{
		{
			name: "mysql",
			ctr:  &container.InspectResponse{Config: &container.Config{Image: "mysql"}},
			want: "MysqlDumper",
		},
		{
			name: "mariadb",
			ctr:  &container.InspectResponse{Config: &container.Config{Image: "mysql"}},
			want: "MysqlDumper",
		},
		{
			name: "postgres",
			ctr:  &container.InspectResponse{Config: &container.Config{Image: "postgres"}},
			want: "PostgresDumper",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDumperFromContainer(tt.ctr)
			if reflect.TypeOf(got).Elem().Name() != tt.want {
				t.Errorf("NewDumperFromContainer() is %s, want %s", got, tt.want)
			}
		})
	}
}
