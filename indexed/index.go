package main

import (
	"fmt"
)

type Index struct {
	loaded   map[string]interface{}
	location string
}

func (i *Index) get(key string) (interface{}, error) {
	if v, exists := loaded[key]; exists {
		return v
	}
	return fmt.Errorf(key + " does not exist")
}
