package main

import (
	"log"
	"time"
)

func main() {
	c, err := getClient()
	if err != nil {
		log.Fatalf("Could not create client: %s", err)
	}

	for {
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

		log.Println("Done.")
		log.Printf("Sleeping for %s\n", config.Interval)

		time.Sleep(config.Interval)
	}
}
