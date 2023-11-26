package main

import (
	"time"

	"crdx.org/lighthouse/services"
	"crdx.org/lighthouse/services/notifier"
	"crdx.org/lighthouse/services/pinger"
	"crdx.org/lighthouse/services/prober"
	"crdx.org/lighthouse/services/reader"
	"crdx.org/lighthouse/services/vendor"
	"crdx.org/lighthouse/services/watcher"
	"crdx.org/lighthouse/services/writer"
)

func startServices() {
	services.Start("writer", &services.Config{
		Service: writer.New(),
	})

	services.Start("reader", &services.Config{
		Service: reader.New(),
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

	services.Start("pinger", &services.Config{
		Service: pinger.New(),
		// Give it time to establish the current state of devices.
		StartDelay: watcherStartDelay + (30 * time.Second),
	})

	services.Start("notifier", &services.Config{
		Service:     notifier.New(),
		RunInterval: 1 * time.Minute,
		// Give it time to establish the current state of devices.
		StartDelay: watcherStartDelay + (30 * time.Second),
		// A fairly high initial restart interval as we don't want to risk spamming notifications.
		InitialRestartInterval: 1 * time.Minute,
	})

	services.Start("prober", &services.Config{
		Service:    prober.New(),
		StartDelay: 5 * time.Minute,
	})
}
