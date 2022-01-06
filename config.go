package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Backup folder
// Backup interval
// Max backups to keep

type conf struct {
	Folder                string // Backup folder _with_ trailing slash
	fullCurrentBackupPath string
	Schedule              string
	MaxBackups            int
}

var (
	config   *conf
	dumpTime time.Time
)

const (
	envBackupFolder = `BACKUP_FOLDER`
	envSchedule     = `BACKUP_SCHEDULE`
	envMax          = `BACKUP_MAX`
)

func init() {
	config = &conf{
		Folder:     "/backups/",
		Schedule:   "0 */6 * * *",
		MaxBackups: 12,
	}

	folder, has := os.LookupEnv(envBackupFolder)
	if has {
		if !strings.HasSuffix(folder, "/") {
			folder = folder + "/"
		}

		config.Folder = folder
	}

	schedule, has := os.LookupEnv(envSchedule)
	if has {
		config.Schedule = schedule
	}

	max, has := os.LookupEnv(envMax)
	if has {
		maxBackups, err := strconv.ParseInt(max, 10, 32)
		if err != nil {
			log.Fatalf("Invalid max: %s\n", err)
		}
		config.MaxBackups = int(maxBackups)
	}

	updateFullBackupPath()
}

func updateFullBackupPath() {
	dumpTime = time.Now()
	config.fullCurrentBackupPath = config.Folder + dumpTime.Format("02-01-2006_15-04-05") + "/"
	err := os.MkdirAll(config.fullCurrentBackupPath, 0744)
	if err != nil {
		log.Fatalf("Could not create backup folder: %s\n", err)
	}
}
