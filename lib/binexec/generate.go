package binexec

//go:generate mkdir -p internal/linux_amd64
//go:generate go-bindata -o internal/linux_amd64/bindata.go -prefix bin/linux_amd64 -pkg bindata bin/linux_amd64/

//go:generate mkdir -p internal/linux_386
//go:generate go-bindata -o internal/linux_386/bindata.go -prefix bin/linux_386 -pkg bindata bin/linux_386/

// clean up after go-bindata, which produces
// improperly-formatted source files
//go:generate go fmt ./...
