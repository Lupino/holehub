BIN=assets/download
PKG=github.com/Lupino/holehub/holehub

all: build-all zip zip-exe


build-all:
	GOOS=linux GOARCH=386 make build
	GOOS=linux GOARCH=amd64 make build
	GOOS=linux GOARCH=arm make build
	GOOS=darwin GOARCH=386 make build
	GOOS=darwin GOARCH=amd64 make build
	GOOS=windows GOARCH=386 make build-exe
	GOOS=windows GOARCH=amd64 make build-exe

$(GOPATH)/src/$(PKG):
	go get -v $(PKG)


build: $(GOPATH)/src/$(PKG)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BIN)/$(GOOS)-$(GOARCH)/holehub $(PKG)

build-exe: $(GOPATH)/src/$(PKG)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BIN)/$(GOOS)-$(GOARCH)/holehub.exe $(PKG)

zip:
	find -name "holehub"  | while read F;do zip -D -j $$F $$F; rm $$F; done

zip-exe:
	find -name "holehub.exe"  | while read F;do zip -D -j $${F%.exe} $$F; rm $$F; done

clean:
	rm -rf $(BIN)
	gulp clean

prepare:
	npm install

dist:
	gulp

serve:
	gulp serve
