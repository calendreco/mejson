package main

import (
	"fmt"
	"github.com/supershabam/mejson"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func marsh() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("test").C("people")
	result := []bson.M{}
	err = c.Find(bson.M{}).All(&result)
	if err != nil {
		panic(err)
	}
	bytes, err := mejson.Marshal(result)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", bytes)
}

func bsonify() {
	in := map[string]interface{}{
		"_id": map[string]interface{}{
			"$oid": "1234",
		},
	}
	m, err := mejson.Bsonify(in)
	if err != nil {
		panic(err)
	}
	fmt.Printf("bson: %+v\n", m)
}

func main() {
	bsonify()
}
