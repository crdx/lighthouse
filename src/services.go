package main

import (
	"time"

	"crdx.org/lighthouse/m/repo/settingR"
	"crdx.org/lighthouse/services"
	"crdx.org/lighthouse/services/notifier"
	"crdx.org/lighthouse/services/reader"
	"crdx.org/lighthouse/services/vendor"
	"crdx.org/lighthouse/services/watcher"
	"crdx.org/lighthouse/services/writer"
)

func startServices() {
	services.Start("reader", &services.Config{
		Service:     reader.New(),
		RunInterval: 0,
	})

	services.Start("writer", &services.Config{
		Service:     writer.New(),
		RunInterval: settingR.ScanInterval(),
	})

	services.Start("vendor", &services.Config{
		Service:     vendor.New(),
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
