package mejson

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"labix.org/v2/mgo/bson"
	"time"
)

// Bsonify turns a naive json map[string]interface{} object into a bson.M object.
// It will turn mongodb extended json fields (like $oid) into an appropriate golang object.
// e.g. {"_id":{"$oid":"52ddcd077854f173de9429b3"}} will become {"_id":bson.ObjectId("52ddcd077854f173de9429b3")}
// Find the documentation on MongoDB Extended JSON here: http://docs.mongodb.org/manual/reference/mongodb-extended-json/
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
		} else if timestamp, ok := Timestamp(value); ok {
			result[key] = timestamp
		} else if binary, ok := Binary(value); ok {
			result[key] = binary
		} else if v, ok := value.(map[string]interface{}); ok {
			result[key], err = Bsonify(v)
		} else {
			result[key] = value
		}
	}
	return
}

func getOnly(m map[string]interface{}, s string) (interface{}, bool) {
	value, ok := m[s]
	if ok && len(m) == 1 {
		return value, true
	}
	return nil, false
}

func getTimestamp(m map[string]interface{}) (t uint32, i uint32, ok bool) {
	if len(m) == 2 {
		if t, ok := getInt(m, "t"); ok {
			if i, ok := getInt(m, "i"); ok {
				return t, i, true
			}
		}
	}
	return
}

func getInt(m map[string]interface{}, k string) (uint32, bool) {
	v, ok := m[k]
	if !ok {
		return 0, false
	}
	i, ok := v.(int)
	if !ok {
		return 0, false
	}
	return uint32(i), true
}

func Timestamp(in interface{}) (timestamp bson.MongoTimestamp, ok bool) {
	switch v := in.(type) {
	case map[string]interface{}:
		if value, contains := getOnly(v, "$timestamp"); contains {
			if m, ismap := value.(map[string]interface{}); ismap {
				t, i, istimestamp := getTimestamp(m)
				if istimestamp {
					ok = true
					var concat int64
					concat = int64(uint64(t)<<32 | uint64(i))
					timestamp = bson.MongoTimestamp(concat)
				}
			}
		}
	}
	return
}

func getBinaryKind(m map[string]interface{}) (kind byte, ok bool) {
	v, contains := m["$type"]
	if !contains {
		return
	}
	hexstr, isstr := v.(string)
	if !isstr {
		return
	}
	hexbytes, err := hex.DecodeString(hexstr)
	if err != nil || len(hexbytes) != 1 {
		return
	}
	kind = hexbytes[0]
	ok = true
	return
}

func getBinaryData(m map[string]interface{}) (data []byte, ok bool) {
	v, contains := m["$binary"]
	if !contains {
		return
	}
	binarystr, isstr := v.(string)
	if !isstr {
		return
	}
	bytes, err := base64.StdEncoding.DecodeString(binarystr)
	if err != nil {
		return
	}
	data = bytes
	ok = true
	return
}

func Binary(in interface{}) (binary bson.Binary, ok bool) {
	switch v := in.(type) {
	case map[string]interface{}:
		if len(v) != 2 {
			return
		}
		kind, kindok := getBinaryKind(v)
		if !kindok {
			return
		}
		data, dataok := getBinaryData(v)
		if !dataok {
			return
		}
		binary.Kind = kind
		binary.Data = data
		ok = true
	}
	return
}

func Date(in interface{}) (date time.Time, ok bool) {
	switch v := in.(type) {
	case map[string]interface{}:
		if value, contains := getOnly(v, "$date"); contains {
			if milli, isint := value.(int); isint {
				ok = true
				date = time.Unix(0, int64(milli)*int64(time.Millisecond))
			}
		}
	}
	return
}

// ok if in == {"$oid": "#{ObjectId.Hex()}"}
func Oid(in interface{}) (oid bson.ObjectId, ok bool) {
	switch v := in.(type) {
	case map[string]interface{}:
		if value, contains := getOnly(v, "$oid"); contains {
			if hex, isstr := value.(string); isstr && bson.IsObjectIdHex(hex) {
				oid = bson.ObjectIdHex(hex)
				ok = true
			}
		}
	}
	return
}
