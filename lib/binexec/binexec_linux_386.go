package binexec

import (
	"fmt"

	"github.com/joshlf/kudos/lib/binexec/internal/linux_386"
)

func init() {
	loadAsset := func(name string) {
		bytes, err := bindata.Asset(name)
		if err != nil {
			panic(fmt.Errorf("lib/binexec: unexpected error loading asset: %v", err))
		}
		bins[name] = &bin{contents: bytes}
	}

	names := []string{
		"tar",
	}
	for _, name := range names {
		loadAsset(name)
	}
}
