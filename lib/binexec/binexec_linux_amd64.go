package binexec

import "github.com/joshlf/kudos/lib/binexec/internal/linux_amd64"

func init() {
	assets = []asset{
		{"tar", []string{
			"libacl.so",
			"libattr.so",
			"libdl.so",
			"libpcre.so",
			"libpthread.so",
			"libselinux.so"},
		},
	}

	MustAsset = bindata.MustAsset
}
