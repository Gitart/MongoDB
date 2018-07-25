## Sample

```golang
package main

import (
    "fmt"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

type Person struct {
    Id      bson.ObjectId `json:"id" bson:"_id,omitempty"`
    Name  string
    Phone string
}

func checkError(err error) {
    if err != nil {
        panic(err)
    }
}

const (
    DB_NAME       = "gotest"
    DB_COLLECTION = "pepole_new1"
)

func main() {
    session, err := mgo.Dial("localhost")
    checkError(err)
    defer session.Close()

    session.SetMode(mgo.Monotonic, true)

    c := session.DB(DB_NAME).C(DB_COLLECTION)
    err = c.DropCollection()
    checkError(err)

    ale := Person{Name:"Ale", Phone:"555-5555"}
    cla := Person{Name:"Cla", Phone:"555-1234-2222"}
    kasaun := Person{Name:"kasaun", Phone:"533-12554-2222"}
    chamila := Person{Name:"chamila", Phone:"533-545-6784"}

    fmt.Println("Inserting")
    err = c.Insert(&ale, &cla, &kasaun, &chamila)
    checkError(err)

    fmt.Println("findbyID")
    var resultsID []Person
    //err = c.FindId(bson.ObjectIdHex("56bdd27ecfa93bfe3d35047d")).One(&resultsID)
    err = c.FindId(bson.M{"Id": bson.ObjectIdHex("56bdd27ecfa93bfe3d35047d")}).One(&resultsID)
    checkError(err)
    if err != nil {
        panic(err)
    }
    fmt.Println("Phone:", resultsID)



    fmt.Println("Queryingall")
    var results []Person
    err = c.Find(nil).All(&results)

    if err != nil {
        panic(err)
    }
    fmt.Println("Results All: ", results)
}
```
