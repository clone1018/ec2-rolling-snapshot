package main

import (
	"flag"
	"log"
)

var (
	config     = &Configuration{}
	configFile = flag.String("config", DefaultConfigFile, "specify a config file, it will be created if not existing")
)

func main() {
	flag.Parse()

	err := config.load(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	for taskName, task := range config.Snapshot_Task {
		task.CreateSvc()

		task.CreateSnapshot(taskName)

		_, err := task.DeleteOldSnapshots(taskName)
		if err != nil {
			log.Fatal(err)
		}

	}
}
