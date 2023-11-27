package prober

import (
	"log/slog"
	"slices"
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/deviceR"
	"crdx.org/lighthouse/m/repo/probeR"
	"crdx.org/lighthouse/m/repo/settingR"
	"crdx.org/lighthouse/pkg/probe"
	"crdx.org/lighthouse/services"
	"github.com/samber/lo"
)

type Prober struct {
	log *slog.Logger
}

func New() *Prober {
	return &Prober{}
}

func (self *Prober) Init(args *services.Args) error {
	self.log = args.Logger
	return nil
}

func (self *Prober) Run() error {
	for {
		if settingR.EnableServiceScan() {
			self.run()
		}

		time.Sleep(settingR.ServiceScanInterval())
	}
}

func (self *Prober) run() {
	devices := deviceR.Scannable()
	lo.Shuffle(devices)

	scan := db.Save(&m.Scan{})

	self.log.Info("service scan started")

	for _, device := range devices {
		for _, adapter := range device.Adapters() {
			if !adapter.IsOnline() {
				continue
			}

			ports := lo.Must(probe.Scan(adapter.MACAddress, adapter.IPAddress)).Ports

			for _, port := range ports {
				self.log.Info(
					"found service",
					"device", device.DisplayName(),
					"name", probe.ServiceName(port),
					"port", port,
					"ip", adapter.IPAddress,
				)

				db.Save(&m.ScanResult{
					ScanID:   scan.ID,
					DeviceID: device.ID,
					Port:     port,
				})
			}

			process(device, ports)
			break
		}
	}

	scan.Update("completed_at", time.Now())
	self.log.Info("service scan completed")
}

func process(device *m.Device, foundPorts []uint) {
	services := db.B[m.Service]("device_id = ?", device.ID).Find()

	var servicePorts []uint
	for _, service := range services {
		if !slices.Contains(foundPorts, service.Port) {
			if service.UpdatedAt.Before(time.Now().Add(-probeR.TTL())) {
				service.Delete()
			}
		} else {
			service.Update("last_seen_at", time.Now())
		}

		servicePorts = append(servicePorts, service.Port)
	}

	for _, port := range foundPorts {
		if !slices.Contains(servicePorts, port) {
			service := db.Create(&m.Service{
				DeviceID:   device.ID,
				Port:       port,
				LastSeenAt: time.Now(),
			})

			if settingR.NotifyOnNewService() {
				previouslyFound := db.B[m.Service](
					"device_id = ? and port = ? and id != ?",
					device.ID,
					port,
					service.ID,
				).Unscoped().Exists()

				if !previouslyFound {
					db.Create(&m.DeviceServiceNotification{
						DeviceID:  device.ID,
						ServiceID: service.ID,
					})
				}
			}
		}
	}
}
