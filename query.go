package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Query struct {
	Service string   // where to request data...
	Limit   int      // maxiumum features to request, prior to filtering
	Filters []string // extra filtering to apply
	SortBy  string   // whether to sort
}

func TimeOffset(from time.Time, offset time.Duration) string {
	return from.UTC().Add(-offset).Format("2006-01-02T15:04:05")
}

func TimeOffsetNow(offset time.Duration) string {
	return TimeOffset(time.Now(), offset)
}

func (q *Query) Values() *url.Values {

	v := url.Values{}
	v.Set("service", "WFS")
	v.Set("version", "1.0.0")
	v.Set("request", "GetFeature")
	v.Set("typeName", "geonet:quake_search_v1")
	if q.Limit > 0 {
		v.Set("maxFeatures", strconv.Itoa(q.Limit))
	}
	v.Set("outputFormat", "json")

	return &v
}

func (q *Query) Filter(filter string) {
	q.Filters = append(q.Filters, filter)
}

func (q *Query) AddFilter(k, o, v string) {
	q.Filter(k + "+" + o + "+" + v)
}

func (q *Query) URL() *url.URL {

	v := q.Values()

	var cql_filter string
	if len(q.Filters) > 0 {
		cql_filter = "&cql_filter=" + strings.Join(q.Filters, "+and+")
	}

	var sortBy string
	if len(q.Filters) > 0 && len(q.SortBy) > 0 {
		sortBy = "+and+sortBy+" + q.SortBy
	}
	u := url.URL{
		Scheme:   "http",
		Host:     q.Service,
		Path:     "geonet/ows",
		RawQuery: v.Encode() + cql_filter + sortBy,
	}

	return &u
}

func (q *Query) Get() ([]byte, error) {
	res, err := http.Get(q.URL().String())
	if err != nil {
		return nil, err
	}
	var body []byte
	if res.Body != nil {
		defer res.Body.Close()
		if body, err = ioutil.ReadAll(res.Body); err != nil {
			return nil, err
		}
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error [%d]: %s", res.StatusCode, string(body))
	}

	return body, nil
}

func (q *Query) Search() (*Search, error) {
	s := Search{}

	out, err := q.Get()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(out, &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
