package main

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"testing"
)

func TestPostgresDumper_buildConnStr(t *testing.T) {
	nw := &types.NetworkSettings{
		DefaultNetworkSettings: types.DefaultNetworkSettings{
			IPAddress: "1.2.3.4",
		},
	}
	tests := []struct {
		name      string
		container *types.ContainerJSON
		want      string
	}{
		{
			name: "values for everything",
			container: &types.ContainerJSON{
				NetworkSettings: nw,
				Config: &container.Config{
					Env: []string{
						"POSTGRES_USER=loremipsum",
						"POSTGRES_DB=ipsum",
						"POSTGRES_PASSWORD=notapassword",
						"POSTGRES_PORT=1234",
					},
				},
			},
			want: "postgresql://loremipsum:notapassword@1.2.3.4:1234/ipsum",
		},
		{
			name: "no user",
			container: &types.ContainerJSON{
				NetworkSettings: nw,
				Config: &container.Config{
					Env: []string{
						"POSTGRES_DB=ipsum",
						"POSTGRES_PASSWORD=notapassword",
						"POSTGRES_PORT=1234",
					},
				},
			},
			want: "postgresql://root:notapassword@1.2.3.4:1234/ipsum",
		},
		{
			name: "no password",
			container: &types.ContainerJSON{
				NetworkSettings: nw,
				Config: &container.Config{
					Env: []string{
						"POSTGRES_USER=loremipsum",
						"POSTGRES_DB=ipsum",
						"POSTGRES_PORT=1234",
					},
				},
			},
			want: "postgresql://loremipsum:@1.2.3.4:1234/ipsum",
		},
		{
			name: "no password, but root",
			container: &types.ContainerJSON{
				NetworkSettings: nw,
				Config: &container.Config{
					Env: []string{
						"POSTGRES_USER=loremipsum",
						"POSTGRES_DB=ipsum",
						"POSTGRES_PORT=1234",
						"POSTGRES_ROOT_PASSWORD=roooot",
					},
				},
			},
			want: "postgresql://loremipsum:roooot@1.2.3.4:1234/ipsum",
		},
		{
			name: "no port",
			container: &types.ContainerJSON{
				NetworkSettings: nw,
				Config: &container.Config{
					Env: []string{
						"POSTGRES_USER=loremipsum",
						"POSTGRES_DB=ipsum",
						"POSTGRES_PASSWORD=notapassword",
					},
				},
			},
			want: "postgresql://loremipsum:notapassword@1.2.3.4:5432/ipsum",
		},
		{
			name: "no db",
			container: &types.ContainerJSON{
				NetworkSettings: nw,
				Config: &container.Config{
					Env: []string{
						"POSTGRES_USER=loremipsum",
						"POSTGRES_PASSWORD=notapassword",
						"POSTGRES_PORT=1234",
					},
				},
			},
			want: "postgresql://loremipsum:notapassword@1.2.3.4:1234/postgres",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &PostgresDumper{
				Container: tt.container,
			}
			if got := d.buildConnStr(); got != tt.want {
				t.Errorf("buildConnStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
