package nearby

import (
	"errors"
	"fmt"
	"github.com/hailocab/go-geoindex"
	"github.com/whosonfirst/go-whosonfirst-csv"
	"github.com/whosonfirst/go-whosonfirst-geojson"
	"io"
	"strconv"
)

type Callback func(p geoindex.Point) bool

func NewRecordFromFeature(f *geojson.WOFFeature) (*Record, error) {

	id := f.Id()

	body := f.Body()

	var lat float64
	var lon float64
	var ok bool

	/*
	   TODO: check for lbl: (and other) centroids
	*/

	lat, ok = body.Path("properties.geom:latitude").Data().(float64)

	if !ok {
		return nil, errors.New("failed to determine geom:latitude")
	}

	lon, ok = body.Path("properties.geom:longitude").Data().(float64)

	if !ok {
		return nil, errors.New("failed to determine geom:longitude")
	}

	str_id := strconv.Itoa(id)

	return NewRecord(str_id, lat, lon)
}

func NewRecord(id string, lat float64, lon float64) (*Record, error) {

	record := Record{id, lat, lon}
	return &record, nil
}

type Record struct {
	id  string
	lat float64
	lon float64
}

func (r *Record) Lat() float64 {
	return r.lat
}

func (r *Record) Lon() float64 {
	return r.lon
}

func (r *Record) Id() string {
	return r.id
}

func (r *Record) Stringer() string {
	return fmt.Sprintf("%s %06f %06f", r.Id(), r.Lat(), r.Lon())
}

func NewIndex() *Index {

	geoindex := geoindex.NewPointsIndex(geoindex.Km(0.5))
	return &Index{geoindex}
}

type Index struct {
	geoindex *geoindex.PointsIndex
}

func (i *Index) IndexCSVFile(csv_file string, key map[string]string) (bool, error) {

	reader, reader_err := csv.NewDictReader(csv_file)

	if reader_err != nil {
		return false, reader_err
	}

	for {
		row, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			break
		}

		id, ok := row[key["id"]]

		if !ok {
			continue
		}

		str_lat, ok := row[key["latitude"]]

		if !ok {
			continue
		}

		str_lon, ok := row[key["longitude"]]

		if !ok {
			continue
		}

		lat, _ := strconv.ParseFloat(str_lat, 64)
		lon, _ := strconv.ParseFloat(str_lon, 64)

		record, err := NewRecord(id, lat, lon)

		if err != nil {
			continue
		}

		i.geoindex.Add(record)
	}

	return true, nil
}

func (i *Index) IndexFeature(f *geojson.WOFFeature) (bool, error) {

	record, err := NewRecordFromFeature(f)

	if err != nil {
		return false, err
	}

	i.geoindex.Add(record)
	return true, nil
}

func (i *Index) Nearby(lat float64, lon float64) []geoindex.Point {

	id := "u r on first"
	pt := &geoindex.GeoPoint{id, lat, lon}

	cb := func(p geoindex.Point) bool {
		return true
	}

	return i.nearby(pt, cb)
}

func (i *Index) NearbyWithCallback(lat float64, lon float64, cb Callback) []geoindex.Point {

	id := "u r on first"
	pt := &geoindex.GeoPoint{id, lat, lon}

	return i.nearby(pt, cb)
}

func (i *Index) nearby(pt *geoindex.GeoPoint, cb Callback) []geoindex.Point {

	points := i.geoindex.KNearest(pt, 5, geoindex.Km(5), cb)
	return points
}
