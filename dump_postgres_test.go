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

func TestFindPGVersionFromEnv(t *testing.T) {
	t.Run("no PG_MAJOR", func(t *testing.T) {
		pgVersion := findPgVersion([]string{})
		if pgVersion != "" {
			t.Errorf("Version is not empty")
		}
	})
	t.Run("pg 14", func(t *testing.T) {
		pgVersion := findPgVersion([]string{
			"POSTGRES_PASSWORD=test",
			"POSTGRES_USER=user",
			"POSTGRES_DB=test",
			"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/lib/postgresql/14/bin",
			"GOSU_VERSION=1.14",
			"LANG=en_US.utf8",
			"PG_MAJOR=14",
			"PG_VERSION=14.1-1.pgdg110+1",
			"PGDATA=/var/lib/postgresql/data",
		})
		if pgVersion != "14" {
			t.Errorf("Version is not 14")
		}
	})
}
