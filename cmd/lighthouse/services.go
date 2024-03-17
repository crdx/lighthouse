package main

import (
	"time"

	"crdx.org/lighthouse/cmd/lighthouse/services"
	"crdx.org/lighthouse/cmd/lighthouse/services/notifier"
	"crdx.org/lighthouse/cmd/lighthouse/services/pinger"
	"crdx.org/lighthouse/cmd/lighthouse/services/prober"
	"crdx.org/lighthouse/cmd/lighthouse/services/reader"
	"crdx.org/lighthouse/cmd/lighthouse/services/vendor"
	"crdx.org/lighthouse/cmd/lighthouse/services/watcher"
	"crdx.org/lighthouse/cmd/lighthouse/services/writer"
	"crdx.org/lighthouse/db/repo/settingR"
)

func startServices() {
	services.Start("writer", &services.Config{
		Service:         writer.New(),
		RunIntervalFunc: settingR.DeviceScanInterval,
		Enabled:         settingR.EnableDeviceScan,
	})

	services.Start("reader", &services.Config{
		Service: reader.New(), // Reader service runs forever.
	})

	services.Start("vendor", &services.Config{
		Service:     vendor.New(),
		RunInterval: 5 * time.Second,
		Enabled:     func() bool { return settingR.MACVendorsAPIKey() != "" },
		Quiet:       true,
	})

	watcherStartDelay := 30 * time.Second

	services.Start("watcher", &services.Config{
		Service:     watcher.New(),
		RunInterval: 10 * time.Second,
		StartDelay:  watcherStartDelay,
		Quiet:       true,
	})

	services.Start("pinger", &services.Config{
		Service:         pinger.New(),
		RunIntervalFunc: settingR.DeviceScanInterval,
		Enabled:         settingR.EnableDeviceScan,

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

		RunIntervalFunc: settingR.ServiceScanInterval,
		Enabled:         settingR.EnableServiceScan,
	})
}
