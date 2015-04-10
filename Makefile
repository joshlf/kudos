build: deps
	go build -o bin/kudos -a cmd/*.go

dev: deps
	go build -o bin/kudos -a -tags dev cmd/*.go

debug: deps
	go build -o bin/kudos -a -tags debug cmd/*.go

# TODO: Figure out clean way of doing `-tags dev debug`

deps: bin-dir

bin-dir:
	mkdir -p bin

clean:
	rm -rf bin
