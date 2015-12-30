LDFLAGS="-X github.com/joshlf/kudos/lib/build.Version=$(shell git rev-parse HEAD) \
		-X github.com/joshlf/kudos/lib/build.Commit=$(shell git rev-parse HEAD)"

build: deps
	go build -ldflags $(LDFLAGS) -o bin/kudos -a cmd/*.go

dev: deps
	go build -ldflags $(LDFLAGS) -o bin/kudos -a -tags dev cmd/*.go

debug: deps
	go build -ldflags $(LDFLAGS) -o bin/kudos -a -tags debug cmd/*.go

# TODO: Figure out clean way of doing `-tags dev debug`

deps: bin-dir

bin-dir:
	mkdir -p bin

clean:
	rm -rf bin
