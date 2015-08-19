package build

func DevDo(f func()) {
	if DevMode {
		f()
	}
}

func DebugDo(f func()) {
	if DebugMode {
		f()
	}
}

// Version is set by linking from the build script

var Version string
