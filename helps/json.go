package helps

import (
	"encoding/json"
	"log"
)

func DumpJSON(v interface{}) string {
	dump, err := DumpJSONSafe(v)
	if err != nil {
		log.Panic(err)
	}
	return dump
}

func DumpJSONSafe(v interface{}) (string, error) {
	dump, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}

	return string(dump), nil
}
