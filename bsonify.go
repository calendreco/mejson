package mejson

import (
	"fmt"
	"labix.org/v2/mgo/bson"
)

func Bsonify(in map[string]interface{}) (result bson.M, err error) {
	result = make(bson.M)
	for key, value := range in {
		if err != nil {
			break
		}
		switch key {
		case "$oid":
			result[key], err = ParseOid(value)
		default:
			switch v := value.(type) {
			case map[string]interface{}:
				result[key], err = Bsonify(v)
			default:
				result[key] = value
			}
		}
	}
	return
}

func ParseOid(in interface{}) (result bson.ObjectId, err error) {
	switch v := in.(type) {
	case string:
		if bson.IsObjectIdHex(v) {
			result = bson.ObjectIdHex(v)
		} else {
			err = fmt.Errorf("%s is not a valid ObjectId hex value", v)
		}
	default:
		err = fmt.Errorf("expected $oid value to be string")
	}
	return
}
