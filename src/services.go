package main

import (
	"time"

	"crdx.org/lighthouse/services"
	"crdx.org/lighthouse/services/notifier"
	"crdx.org/lighthouse/services/scanner"
	"crdx.org/lighthouse/services/vendordb"
	"crdx.org/lighthouse/services/watcher"
)

func startServices() {
	services.Start("scanner", &services.Config{
		Service:     scanner.New(),
		RunInterval: 0,
	})

	services.Start("vendordb", &services.Config{
		Service:     vendordb.New(),
		RunInterval: 5 * time.Second,
	})

	watcherStartDelay := 30 * time.Second

	services.Start("watcher", &services.Config{
		Service:     watcher.New(),
		RunInterval: 10 * time.Second,
		StartDelay:  watcherStartDelay,
	})

	services.Start("notifier", &services.Config{
		Service:     notifier.New(),
		RunInterval: 1 * time.Minute,

		// Give the watcher time to establish the new state of devices.
		StartDelay: watcherStartDelay * 2,

		// A fairly high initial restart interval should be used, as we don't want to risk spamming
		// notifications.
		InitialRestartInterval: 1 * time.Minute,
	})
}
