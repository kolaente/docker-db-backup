package main

import (
	"github.com/robfig/cron/v3"
	"log"
)

func main() {
	c, err := getClient()
	if err != nil {
		log.Fatalf("Could not create client: %s", err)
	}

	cr := cron.New()
	_, err = cr.AddFunc(config.Schedule, func() {
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

		err = dumpAllDatabases(c)
		if err != nil {
			log.Fatalf("Could not dump databases: %s", err)
		}

		err = callWebhook()
		if err != nil {
			log.Fatalf("Could not call completion webhook: %s", err)
		}

		log.Println("Done.")
	})
	if err != nil {
		log.Fatalf("Could not create cron job: %s\n", err)
	}
	cr.Run()
}
