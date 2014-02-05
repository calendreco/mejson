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
		return nil, fmt.Errorf("you gave me bson.D")
	case string:
		return json.Marshal(v)
	case bson.ObjectId:
		return MarshalObjectId(v)
	default:
		fmt.Printf("type %T\n", v)
		fmt.Printf("stuff: %+v\n", v)
		return json.Marshal(v)
	}
	return nil, nil
}

func MarshalObjectId(in bson.ObjectId) ([]byte, error) {
	result := map[string]string{}
	result["$oid"] = in.Hex()
	return json.Marshal(result)
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
