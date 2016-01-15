package model

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/geotrace/geo"
	"github.com/kr/pretty"
	"github.com/mdigger/rest"
	"github.com/ugorji/go/codec"
)

func TestData(t *testing.T) {
	hjson := new(codec.JsonHandle)
	hjson.Canonical = true // сортировать ключи в словаре
	hjson.Indent = -1      // отступ с табуляцией
	enc := codec.NewEncoder(os.Stdout, hjson)

	circle := geo.NewCircle(55.3980239842, 88.9283429834, 500)
	err := enc.Encode(&Place{
		ID:      "id",
		GroupID: "group_id",
		Name:    "name",
		Circle:  &circle,
	})
	if err != nil {
		t.Error(err)
	}
	err = enc.Encode([]*Place{
		{
			ID:      "id",
			GroupID: "group_id",
			Name:    "name",
			Circle:  &circle,
		},
		{
			ID:      "id",
			GroupID: "group_id",
			Name:    "name",
			Circle:  &circle,
		},
		{
			ID:      "id",
			GroupID: "group_id",
			Name:    "name",
			Circle:  &circle,
		},
	})
	if err != nil {
		t.Error(err)
	}
	err = enc.Encode(circle.Geo())
	if err != nil {
		t.Error(err)
	}

	data, err := json.Marshal(rest.JSON{
		"id":      "id",
		"groupId": "group_id",
		"name":    "test_place",
		"circle": rest.JSON{
			"center": geo.Point{88, 55},
			"radius": 500,
		},
	})
	if err != nil {
		t.Error(err)
	}
	dec := codec.NewDecoderBytes(data, hjson)
	place := new(Place)
	err = dec.Decode(place)
	if err != nil {
		t.Error(err)
	}
	pretty.Println(place)
}
