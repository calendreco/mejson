package mejson

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"labix.org/v2/mgo/bson"
	"reflect"
	"time"
)

func Marshal(in interface{}) ([]byte, error) {
	if reflect.TypeOf(in).Kind() == reflect.Slice {
		v := reflect.ValueOf(in)
		slice := []interface{}{}
		for i := 0; i < v.Len(); i++ {
			slice = append(slice, v.Index(i).Interface())
		}
		return MarshalSlice(slice)
	} else {
		switch v := in.(type) {
		case bson.M:
			return MarshalMap(v)
		case bson.D:
			// todo write marshaller for doc to ensure serialization order
			return MarshalMap(v.Map())
		case bson.Binary:
			return MarshalBinary(v)
		case bson.ObjectId:
			return MarshalObjectId(v)
		case time.Time:
			return nil, fmt.Errorf("unimplemented type time.Time")
		case string:
			return json.Marshal(v)
		default:
			fmt.Printf("unknown type: %T\n", v)
			return json.Marshal(v)
		}
	}
}

func MarshalSlice(in []interface{}) ([]byte, error) {
	result := []*json.RawMessage{}
	for _, value := range in {
		bytes, err := Marshal(value)
		if err != nil {
			return nil, err
		}
		message := &json.RawMessage{}
		message.UnmarshalJSON(bytes)
		result = append(result, message)
	}
	return json.Marshal(result)
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

func MarshalBinary(in bson.Binary) ([]byte, error) {
	result := map[string]string{
		"$type":   fmt.Sprintf("%x", in.Kind),
		"$binary": base64.StdEncoding.EncodeToString(in.Data),
	}
	return json.Marshal(result)
}
