BIN=assets/download
PKG=github.com/Lupino/holehub/holehub

all: build-all zip


build-all:
	GOOS=linux GOARCH=amd64 make build
	GOOS=linux GOARCH=arm make build
	GOOS=darwin GOARCH=amd64 make build
	# GOOS=windows GOARCH=amd64 make build

$(GOPATH)/src/$(PKG):
	go get -v $(PKG)


build: $(GOPATH)/src/$(PKG)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BIN)/$(GOOS)-$(GOARCH)/holehub $(PKG)

zip:
	find -name "holehub"  | while read F;do zip -D -j $$F $$F;done

clean:
	rm -r $(BIN)
