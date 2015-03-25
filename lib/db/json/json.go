// Package db.json provides a JSON-backed DB implementation
// as expected by the interfaces defined in package db.
package db.json

import (
  "github.com/synful/kudos/lib/db"
  "os"
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
  handle File
}

type jsonConfig struct {
  path string
}

// exists checks to see if the database exists.
// It returns true if it exists, false if it doesn't,
// and propagates any I/O errors that occur while checking.
func (jp *jsonDB) exists() (bool, error) {
  if _, err := os.Stat(jp.config.path); os.IsNotExists(err) {
    return false, nil
  } else if err != nil {
    return false, err
  } else {
    return true, nil
  }
}

func (jp *jsonDB) Init() error {
  if present, err := jp.exists(); err != nil {
    return fmt.errorf("db.json.Init: Failed to init db.\n%v", err)
  } else if present {
    return fmt.errorf("db.json.Init: Database already exists.")
  } else {
    return nil
  }
}

func (jp *jsonDB) Connect() (db.Conn, error) {
  if present, err := jp.exists(); err != nil {
    return nil, err
  } else if !present {
    return nil, fmt.errorf("db.json.Connect: Database not found.")
  } else {
    if jp.handle, err := os.OpenFile(jp.config.path, O_RDWR, 0660); err != nil {
      return nil, err
    } else {
      jp.status = connected
      return jp, nil
    }
  }
}

func (jp *jsonDB) Destroy() error {
  if present, err := jp.exists(); err != nil {
    return fmt.errorf("db.json.Destroy: Failed to delete db.\n%v", err)
  } else if present {
    os.Remove(jp.config.path)
    return nil
  } else {
    return fmt.errorf("db.json.Destroy: Database does not exist.")
  }
}

func connect(config interface{}) (db.DB, error) {
  cfg := config.(jsonConfig)

  return jsonDB {
    config: cfg,
    status: unconnected,
  }
}

func init() {
  db.RegisterProvider(ProviderName, connect)
}
