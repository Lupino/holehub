BIN=assets/download
SRC=github.com/Lupino/holehub/holehub

all: build-all zip


build-all:
	GOOS=linux GOARCH=amd64 make build
	GOOS=linux GOARCH=arm make build
	GOOS=darwin GOARCH=amd64 make build
	# GOOS=windows GOARCH=amd64 make build


build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BIN)/$(GOOS)-$(GOARCH)/holehub $(SRC)

zip:
	find -name "holehub"  | while read F;do zip -D -j $$F $$F;done

clean:
	rm -r download
