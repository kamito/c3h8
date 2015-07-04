.PHONY: updatedeps build assets clean server run

BINNAME := c3h8
MAINFILE := c3h8.go
ASSET_PKG := "propane"
ASSET_BINDATA := "./propane/bindata.go"

updatedeps:
	go get -u github.com/codegangsta/cli
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u github.com/russross/blackfriday

build: assets
	go build -o $(BINNAME) $(MAINFILE)

assets:
	go-bindata -pkg=$(ASSET_PKG) -o=$(ASSET_BINDATA) ./assets/...

clean:
	rm $(BINNAME)

server:
	./c3h8 server

run: build server
