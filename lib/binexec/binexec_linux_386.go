package binexec

import "github.com/joshlf/kudos/lib/binexec/internal/linux_386"

func init() {
	assets = []asset{
		{"tar", []string{"libacl.so", "libselinux.so"}},
	}

	MustAsset = bindata.MustAsset
}
