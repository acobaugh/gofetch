package main

import (
	"crypto/tls"
	"fmt"
	"github.com/acobaugh/gofetch/pkg/transport"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	flag.Usage = usage

	fQuiet := flag.BoolP("quiet", "q", false, "Suppress non-errors")
	fSkipVerify := flag.BoolP("skip-verify", "k", false, "Disable TLS verification")
	fTimeout := flag.DurationP("timeout", "t", 30*time.Second, "Total request duration timeout")
	flag.Parse()

	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}
	url := flag.Arg(0)

	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	tlsConfig := tls.Config{
		InsecureSkipVerify: *fSkipVerify,
	}

	tp := transport.NewTransport(&tlsConfig)
	c := &http.Client{
		Transport: tp,
		Timeout:   *fTimeout,
	}

	resp, err := c.Get(url)

	tpFields := log.Fields{
		"duration":     tp.Duration(),
		"reqDuration":  tp.ReqDuration(),
		"connDuration": tp.ConnDuration(),
	}

	if err != nil {
		log.WithError(err).WithFields(tpFields).Fatalf("GET error")
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	respFields := log.Fields{
		"content-length": resp.ContentLength,
		"status":         resp.StatusCode,
		"encoding":       resp.TransferEncoding,
		"body-size":      len(b),
	}

	if err != nil {
		log.WithFields(tpFields).WithFields(respFields).WithFields(log.Fields{"status": resp.StatusCode}).WithError(err).Error("error reading body")
	}

	if resp.StatusCode == http.StatusOK {
		if !*fQuiet {
			log.WithFields(tpFields).WithFields(respFields).Info("success")
		}
	} else {
		log.WithFields(tpFields).WithFields(respFields).Error("non-200 status")
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <url>\n", os.Args[0])
	flag.PrintDefaults()
}
