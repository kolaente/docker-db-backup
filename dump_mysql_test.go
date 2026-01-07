package main

import (
	"reflect"
	"testing"

	"github.com/docker/docker/api/types/container"
)

func TestMysqlDumper_buildDumpArgs(t *testing.T) {
	nw := &container.NetworkSettings{
		DefaultNetworkSettings: container.DefaultNetworkSettings{
			IPAddress: "1.2.3.4",
		},
	}
	tests := []struct {
		name string
		ctr  *container.InspectResponse
		want []string
	}{
		{
			name: "values for everything",
			ctr: &container.InspectResponse{
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
			want: []string{"--lock-tables=0", "--dump-date", "--single-transaction", "--routines", "--triggers", "--events", "-u", "loremipsum", "-p" + "notapassword", "--port", "1234", "-h", "1.2.3.4", "ipsum"},
		},
		{
			name: "no user",
			ctr: &container.InspectResponse{
				NetworkSettings: nw,
				Config: &container.Config{
					Env: []string{
						"MYSQL_DATABASE=ipsum",
						"MYSQL_PASSWORD=notapassword",
						"MYSQL_PORT=1234",
					},
				},
			},
			want: []string{"--lock-tables=0", "--dump-date", "--single-transaction", "--routines", "--triggers", "--events", "-u", "root", "-p" + "notapassword", "--port", "1234", "-h", "1.2.3.4", "ipsum"},
		},
		{
			name: "no password",
			ctr: &container.InspectResponse{
				NetworkSettings: nw,
				Config: &container.Config{
					Env: []string{
						"MYSQL_USER=loremipsum",
						"MYSQL_DATABASE=ipsum",
						"MYSQL_PORT=1234",
					},
				},
			},
			want: []string{"--lock-tables=0", "--dump-date", "--single-transaction", "--routines", "--triggers", "--events", "-u", "loremipsum", "--port", "1234", "-h", "1.2.3.4", "ipsum"},
		},
		{
			name: "no password, but root",
			ctr: &container.InspectResponse{
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
			want: []string{"--lock-tables=0", "--dump-date", "--single-transaction", "--routines", "--triggers", "--events", "-u", "loremipsum", "-p" + "roooot", "--port", "1234", "-h", "1.2.3.4", "ipsum"},
		},
		{
			name: "no port",
			ctr: &container.InspectResponse{
				NetworkSettings: nw,
				Config: &container.Config{
					Env: []string{
						"MYSQL_USER=loremipsum",
						"MYSQL_DATABASE=ipsum",
						"MYSQL_PASSWORD=notapassword",
					},
				},
			},
			want: []string{"--lock-tables=0", "--dump-date", "--single-transaction", "--routines", "--triggers", "--events", "-u", "loremipsum", "-p" + "notapassword", "--port", "3306", "-h", "1.2.3.4", "ipsum"},
		},
		{
			name: "no db",
			ctr: &container.InspectResponse{
				NetworkSettings: nw,
				Config: &container.Config{
					Env: []string{
						"MYSQL_USER=loremipsum",
						"MYSQL_PASSWORD=notapassword",
						"MYSQL_PORT=1234",
					},
				},
			},
			want: []string{"--lock-tables=0", "--dump-date", "--single-transaction", "--routines", "--triggers", "--events", "-u", "loremipsum", "-p" + "notapassword", "--port", "1234", "-h", "1.2.3.4", "--all-databases"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MysqlDumper{
				Container: tt.ctr,
			}
			if got := m.buildDumpArgs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildDumpArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
