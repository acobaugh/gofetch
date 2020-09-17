package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/acobaugh/gofetch/pkg/transport"
	log "github.com/sirupsen/logrus"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	url := os.Args[1]

	tp := transport.NewTransport()
	c := &http.Client{Transport: tp}

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
		log.WithFields(tpFields).WithFields(respFields).Info("success")
	} else {
		log.WithFields(tpFields).WithFields(respFields).Error("NOT success")
	}

	// output := ioutil.Discard
	// if show {
	// 	output = os.Stdout
	// }
	// io.Copy(output, resp.Body)

}

func printHelp() {
	fmt.Printf("Usage: %s <url>", os.Args[0])
}
