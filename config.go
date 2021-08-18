package main

import (
	"log"
	"os"
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

func init() {
	config = &conf{
		Folder:     "/backups/",
		Interval:   time.Hour * 6,
		MaxBackups: 12,
	}

	folder, has := os.LookupEnv("BACKUP_FOLDER")
	if has {
		if !strings.HasSuffix(folder, "/") {
			folder = folder + "/"
		}

		config.Folder = folder
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
