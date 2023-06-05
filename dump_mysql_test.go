package main

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"reflect"
	"testing"
)

func TestMysqlDumper_buildDumpArgs(t *testing.T) {
	nw := &types.NetworkSettings{
		DefaultNetworkSettings: types.DefaultNetworkSettings{
			IPAddress: "1.2.3.4",
		},
	}
	tests := []struct {
		name      string
		container *types.ContainerJSON
		want      []string
	}{
		{
			name: "values for everything",
			container: &types.ContainerJSON{
				NetworkSettings: nw,
				Config: &container.Config{
					Env: []string{
						"MYSQL_USER=loremipsum",
						"MYSQL_DATABASE=ipsum",
						"MYSQL_PASSWORD=notapassword",
						"MYSQL_PORT=1234",
					},
				},
			},
			want: []string{"--lock-tables=0", "--dump-date", "--single-transaction", "--source-data=2", "--routines", "--triggers", "--events", "-u", "loremipsum", "-p" + "notapassword", "--port", "1234", "-h", "1.2.3.4", "ipsum"},
		},
		{
			name: "no user",
			container: &types.ContainerJSON{
				NetworkSettings: nw,
				Config: &container.Config{
					Env: []string{
						"MYSQL_DATABASE=ipsum",
						"MYSQL_PASSWORD=notapassword",
						"MYSQL_PORT=1234",
					},
				},
			},
			want: []string{"--lock-tables=0", "--dump-date", "--single-transaction", "--source-data=2", "--routines", "--triggers", "--events", "-u", "root", "-p" + "notapassword", "--port", "1234", "-h", "1.2.3.4", "ipsum"},
		},
		{
			name: "no password",
			container: &types.ContainerJSON{
				NetworkSettings: nw,
				Config: &container.Config{
					Env: []string{
						"MYSQL_USER=loremipsum",
						"MYSQL_DATABASE=ipsum",
						"MYSQL_PORT=1234",
					},
				},
			},
			want: []string{"--lock-tables=0", "--dump-date", "--single-transaction", "--source-data=2", "--routines", "--triggers", "--events", "-u", "loremipsum", "--port", "1234", "-h", "1.2.3.4", "ipsum"},
		},
		{
			name: "no password, but root",
			container: &types.ContainerJSON{
				NetworkSettings: nw,
				Config: &container.Config{
					Env: []string{
						"MYSQL_USER=loremipsum",
						"MYSQL_DATABASE=ipsum",
						"MYSQL_PORT=1234",
						"MYSQL_ROOT_PASSWORD=roooot",
					},
				},
			},
			want: []string{"--lock-tables=0", "--dump-date", "--single-transaction", "--source-data=2", "--routines", "--triggers", "--events", "-u", "loremipsum", "-p" + "roooot", "--port", "1234", "-h", "1.2.3.4", "ipsum"},
		},
		{
			name: "no port",
			container: &types.ContainerJSON{
				NetworkSettings: nw,
				Config: &container.Config{
					Env: []string{
						"MYSQL_USER=loremipsum",
						"MYSQL_DATABASE=ipsum",
						"MYSQL_PASSWORD=notapassword",
					},
				},
			},
			want: []string{"--lock-tables=0", "--dump-date", "--single-transaction", "--source-data=2", "--routines", "--triggers", "--events", "-u", "loremipsum", "-p" + "notapassword", "--port", "3306", "-h", "1.2.3.4", "ipsum"},
		},
		{
			name: "no db",
			container: &types.ContainerJSON{
				NetworkSettings: nw,
				Config: &container.Config{
					Env: []string{
						"MYSQL_USER=loremipsum",
						"MYSQL_PASSWORD=notapassword",
						"MYSQL_PORT=1234",
					},
				},
			},
			want: []string{"--lock-tables=0", "--dump-date", "--single-transaction", "--source-data=2", "--routines", "--triggers", "--events", "-u", "loremipsum", "-p" + "notapassword", "--port", "1234", "-h", "1.2.3.4", "--all-databases"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MysqlDumper{
				Container: tt.container,
			}
			if got := m.buildDumpArgs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildDumpArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
