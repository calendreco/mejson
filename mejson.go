package mejson

import (
	"encoding/json"
	"fmt"
	"labix.org/v2/mgo/bson"
)

func Marshal(in interface{}) ([]byte, error) {
	switch v := in.(type) {
	case bson.M:
		return MarshalMap(v)
	case bson.D:
		return nil, fmt.Errorf("unhandled input type bson.D")
	case string:
		return json.Marshal(v)
	case bson.ObjectId:
		return MarshalObjectId(v)
	default:
		return json.Marshal(v)
	}
}

func MarshalObjectId(in bson.ObjectId) ([]byte, error) {
	return json.Marshal(map[string]string{"$oid": in.Hex()})
}

func MarshalMap(in bson.M) ([]byte, error) {
	result := map[string]*json.RawMessage{}
	for key, value := range in {
		bytes, err := Marshal(value)
		if err != nil {
			return nil, err
		}
		message := &json.RawMessage{}
		message.UnmarshalJSON(bytes)
		result[key] = message
	}
	return json.Marshal(result)
}
