package mejson

import (
	"encoding/base64"
	"encoding/hex"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type M map[string]interface{}

type I interface{}

type S []interface{}

func Bsonify(m map[string]interface{}) (result bson.M, err error) {
	return M(m).Bson()
}

func (m M) Bson() (result bson.M, err error) {
	result = bson.M{}

	for key, value := range m {
		switch v := value.(type) {
		case []interface{}:
			result[key], err = S(v).Bson()
			if err != nil {
				return
			}

		case map[string]interface{}:
			if !M(v).isExtended() {
				result[key], err = M(v).Bson()
				if err != nil {
					return
				}
			} else {
				if oid, ok := M(v).Oid(); ok {
					result[key] = oid
				} else if date, ok := M(v).Date(); ok {
					result[key] = date
				} else if timestamp, ok := M(v).Timestamp(); ok {
					result[key] = timestamp
				} else if binary, ok := M(v).Binary(); ok {
					result[key] = binary
				} else {
					result[key], err = M(v).Bson() // it's ugly to repeat this clause here
					if err != nil {
						return
					}
				}
			}
		default:
			result[key] = v
		}
	}

	return
}

func (m M) isExtended() bool {
	if len(m) != 1 && len(m) != 2 {
		return false
	}

	for k, _ := range m {
		if k[0] != '$' {
			return false
		}
	}

	return true
}

/* $oid type */
func (m M) Oid() (oid bson.ObjectId, ok bool) {
	if len(m) != 1 {
		return
	}
	if value, contains := m["$oid"]; contains {
		if hex, isstr := value.(string); isstr && bson.IsObjectIdHex(hex) {
			oid = bson.ObjectIdHex(hex)
			ok = true
		}
	}
	return
}

// RFC3339Nano with a numeric zone
const ISO8601 = "2006-01-02T15:04:05.999999999-0700"

/* $date type */
func (m M) Date() (date time.Time, ok bool) {
	if len(m) != 1 {
		return
	}

	if value, contains := m["$date"]; contains {
		var millis int
		var err error

		switch m := value.(type) {
		case int:
			millis = m
		case int64:
			millis = int(m)
		case int32:
			millis = int(m)
		case float64:
			millis = int(m)
		case float32:
			millis = int(m)
		// The MongoDB JSON parser currently does not support loading
		// ISO-8601 strings representing dates prior to the Unix epoch.
		// When formatting pre-epoch dates and dates past what your
		// system’s time_t type can hold, the following format is used:
		// { "$date" : { "$numberLong" : "<dateAsMilliseconds>" } }
		case map[string]interface{}:
			if value, contains := m["$numberLong"]; contains {
				millis, err = strconv.Atoi(value.(string))
			} else {
				return
			}
		// In Strict mode, <date> is an ISO-8601 date format with a
		// mandatory time zone field following the template
		// YYYY-MM-DDTHH:mm:ss.mmm<+/-Offset>.
		case string:
			date, err = time.Parse(ISO8601, m)
			ok = (err == nil)
			return
		default:
			return
		}
		ok = (err == nil)
		if ok {
			date = time.Unix(0, int64(millis)*int64(time.Millisecond))
		}
		return
	}

	return
}

/* bsonify a mongo Timestamp */
func (m M) Timestamp() (timestamp bson.MongoTimestamp, ok bool) {
	if len(m) != 1 {
		return
	}

	if value, contains := m["$timestamp"]; contains {
		if ts, ismap := value.(map[string]interface{}); ismap {
			t, isok := ts["t"]
			if !isok {
				return
			}
			tt, isok := t.(int)
			if !isok {
				return
			}

			i, isok := ts["i"]
			if !isok {
				return
			}
			ii, isok := i.(int)
			if !isok {
				return
			}

			ok = true
			var concat int64
			concat = int64(uint64(tt)<<32 | uint64(ii))
			timestamp = bson.MongoTimestamp(concat)
		}
	}

	return
}

/* bsonify a binary data type */
func (m M) Binary() (binary bson.Binary, ok bool) {

	if len(m) != 2 {
		return
	}
	kind, kindok := getBinaryKind(m)
	if !kindok {
		return
	}
	data, dataok := getBinaryData(m)
	if !dataok {
		return
	}
	binary.Kind = kind
	binary.Data = data
	ok = true

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

/* BSONify a slice of somethings */
func (s S) Bson() (out S, err error) {
	out = make(S, len(s))
	for k, v := range s {
		switch elem := v.(type) {
		case []interface{}:
			out[k], err = S(elem).Bson()
			if err != nil {
				return
			}
		case map[string]interface{}:
			out[k], err = M(elem).Bson()
			if err != nil {
				return
			}
		default:
			out[k] = elem
		}
	}
	return
}
