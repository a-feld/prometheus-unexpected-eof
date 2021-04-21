package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

// MetricFamiliesByName is a map of Prometheus metrics family names and their
// representation.
type MetricFamiliesByName map[string]dto.MetricFamily

// HTTPDoer executes http requests. It is implemented by *http.Client.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Get scrapes the given URL and decodes the retrieved payload.
func Get(client HTTPDoer, url string) (MetricFamiliesByName, error) {
	mfs := MetricFamiliesByName{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return mfs, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return mfs, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		return nil, fmt.Errorf("status code returned by the prometheus exporter indicates an error occurred: %d", resp.StatusCode)
	}

	d := expfmt.NewDecoder(resp.Body, expfmt.FmtText)
	for {
		var mf dto.MetricFamily
		if err := d.Decode(&mf); err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}
		mfs[mf.GetName()] = mf
	}

	return mfs, nil
}

func main() {
	url := os.Args[1]
	for {
		client := &http.Client{
			Timeout: 5 * time.Second,
		}
		result, err := Get(client, url)
		if err != nil {
			fmt.Println(err)
			return
		}
		b, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(b))
		fmt.Println("-------------------------------------------------")
		time.Sleep(5 * time.Second)
	}
}
