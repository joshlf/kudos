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

var (
	// Version is set by linking from the build script
	Version string

	// Commit is set by linking from the build script
	Commit string
)
