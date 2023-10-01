package main

import (
	"time"

	"crdx.org/lighthouse/services"
	"crdx.org/lighthouse/services/scanner"
	"crdx.org/lighthouse/services/vendordb"
	"crdx.org/lighthouse/services/watcher"
)

func startServices() {
	go services.Start("scanner", services.Config{
		Service:     scanner.New(),
		RunInterval: 0,
	})

	go services.Start("vendordb", services.Config{
		Service:     vendordb.New(),
		RunInterval: 5 * time.Second,
	})

	go services.Start("watcher", services.Config{
		Service:     watcher.New(),
		RunInterval: 10 * time.Second,
		StartDelay:  1 * time.Minute,
	})

	// go services.Start("prober", services.Config{
	// 	Service:     prober.New(),
	// 	RunInterval: 1 * time.Hour,
	// 	StartDelay:  1 * time.Hour,
	// })
}
