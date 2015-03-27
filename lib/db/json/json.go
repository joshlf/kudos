// Package json provides a JSON-backed DB implementation
// as expected by the interfaces defined in package db.
package json

import (
	"fmt"
	"os"

	"github.com/synful/kudos/lib/db"
)

const (
	ProviderName = "json"
)

// connectionStatus represents the current status
// of a connection to the JSON provider
type connectionStatus int

const (
	unconnected connectionStatus = iota
	connected
	closed
)

type jsonDB struct {
	config jsonConfig
	status connectionStatus
	handle *os.File
}

type jsonConfig struct {
	path string
}

// exists checks to see if the database exists.
// It returns true if it exists, false if it doesn't,
// and propagates any I/O errors that occur while checking.
func (jp *jsonDB) exists() (bool, error) {
	if _, err := os.Stat(jp.config.path); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (jp *jsonDB) Init() error {
	if present, err := jp.exists(); err != nil {
		return fmt.Errorf("json.Init: Failed to init db.\n%v", err)
	} else if present {
		return fmt.Errorf("json.Init: Database already exists")
	}
	return nil
}

func (jp *jsonDB) Connect() (db.Conn, error) {
	if present, err := jp.exists(); err != nil {
		return nil, err
	} else if !present {
		return nil, fmt.Errorf("json.Connect: Database not found")
	} else {
		if jp.handle, err = os.OpenFile(jp.config.path, os.O_RDWR, 0660); err != nil {
			return nil, err
		}
		jp.status = connected
		return jp, nil

	}
}

func (jp *jsonDB) Destroy() error {
	if present, err := jp.exists(); err != nil {
		return fmt.Errorf("json.Destroy: Failed to delete db.\n%v", err)
	} else if present {
		os.Remove(jp.config.path)
		return nil
	}
	return fmt.Errorf("json.Destroy: Database does not exist")
}

func (jp *jsonDB) Close() error {
	// TODO(mdburns)
	return nil
}

func (jp *jsonDB) Query(path db.Path, constraint ...db.Constraint) ([]db.Entity, error) {
	// TODO(mdburns)
	return nil, nil
}

func connect(config interface{}) (db.DB, error) {
	cfg := config.(jsonConfig)

	return &jsonDB{
		config: cfg,
		status: unconnected,
	}, nil
}

func init() {
	db.RegisterProvider(ProviderName, connect)
}
