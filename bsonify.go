package mejson

import (
	"fmt"
	"labix.org/v2/mgo/bson"
)

func Bsonify(in map[string]interface{}) (result bson.M, err error) {
	result = bson.M{}
	for key, value := range in {
		if oid, ok := Oid(value); ok {
			result[key] = oid
		} else {
			result[key] = value
		}
	}
	return
}

// true if in == {"$oid": "#{ObjectId.Hex()}"}
func Oid(in interface{}) (oid bson.ObjectId, ok bool) {
	ok = false
	switch v := in.(type) {
	case map[string]interface{}:
		value, contains := v["$oid"]
		if hex, isstr := value.(string); isstr && contains && len(v) == 1 && bson.IsObjectIdHex(hex) {
			oid = bson.ObjectIdHex(hex)
			ok = true
		}
	}
	return
}
