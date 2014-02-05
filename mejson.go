package mejson

import (
	"encoding/json"
	"fmt"
	"labix.org/v2/mgo/bson"
)

func MarshalMEJSON(in interface{}) ([]byte, error) {
	switch v := in.(type) {
	case bson.M:
		return MarshalMapJSON(v)
	case bson.D:
		return nil, fmt.Errorf("you gave me bson.D")
	case string:
		return json.Marshal(v)
	case bson.ObjectId:
		return MarshalObjectIdMEJSON(v)
	default:
		fmt.Printf("type %T\n", v)
		fmt.Printf("stuff: %+v\n", v)
		return json.Marshal(v)
	}
	return nil, nil
}

func MarshalObjectIdMEJSON(in bson.ObjectId) ([]byte, error) {
	result := map[string]string{}
	result["$oid"] = in.Hex()
	return json.Marshal(result)
}

func MarshalMapJSON(in bson.M) ([]byte, error) {
	result := map[string]*json.RawMessage{}
	for key, value := range in {
		bytes, err := MarshalMEJSON(value)
		if err != nil {
			return nil, err
		}
		message := &json.RawMessage{}
		message.UnmarshalJSON(bytes)
		result[key] = message
	}
	return json.Marshal(result)
}
