package db

type jsonProvider struct {
}

func (jp *jsonProvider) Open() (error,DBProvider) {
}

func (jp *jsonProvider) Commit() (error,DBProvider) {
}

func (jp *jsonProvider) Init() (error,DBProvider) {
}

func (jp *jsonProvider) Destroy() error {
}

func (jp *jsonProvider) Query(path []string, constraints []DBConstraint) (error,[]DBEntity) {
}
