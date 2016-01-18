package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-log"
	"github.com/whosonfirst/go-whosonfirst-nearby"
	"io"
	"net/http"
	"os"
	"strconv"
)

func main() {

	var id = flag.String("id", "id", "")
	var lat = flag.String("latitude", "latitude", "")
	var lon = flag.String("longitude", "longitude", "")
	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 1414, "The port number to listen for requests on")
	var cors = flag.Bool("cors", false, "Enable CORS headers")
	var loglevel = flag.String("loglevel", "info", "Log level for reporting")

	flag.Parse()

	key := make(map[string]string)
	key["id"] = *id
	key["latitude"] = *lat
	key["longitude"] = *lon

	var logger *log.WOFLogger

	/*
		Please wrap all this logic in to SimpleWOFLogger
		See also: https://github.com/whosonfirst/go-whosonfirst-log/issues/1
	*/

	if *loglevel != "" {
		logger = log.NewWOFLogger("[wof-nearby-server]")

		stdout := io.Writer(os.Stdout)
		stderr := io.Writer(os.Stderr)

		logger.AddLogger(stdout, *loglevel)
		logger.AddLogger(stderr, "error")
	} else {
		logger = log.SimpleWOFLogger("[wof-nearby-server]")
	}

	idx := nearby.NewIndex(logger)

	for _, path := range flag.Args() {
		idx.IndexCSVFile(path, key)
	}

	handler := func(rsp http.ResponseWriter, req *http.Request) {

		query := req.URL.Query()

		str_lat := query.Get("latitude")
		str_lon := query.Get("longitude")

		if str_lat == "" {
			http.Error(rsp, "Missing latitude parameter", http.StatusBadRequest)
			return
		}

		if str_lon == "" {
			http.Error(rsp, "Missing longitude parameter", http.StatusBadRequest)
			return
		}

		lat, lat_err := strconv.ParseFloat(str_lat, 64)
		lon, lon_err := strconv.ParseFloat(str_lon, 64)

		if lat_err != nil {
			http.Error(rsp, "Invalid latitude parameter", http.StatusBadRequest)
			return
		}

		if lon_err != nil {
			http.Error(rsp, "Invalid longitude parameter", http.StatusBadRequest)
			return
		}

		if lat > 90.0 || lat < -90.0 {
			http.Error(rsp, "E_IMPOSSIBLE_LATITUDE", http.StatusBadRequest)
			return
		}

		if lon > 180.0 || lon < -180.0 {
			http.Error(rsp, "E_IMPOSSIBLE_LONGITUDE", http.StatusBadRequest)
			return
		}

		rows := 10
		dist := 1.0

		results := idx.Nearby(lat, lon, rows, dist)

		js, err := json.Marshal(results)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		// maybe this although it seems like it adds functionality for a lot of
		// features this server does not need - https://github.com/rs/cors
		// (20151022/thisisaaronland)

		if *cors {
			rsp.Header().Set("Access-Control-Allow-Origin", "*")
		}

		rsp.Header().Set("Content-Type", "application/json")
		rsp.Write(js)
	}

	endpoint := fmt.Sprintf("%s:%d", *host, *port)
	logger.Status("wof-nearby-server listening on %s", endpoint)

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(endpoint, nil)

	if err != nil {
		logger.Error("failed to start server, because %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}
