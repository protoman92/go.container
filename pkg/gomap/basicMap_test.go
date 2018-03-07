package gomap

import (
	"fmt"
	"testing"

	"github.com/protoman92/gocontainer/pkg/gocollection"
)

func TestBasicMapConvertKey(t *testing.T) {
	/// Setup
	convertKeyFn := func(e interface{}) string {
		converted := fmt.Sprint(e)
		return converted + converted
	}

	params := BasicMapParams{ConvertKeyFn: convertKeyFn}
	m := NewBasicMap(params)
	keys := []interface{}{1, 2, 3, 4, 5, 6}
	convertedKeys := make([]interface{}, len(keys))

	for ix := range keys {
		convertedKeys[ix] = convertKeyFn(keys[ix])
	}

	/// When
	for k := range keys {
		key := keys[k]
		m.Set(key, key)
	}

	/// Then
	mapKeys := m.Keys()
	bl := gocollection.NewSliceList(mapKeys)

	for ix := range convertedKeys {
		k := convertedKeys[ix]

		if contains := bl.Contains(k); !contains {
			t.Errorf("Should contain key")
		}
	}
}
