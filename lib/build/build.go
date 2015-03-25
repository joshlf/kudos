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
