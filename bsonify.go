package mejson

import (
	"fmt"
	"labix.org/v2/mgo/bson"
)

func Bsonify(in map[string]interface{}) (result bson.M, err error) {
	if false {
		fmt.Println("meow")
	}
	result = bson.M{}
	for key, value := range in {
		if oid, ok := Oid(value); ok {
			result[key] = oid
		} else if v, ok := value.(map[string]interface{}); ok {
			result[key], err = Bsonify(v)
		} else {
			result[key] = value
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
