package prober

import (
	"log/slog"
	"slices"

	"crdx.org/lighthouse/cmd/lighthouse/services"
	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/db/repo/probeR"
	"crdx.org/lighthouse/db/repo/settingR"
	"crdx.org/lighthouse/pkg/probe"
	"github.com/samber/lo"
	"github.com/samber/lo/mutable"
)

type Prober struct {
	logger *slog.Logger
}

func New() *Prober {
	return &Prober{}
}

func (self *Prober) Init(args *services.Args) error {
	self.logger = args.Logger
	return nil
}

func (self *Prober) Run() error {
	devices := db.FindScannableDevices()
	mutable.Shuffle(devices)

	scan := db.CreateScan(&db.Scan{})

	for _, device := range devices {
		for _, adapter := range device.Adapters() {
			if !db.IsOnline(device, adapter) {
				continue
			}

			ports := lo.Must(probe.Scan(adapter.MACAddress, adapter.IPAddress)).Ports

			for _, port := range ports {
				self.logger.Info(
					"found service",
					"device", device.DisplayName(),
					"name", probe.ServiceName(port),
					"port", port,
					"ip", adapter.IPAddress,
				)

				db.CreateScanResult(&db.ScanResult{
					ScanID:   scan.ID,
					DeviceID: device.ID,
					Port:     port,
				})
			}

			process(device, ports)
			break
		}
	}

	scan.UpdateCompletedAt(db.Now())
	return nil
}

func process(device *db.Device, foundPorts []int64) {
	services := db.FindServicesByDeviceID(device.ID)

	var servicePorts []int64
	for _, service := range services {
		if !slices.Contains(foundPorts, service.Port) {
			if service.LastSeenAt.Before(db.Now().Add(-probeR.TTL())) {
				service.Delete()
			}
		} else {
			service.UpdateLastSeenAt(db.Now())
		}

		servicePorts = append(servicePorts, service.Port)
	}

	for _, port := range foundPorts {
		if !slices.Contains(servicePorts, port) {
			service := db.CreateService(&db.Service{
				DeviceID:   device.ID,
				Port:       port,
				LastSeenAt: db.Now(),
			})

			if settingR.NotifyOnNewService() {
				previouslyFound := db.CountPreviousServices(device.ID, port, service.ID) > 0

				if !previouslyFound {
					db.CreateDeviceServiceNotification(&db.DeviceServiceNotification{
						DeviceID:  device.ID,
						ServiceID: db.N(service.ID),
					})
				}
			}
		}
	}
}
