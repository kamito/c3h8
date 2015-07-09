.PHONY: updatedeps build assets clean server run

BINNAME := c3h8
MAINFILE := c3h8.go
ASSET_PKG := "propane"
ASSET_BINDATA := "./propane/bindata.go"

updatedeps:
	go get -u github.com/codegangsta/cli
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u github.com/russross/blackfriday

assets:
	go-bindata -pkg=$(ASSET_PKG) -o=$(ASSET_BINDATA) ./assets/...

build: assets
	go build -o $(BINNAME) $(MAINFILE)

run: assets
	go run $(MAINFILE)

clean:
	rm $(BINNAME)

server:
	./c3h8 server

run: build server
