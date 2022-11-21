package main

// 我用于调试panic的plugin

import (
	"fmt"
	"sort"
	"strings"

	"6.824/mr"
)

func Map(filename string, contents string) []mr.KeyValue {
	err := fmt.Errorf("new error")
	panic(err)
	var kva []mr.KeyValue
	return kva
}

func Reduce(key string, values []string) string {
	// sort values to ensure deterministic output.
	vv := make([]string, len(values))
	copy(vv, values)
	sort.Strings(vv)

	val := strings.Join(vv, " ")
	return val
}
