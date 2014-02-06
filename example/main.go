package main

import (
	"fmt"
	"github.com/supershabam/mejson"
	"labix.org/v2/mgo/bson"
	"time"
)

func main() {
	fmt.Println("oh hai")
	v := bson.D{{"test", bson.NewObjectId()}, {"time", time.Second}}
	bytes, err := mejson.Marshal(v)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("%s\n", bytes)
}
