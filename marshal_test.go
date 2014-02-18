package mejson

import (
	"fmt"
	"labix.org/v2/mgo/bson"
	"reflect"
	"testing"
	"time"
)

func TestMarshal(t *testing.T) {
	data := []struct {
		in      interface{}
		want    []byte
		wanterr error
	}{
		{
			bson.ObjectIdHex("52dc18556c528d7736000003"),
			[]byte("{\"$oid\":\"52dc18556c528d7736000003\"}"),
			nil,
		},
		{
			time.Now(),
			nil,
			fmt.Errorf("unimplemented type time.Time"),
		},
	}

	for _, d := range data {
		b, err := Marshal(d.in)
		if err != nil && err.Error() != d.wanterr.Error() {
			t.Errorf("wanted!: %s, got: %s", d.wanterr, err)
			t.FailNow()
		}
		if err == nil && err != d.wanterr {
			t.Errorf("wanted: %s, got: %s", d.wanterr, err)
			t.FailNow()
		}
		if !reflect.DeepEqual(b, d.want) {
			t.Errorf("wanted: %s, got: %s", d.want, b)
		}
	}
}

func TestMarshalObjectId(t *testing.T) {
	data := []struct {
		in   bson.ObjectId
		want []byte
	}{
		{bson.ObjectIdHex("52dc18556c528d7736000003"), []byte("{\"$oid\":\"52dc18556c528d7736000003\"}")},
		{bson.ObjectIdHex("deadbeefcafedeadbeedcafe"), []byte("{\"$oid\":\"deadbeefcafedeadbeedcafe\"}")},
	}

	for _, d := range data {
		b, err := MarshalObjectId(d.in)
		if err != nil {
			t.FailNow()
		}
		if !reflect.DeepEqual(b, d.want) {
			t.Errorf("wanted: %s, got: %s", d.want, b)
		}
	}
}

func TestMarshalBinary(t *testing.T) {
	data := []struct {
		in   bson.Binary
		want []byte
	}{
		{bson.Binary{Kind: 0x80, Data: []byte("52dc18556c528d7736000003")}, []byte("{\"$binary\":\"NTJkYzE4NTU2YzUyOGQ3NzM2MDAwMDAz\",\"$type\":\"80\"}")},
	}

	for _, d := range data {
		b, err := MarshalBinary(d.in)
		if err != nil {
			t.FailNow()
		}
		if !reflect.DeepEqual(b, d.want) {
			t.Errorf("wanted: %s, got: %s", d.want, b)
		}
	}
}
