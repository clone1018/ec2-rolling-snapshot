package main

import (
	"gopkg.in/gcfg.v1"
	"io/ioutil"
	"os"
)

var (
	DefaultConfigFile = "config.gcfg"
)

const (
	// The example file kept in version control. We'll copy and load from this
	// by default.
	CONFIG_EXAMPLE = `; ec2-rolling-snapshot
; You can create as many of these as you want
[snapshot-task "some-easy-to-understand-name"]
volume = vol-1234567
snapshots = 6 
region = us-east-1
`
)

type Configuration struct {
	Snapshot_Task map[string]*SnapshotTask
}

// Reads the configuration from the config file, copying a config into
// place from the example if one does not yet exist.
func (c *Configuration) load(file string) error {
	err := c.ensureConfigExists(file)
	if err != nil {
		return err
	}

	return gcfg.ReadFileInto(c, file)
}

// Creates the config.gcfg if it does not exist.
func (c *Configuration) ensureConfigExists(file string) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return ioutil.WriteFile(file, []byte(CONFIG_EXAMPLE), 0644)
	} else {
		return nil
	}
}
