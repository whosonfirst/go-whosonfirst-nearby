package nearby

import (
       "errors"
       "github.com/whosonfirst/go-whosonfirst-geojson"
       "github.com/hailocab/go-geoindex"
       "strconv"
)

func NewRecord(f *geojson.WOFFeature) (*Record, error) {

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

     record := Record{id, lat, lon}
     return &record, nil
}

type Record struct {
     id int
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
     return strconv.Itoa(r.id)
     //    return strconv.FormatInt(r.id, 10)
}

func NewIndex() *Index {

     geoindex := geoindex.NewPointsIndex(geoindex.Km(0.5))
     return &Index{geoindex}
}

type Index struct {
     geoindex *geoindex.PointsIndex
}

/*

func (i *Index) IndexMetaFile(meta string) (bool, error) {
     // please write me
}

*/

func (i *Index) IndexFeature(f *geojson.WOFFeature) (bool, error) {

     record, err := NewRecord(f)

     if err != nil {
     	return false, err
     }

     i.geoindex.Add(record)
     return true, nil
}
