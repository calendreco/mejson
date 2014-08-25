package mejson

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"testing"
	"time"
)

func TestBson(t *testing.T) {

	data := []struct {
		in   M
		want bson.M
	}{
		{
			map[string]interface{}{"one": 1},
			bson.M{"one": 1},
		},
		{
			map[string]interface{}{"one": []interface{}{"one", "two"}},
			bson.M{"one": S{"one", "two"}},
		},
		{
			map[string]interface{}{"one": map[string]interface{}{"two": 2}},
			bson.M{"one": bson.M{"two": 2}},
		},
	}

	for _, d := range data {
		b, err := d.in.Bson()
		if err != nil {
			t.FailNow()
		}
		if !reflect.DeepEqual(b, d.want) {
			t.Errorf("wanted: %v (%T), got: %v (%T)", d.want, d.want, b, b)
			t.Errorf("one: %v (%T), one: %v (%T)", d.want["one"], d.want["one"], b["one"], b["one"])
		}
	}
}

func TestIsExtended(t *testing.T) {
	data := []struct {
		in   M
		want bool
	}{
		{
			map[string]interface{}{"one": 1},
			false,
		},
		{
			map[string]interface{}{"one": 1, "two": 2},
			false,
		},
		{
			map[string]interface{}{"one": 1, "two": 2, "tree": 3},
			false,
		},
		{
			map[string]interface{}{"$oid": 1, "two": 2},
			false,
		},
		{
			map[string]interface{}{"$oid": 1},
			true,
		},
		{
			map[string]interface{}{"$type": 1, "$binary": 2},
			true,
		},
	}

	for _, d := range data {
		if d.in.isExtended() != d.want {
			t.Errorf("wanted: %v, got: %v", d.want, d.in.isExtended())
		}
	}
}

func TestOid(t *testing.T) {

	data := []struct {
		in     M
		want   bson.ObjectId
		wantok bool
	}{
		{
			map[string]interface{}{"$oid": "52dc18556c528d7736000003"},
			bson.ObjectIdHex("52dc18556c528d7736000003"),
			true,
		},
		{
			map[string]interface{}{"$odd": "52dc18556c528d7736000003"},
			bson.ObjectId(""),
			false,
		},
		{
			map[string]interface{}{"$odd": "52dc18556c528d773600000r"},
			bson.ObjectId(""),
			false,
		},
		{
			map[string]interface{}{"$odd": "52dc18556c528d773600000RRRRRR"},
			bson.ObjectId(""),
			false,
		},
	}

	for _, d := range data {
		b, ok := d.in.Oid()
		if ok != d.wantok {
			t.FailNow()
		}
		if !reflect.DeepEqual(b, d.want) {
			t.Errorf("wanted: %v (%T), got: %v (%T)", d.want, d.want, b, b)
		}
	}
}

func TestDate(t *testing.T) {
	sample_time, _ := time.Parse(time.RFC3339, "2014-02-19T15:14:41.288Z")
	sample_time2, _ := time.Parse(time.RFC3339, "2007-02-19T15:14:41.288Z")
	sample_time3, _ := time.Parse(time.RFC3339, "2014-02-20T11:23:56Z")

	data := []struct {
		in     M
		want   time.Time
		wantok bool
	}{
		{
			map[string]interface{}{"$date": 1392822881288},
			sample_time,
			true,
		},
		{
			map[string]interface{}{"$milwaukee": 1392822881288},
			sample_time2,
			false,
		},
		{
			map[string]interface{}{"$date": 1392895436000},
			sample_time3,
			true,
		},
	}

	for _, d := range data {
		b, ok := d.in.Date()
		if ok != d.wantok {
			t.FailNow()
		}
		if ok && b.UnixNano() != d.want.UnixNano() {
			t.Errorf("wanted: %v (%d) (%T), got: %v (%d) (%T)", d.want, d.want.UnixNano()/1e6, d.want, b, b.UnixNano()/1e6, b)
		}
	}
}

func TestTimestamp(t *testing.T) {
	data := []struct {
		in     M
		want   bson.MongoTimestamp
		wantok bool
	}{
		{
			map[string]interface{}{"$timestamp": map[string]interface{}{"t": 1392822881, "i": 1}},
			bson.MongoTimestamp(5982128723015499777),
			true,
		},
		{
			map[string]interface{}{"$ugh": map[string]interface{}{"t": 1392822881, "i": 1}},
			bson.MongoTimestamp(5982128723015499777),
			false,
		},
	}

	for _, d := range data {
		b, ok := d.in.Timestamp()
		if ok != d.wantok {
			t.Errorf("got %t, want %t, (%v)", ok, d.wantok, d.in)
			t.FailNow()
		}
		if ok && b != d.want {
			t.Errorf("wanted: %v (%T), got: %v (%T)", d.want, d.want, b, b)
		}
	}
}

func TestBinary(t *testing.T) {
	data := []struct {
		in     M
		want   bson.Binary
		wantok bool
	}{
		{
			map[string]interface{}{
				"$binary": "b2ggaGk=",
				"$type":   "00",
			},
			bson.Binary{Kind: 0x00, Data: []byte{111, 104, 32, 104, 105}},
			true,
		},
	}

	for _, d := range data {
		b, ok := d.in.Binary()
		if ok != d.wantok {
			t.Errorf("got %t, want %t, (%v)", ok, d.wantok, d.in)
			t.FailNow()
		}
		if ok && !reflect.DeepEqual(b, d.want) {
			t.Errorf("wanted: %v (%T), got: %v (%T)", d.want, d.want, b, b)
		}
	}
}

func TestBsonify(t *testing.T) {
	sample_time3 := time.Unix(0, int64(1392895436000)*int64(time.Millisecond)) //:= time.Parse(time.RFC3339, "2014-02-20T11:23:56Z")

	data := []struct {
		in   []byte
		want bson.M
	}{
		{
			[]byte("{\"name\":\"jp_with_date\",\"created_at\":{\"$date\":1392895436000}}"),
			bson.M{"name": "jp_with_date", "created_at": sample_time3},
		},
		{
			[]byte("{\"vancouver\":{\"$lt\":5}}"),
			bson.M{"vancouver": bson.M{"$lt": float64(5)}},
		},
	}

	for _, d := range data {
		v := make(map[string]interface{})

		err := json.Unmarshal(d.in, &v)
		if err != nil {
			t.Errorf("%s", err)
			t.FailNow()
		}

		b, err := Bsonify(v)
		if err != nil {
			t.FailNow()
		}
		if !reflect.DeepEqual(b, d.want) {
			t.Errorf("wanted: %v (%T), got: %v (%T)", d.want, d.want, b, b)
		}
	}

}
