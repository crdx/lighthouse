package vendor

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"crdx.org/lighthouse/cmd/lighthouse/services"
	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/db/repo/settingR"
	"github.com/imroc/req/v3"
)

type Vendor struct {
	logger *slog.Logger
}

func New() *Vendor {
	return &Vendor{}
}

const backoff = 10 * time.Second

func (self *Vendor) Init(args *services.Args) error {
	self.logger = args.Logger
	return nil
}

func (self *Vendor) Run() error {
	for _, lookup := range db.FindVendorLookupsByProcessed(false) {
		adapter, found := db.FindAdapter(lookup.AdapterID)
		if !found {
			lookup.Delete()
			continue
		}

		logger := self.logger.With(slog.Group("adapter", "id", adapter.ID, "mac", adapter.MACAddress))

		if err := doLookup(lookup, adapter, logger); err != nil {
			return err
		}
	}

	return nil
}

func doLookup(lookup *db.VendorLookup, adapter *db.Adapter, log *slog.Logger) error {
	// MacVendors API free plan allows 2 requests per second, so to be safe limit to 1 per second.
	defer time.Sleep(time.Second)

	update := func(vendor string) {
		log.Info("lookup complete", "vendor", vendor)

		var succeeded bool
		if vendor != "" {
			adapter.UpdateVendor(vendor)
			succeeded = true
		}

		lookup.UpdateProcessed(true)
		lookup.UpdateSucceeded(succeeded)
	}

retry:

	res, err := getVendor(adapter.MACAddress)

	if err != nil || res.StatusCode == http.StatusNotFound {
		update("")
		return nil //nolint
	}

	if res.StatusCode == http.StatusUnauthorized {
		log.Error("authorisation failed", "response_code", res.StatusCode, "body", res.String())
		return errors.New("authorisation failed")
	}

	if res.StatusCode == http.StatusTooManyRequests {
		log.Info("throttling", "response_code", res.StatusCode, "delay", backoff)
		time.Sleep(backoff)
		goto retry
	}

	if res.StatusCode != http.StatusOK {
		log.Error("request failed", "response_code", res.StatusCode, "body", res.String())
		// A non-OK response could be an intermittent networking failure, or perhaps the API is
		// down. Panic as we don't know which.
		panic("request failed")
	}

	vendor := res.String()

	if vendor == "" || res.StatusCode == http.StatusNotFound {
		update("")
		return nil
	}

	update(vendor)
	return nil
}

func getVendor(macAddress string) (*req.Response, error) {
	return req.R().
		SetBearerAuthToken(settingR.MACVendorsAPIKey()).
		SetHeader("Accept", "text/plain").
		Get("https://api.macvendors.com/v1/lookup/" + macAddress)
}
