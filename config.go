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
	Interval              time.Duration
	MaxBackups            int64
}

var (
	config   *conf
	dumpTime time.Time
)

const (
	envBackupFolder = `BACKUP_FOLDER`
	envInterval     = `BACKUP_INTERVAL`
	envMax          = `BACKUP_MAX`
)

func init() {
	config = &conf{
		Folder:     "/backups/",
		Interval:   time.Hour * 6,
		MaxBackups: 24,
	}

	folder, has := os.LookupEnv(envBackupFolder)
	if has {
		if !strings.HasSuffix(folder, "/") {
			folder = folder + "/"
		}

		config.Folder = folder
	}

	var err error

	interval, has := os.LookupEnv(envInterval)
	if has {
		config.Interval, err = time.ParseDuration(interval)
		if err != nil {
			log.Fatalf("Invalid interval: %s\n", err)
		}
	}

	max, has := os.LookupEnv(envMax)
	if has {
		config.MaxBackups, err = strconv.ParseInt(max, 10, 64)
		if err != nil {
			log.Fatalf("Invalid max: %s\n", err)
		}
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
