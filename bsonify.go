package mejson

import (
	"fmt"
	"labix.org/v2/mgo/bson"
	"time"
)

func Bsonify(in map[string]interface{}) (result bson.M, err error) {
	if false {
		fmt.Println("meow")
	}
	result = bson.M{}
	for key, value := range in {
		if oid, ok := Oid(value); ok {
			result[key] = oid
		} else if date, ok := Date(value); ok {
			result[key] = date
		} else if v, ok := value.(map[string]interface{}); ok {
			result[key], err = Bsonify(v)
		} else {
			result[key] = value
		}
	}
	return
}

func Date(in interface{}) (date time.Time, ok bool) {
	ok = false
	switch v := in.(type) {
	case map[string]interface{}:
		value, contains := v["$date"]
		milli, isint := value.(int)
		if isint && contains && len(v) == 1 {
			ok = true
			date = time.Unix(0, int64(milli)*int64(time.Millisecond))
		}
	}
	return
}

// ok if in == {"$oid": "#{ObjectId.Hex()}"}
func Oid(in interface{}) (oid bson.ObjectId, ok bool) {
	ok = false
	switch v := in.(type) {
	case map[string]interface{}:
		value, contains := v["$oid"]
		hex, isstr := value.(string)
		if isstr && contains && len(v) == 1 && bson.IsObjectIdHex(hex) {
			oid = bson.ObjectIdHex(hex)
			ok = true
		}
	}
	return
}
