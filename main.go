// Copyright © 2016 Charles Phillips <charles@doublerebel.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package bellows

import (
	"fmt"
	"reflect"
	"strings"
)

func Expand(value map[string]interface{}) map[string]interface{} {
	return ExpandPrefixed(value, "")
}

func ExpandPrefixed(value map[string]interface{}, prefix string) map[string]interface{} {
	m := make(map[string]interface{})
	ExpandPrefixedToResult(value, prefix, m)
	return m
}

func ExpandPrefixedToResult(value map[string]interface{}, prefix string, result map[string]interface{}) {
	if prefix != "" {
		prefix += "."
	}
	for k, val := range value {
		if !strings.HasPrefix(k, prefix) {
			continue
		}

		key := k[len(prefix):]
		idx := strings.Index(key, ".")
		if idx != -1 {
			key = key[:idx]
		}
		if _, ok := result[key]; ok {
			continue
		}
		if idx == -1 {
			result[key] = val
			continue
		}

		// It contains a period, so it is a more complex structure
		result[key] = ExpandPrefixed(value, k[:len(prefix)+len(key)])
	}
}

func Flatten(value interface{}) map[string]interface{} {
	return FlattenPrefixed(value, "")
}

func FlattenPrefixed(value interface{}, prefix string) map[string]interface{} {
	m := make(map[string]interface{})
	FlattenPrefixedToResult(value, prefix, m)
	return m
}

func FlattenPrefixedToResult(value interface{}, prefix string, m map[string]interface{}) {
	original := reflect.ValueOf(value)
	kind := original.Kind()
	if kind == reflect.Ptr || kind == reflect.Interface {
		original = reflect.Indirect(original)
		kind = original.Kind()
	}

	if !original.IsValid() {
		if prefix != "" {
			m[prefix] = nil
		}
		return
	}

	t := original.Type()

	switch kind {
	case reflect.Map:
		if t.Key().Kind() != reflect.String {
			break
		}
		keys := original.MapKeys()
		base := ""
		if prefix != "" {
			base = prefix + "."
		}
		for _, childKey := range keys {
			childValue := original.MapIndex(childKey)
			FlattenPrefixedToResult(childValue.Interface(), base+childKey.String(), m)
		}
	case reflect.Struct:
		numField := original.NumField()
		base := ""
		if prefix != "" {
			base = prefix + "."
		}
		for i := 0; i < numField; i += 1 {
			childValue := original.Field(i)
			childKey := t.Field(i).Name
			FlattenPrefixedToResult(childValue.Interface(), base+childKey, m)
		}
	case reflect.Array, reflect.Slice:
		l := original.Len()
		base := prefix
		for i := 0; i < l; i++ {
			childValue := original.Index(i)
			FlattenPrefixedToResult(childValue.Interface(), fmt.Sprintf("%s[%d]", base, i), m)
		}
	default:
		if prefix != "" {
			m[prefix] = value
		}
	}
}
