package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-log"
	"github.com/whosonfirst/go-whosonfirst-nearby"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	var id = flag.String("id", "id", "")
	var lat = flag.String("latitude", "latitude", "")
	var lon = flag.String("longitude", "longitude", "")

	flag.Parse()

	key := make(map[string]string)
	key["id"] = *id
	key["latitude"] = *lat
	key["longitude"] = *lon

	logger := log.SimpleWOFLogger("[wof-nearby-csv]")
	idx := nearby.NewIndex(logger)

	t1 := time.Now()

	for _, path := range flag.Args() {
		idx.IndexCSVFile(path, key)
	}

	t2 := time.Since(t1)

	fmt.Printf("# time to index %v\n", t2)

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("# query <lat>,<lon>")
	fmt.Println("")

	for scanner.Scan() {

		input := scanner.Text()
		query := strings.Split(input, ",")

		if len(query) != 2 {
			fmt.Println("invalid query")
			continue
		}

		str_lat := query[0]
		str_lon := query[1]

		lat, _ := strconv.ParseFloat(str_lat, 64)
		lon, _ := strconv.ParseFloat(str_lon, 64)

		t1 := time.Now()

		points := idx.Nearby(lat, lon, 10, 1)
		t2 := time.Since(t1)

		count := len(points)

		fmt.Printf("# time to retrieve %d records: %v\n", count, t2)

		for _, pt := range points {
			fmt.Println(pt)
		}

		fmt.Println("")
		fmt.Println("# query <lat>,<lon>")
		fmt.Println("")

	}
}
