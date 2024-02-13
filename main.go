package main

import (
	"github.com/robfig/cron/v3"
	"log"
	"os"
)

func runBackup() {
	c, err := getClient()
	if err != nil {
		log.Fatalf("Could not create client: %s", err)
	}

	updateFullBackupPath()

	containers, err := getContainers(c)
	if err != nil {
		log.Fatalf("Could not get containers: %s", err)
	}

	storeContainers(c, containers)

	err = cleanupOldBackups()
	if err != nil {
		log.Fatalf("Could not clean old backups: %s", err)
	}

	dumpAllDatabases(c)

	err = callWebhook()
	if err != nil {
		log.Fatalf("Could not call completion webhook: %s", err)
	}

	log.Println("Done.")
}

func main() {
	noCron, has := os.LookupEnv("BACKUP_NO_CRON")
	if has && (noCron == "true" || noCron == "1") {
		log.Println("BACKUP_NO_CRON set, running backup once, then exiting")
		runBackup()
		return
	}

	cr := cron.New()
	_, err := cr.AddFunc(config.Schedule, runBackup)
	if err != nil {
		log.Fatalf("Could not create cron job: %s\n", err)
	}
	log.Println("DB backup service started.")

	cr.Run()
}
