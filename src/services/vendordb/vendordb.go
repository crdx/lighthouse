package vendordb

import (
	"errors"
	"net/http"
	"time"

	"log/slog"

	"crdx.org/db"
	"crdx.org/lighthouse/constants"
	"crdx.org/lighthouse/env"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/services"
	"crdx.org/lighthouse/util"
	"github.com/imroc/req/v3"
)

type VendorDB struct {
	log *slog.Logger
}

func New() *VendorDB {
	return &VendorDB{}
}

func (self *VendorDB) Init(args *services.Args) error {
	self.log = args.Logger

	if env.MACVendorsAPIKey == "" {
		return errors.New("missing API key")
	}

	return nil
}

func (self *VendorDB) Run() error {
	for _, device := range db.B[m.Device]().Where(`vendor = ""`).Find() {
		log := self.log.With("mac", device.MACAddress)

		update := func(deviceID uint, vendor string) {
			log.Info("lookup complete", "vendor", vendor)

			columns := db.Map{}
			columns["vendor"] = vendor

			if device.Name == "" || device.Name == constants.UnknownDeviceLabel {
				columns["name"] = vendor
			}

			db.B(m.Device{ID: deviceID}).Update(columns)
		}

	retry:

		res, err := self.getVendor(device.MACAddress)

		if err != nil || res.StatusCode == http.StatusNotFound {
			update(device.ID, constants.UnknownVendorLabel)
			continue
		}

		if res.StatusCode == http.StatusUnauthorized {
			log.Error("authorisation failed", "response_code", res.StatusCode, "body", res.String())
			return errors.New("authorisation failed")
		}

		if res.StatusCode == http.StatusTooManyRequests {
			delay := 5 * time.Second
			self.log.Info("throttling", "response_code", res.StatusCode, "delay", delay)
			util.Sleep(delay)

			goto retry
		}

		if res.StatusCode != http.StatusOK {
			log.Error("request failed", "response_code", res.StatusCode, "body", res.String())
			continue
		}

		vendor := res.String()

		if vendor == "" || res.StatusCode == http.StatusNotFound {
			update(device.ID, constants.UnknownVendorLabel)
			continue
		}

		update(device.ID, vendor)

		// MacVendors API free plan allows 2 requests per second, so to be safe limit to 1 per second.
		time.Sleep(time.Second)
	}

	return nil
}

func (*VendorDB) getVendor(macAddress string) (*req.Response, error) {
	return req.C().R().
		SetBearerAuthToken(env.MACVendorsAPIKey).
		SetHeader("Accept", "text/plain").
		Get("https://api.macvendors.com/v1/lookup/" + macAddress)
}
