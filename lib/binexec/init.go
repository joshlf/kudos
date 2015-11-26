package binexec

// safe to leave nil by default because if it's not set,
// assets will also be nil, and thus the for loop in the
// init function will not execute
var MustAsset func(string) []byte

type asset struct {
	name string
	libs []string
}

var assets []asset

func init() {
	for _, asset := range assets {
		bytes := MustAsset(asset.name)
		l := make([]lib, 0)
		for _, ll := range asset.libs {
			bytes := MustAsset("lib/" + ll)
			l = append(l, lib{ll, bytes})
		}
		bins[asset.name] = &bin{contents: bytes, libs: l}
	}
}
