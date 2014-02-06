package main

import (
	"fmt"
	"github.com/supershabam/mejson"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func main() {
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
