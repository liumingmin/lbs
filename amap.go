package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/liumingmin/goutils/log"
)

type AMap struct {
	Key string
}

func (t *AMap) GeoByAddresses(ctx context.Context, addresses []string, city string) []string {
	uStr := "https://restapi.amap.com/v3/geocode/geo?parameters"
	values := url.Values{}
	values.Add("key", t.Key)
	values.Add("address", strings.Join(addresses, "|"))
	values.Add("city", city)
	values.Add("batch", "true")

	var result *GeocodesResult
	err := t.getJson(ctx, uStr, values, &result)
	if err != nil || result == nil && result.Status != 1 {
		return []string{}
	}

	locations := make([]string, 0, len(result.Geocodes))
	for _, geoCode := range result.Geocodes {
		locations = append(locations, geoCode.Location)
	}
	return locations
}

func (t *AMap) getJson(ctx context.Context, url string, values url.Values, result interface{}) (err error) {
	url += "&" + values.Encode()
	log.Debug(ctx, "post json to url %s, values: %+v", url, values)
	r, err := http.NewRequest("GET", url, bytes.NewReader([]byte{}))
	if err != nil {
		log.Error(ctx, "Http create new request failed. url:%s, error: %v", url, err)
		return
	}

	client := http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		log.Error(ctx, "Http client.do get response failed. error: %v", err)
		return
	}

	defer resp.Body.Close()
	bd, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(ctx, "Read all in response failed. error: %v", err)
		return
	}

	log.Debug(ctx, "post json to url %s, result: %+v", url, string(bd))
	err = json.Unmarshal(bd, &result)
	if err != nil {
		log.Error(ctx, "Response unmarshal failed. error: %v", err)
		return
	}

	return
}

type GeocodesResult struct {
	Count    int64      `json:"count,string"` // 1
	Geocodes []*Geocode `json:"geocodes"`
	Info     string     `json:"info"`            // OK
	Infocode int64      `json:"infocode,string"` // 10000
	Status   int64      `json:"status,string"`   // 1
}

type Geocode struct {
	City     string `json:"city"`     // 北京市
	Location string `json:"location"` // 116.482086,39.990496
}
