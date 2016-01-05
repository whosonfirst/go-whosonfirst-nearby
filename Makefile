prep:
	if test -d pkg; then rm -rf pkg; fi

self:	prep
	if test -d src/github.com/whosonfirst/go-whosonfirst-nearby; then rm -rf src/github.com/whosonfirst/go-whosonfirst-nearby; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-nearby
	cp nearby.go src/github.com/whosonfirst/go-whosonfirst-nearby/

deps: 	self
	@GOPATH=$(shell pwd) go get -u "github.com/whosonfirst/go-whosonfirst-crawl"
	@GOPATH=$(shell pwd) go get -u "github.com/whosonfirst/go-whosonfirst-log"
	@GOPATH=$(shell pwd) go get -u "github.com/whosonfirst/go-whosonfirst-geojson"
	@GOPATH=$(shell pwd) go get -u "github.com/hailocab/go-geoindex"

fmt:
	go fmt *.go