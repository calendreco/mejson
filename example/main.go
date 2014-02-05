package main

import (
	"fmt"
	"github.com/supershabam/mejson"
	"labix.org/v2/mgo/bson"
)

func main() {
	fmt.Println("oh hai")
	v := bson.M{"test": bson.NewObjectId()}
	bytes, err := mejson.Marshal(v)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("%s\n", bytes)
}
