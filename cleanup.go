package main

import (
	"io/ioutil"
	"log"
	"os"
	"sort"
)

func cleanupOldBackups() error {
	files, err := ioutil.ReadDir(config.Folder)
	if err != nil {
		return err
	}

	if len(files) < config.MaxBackups {
		return nil
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Unix() < files[j].ModTime().Unix()
	})

	oldest := files[:len(files)-config.MaxBackups]

	for _, file := range oldest {
		log.Printf("Removing old backup folder %s...\n", file.Name())
		err = os.RemoveAll(config.Folder + file.Name())
		if err != nil {
			return err
		}
	}

	return nil
}
