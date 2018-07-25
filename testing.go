package main

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

/*
After running this, and then going to mongo shell and running produces this error
> db.people.find();
Thu Aug  9 12:01:36 decode failed. probably invalid utf-8 string [P#��/U�
                                                                         ]
Thu Aug  9 12:01:36 	 why: TypeError: malformed UTF-8 character sequence at offset 2
Error: invalid utf8
*/
type Person struct {
	Id string `bson:"_id,omitempty"`
	Mid bson.ObjectId  `bson:"mid,omitempty"`
	Name string
}

func (p *Person) String() string {
	oid := bson.ObjectId(p.Id)
	return fmt.Sprintf(`<Person Name:"%s" Id:"%s"   hex="%s" mid="%s">`, p.Name, oid.String(), oid.Hex(), p.Mid.String())
}

func main() {
	session, err := mgo.Dial("localhost")

	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB("test2").C("people")
	// How do i get the Id of an object that has just been inserted without querying for it?
	// 1.  Create the ObjectId() in advance?   Which doesn't serialize correctly (invalid errror above)
	// 2.  Id get filled out?   
	bid := []byte{80, 36, 13, 102, 47, 85, 184, 17, 20, 0, 0, 1}
	p1 := &Person{Id: string(bid), Name: "Bob"}
	p2 := &Person{Id: string(bson.NewObjectId()), Name: "Jan"}
	p3 := &Person{Name: "Jim", Mid: bson.NewObjectId()}
	err = c.Insert( p1, p2, p3)

	fmt.Println("initial id1 ", p1)
	fmt.Println("initial id2 ", p2)
	fmt.Println("initial id3 ", p3)
	if err != nil {
		panic(err)
	}
	people := make([]*Person, 0)
	iter := c.Find(nil).Iter()
	err = iter.All(&people)
	if err != nil {
		panic(err)
	}
	for _, p := range people {
		fmt.Println(p)
		fmt.Println([]byte(p.Id))
	}
}
