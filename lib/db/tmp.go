// +build ignore

package db

import "encoding/json"

func concreteToInterface(v interface{}) (interface{}, error) {
	var ret interface{}
	buf, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buf, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
