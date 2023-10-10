package vendordb

import (
	"errors"
	"net/http"
	"time"

	"log/slog"

	"crdx.org/db"
	"crdx.org/lighthouse/env"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/repos/vendorLookupR"
	"crdx.org/lighthouse/services"
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
	for _, lookup := range vendorLookupR.Unprocessed() {
		adapter, found := db.First[m.Adapter](lookup.AdapterID)
		if !found {
			lookup.Delete()
			continue
		}

		log := self.log.With("adapter_id", adapter.ID, "mac", adapter.MACAddress)

		update := func(vendor string) {
			log.Info("lookup complete", "vendor", vendor)

			var succeeded bool
			if vendor != "" {
				adapter.Update("vendor", vendor)
				succeeded = true
			}

			lookup.Update("processed", true, "succeeded", succeeded)
		}

	retry:

		res, err := self.getVendor(adapter.MACAddress)

		if err != nil || res.StatusCode == http.StatusNotFound {
			update("")
			continue
		}

		if res.StatusCode == http.StatusUnauthorized {
			log.Error("authorisation failed", "response_code", res.StatusCode, "body", res.String())
			return errors.New("authorisation failed")
		}

		if res.StatusCode == http.StatusTooManyRequests {
			delay := 5 * time.Second
			self.log.Info("throttling", "response_code", res.StatusCode, "delay", delay)
			time.Sleep(delay)

			goto retry
		}

		if res.StatusCode != http.StatusOK {
			log.Error("request failed", "response_code", res.StatusCode, "body", res.String())
			continue
		}

		vendor := res.String()

		if vendor == "" || res.StatusCode == http.StatusNotFound {
			update("")
			continue
		}

		update(vendor)

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
