### Sample for MongoDB

```golang
package main

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

type Article struct {
	Id      bson.ObjectId `json:"id"        bson:"_id,omitempty"`
	Title   string        `json:"title"`
	Author  string        `json:"author"`
	Date    string        `json:"date"`
	Tags    string        `json:"tags"`
	Content string        `json:"content"`
	Status  string        `json:"status"`
}

func AllArticles() []Article {
	articles := []Article{}
	err := c_articles.Find(bson.M{}).All(&articles)
	if err != nil {
		panic(err)
	}

	return articles
}

var c_articles *mgo.Collection

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	}
	c_articles = session.DB("test").C("articles")
	err = c_articles.Insert(bson.M{"title": "Some Title"})
	if err != nil {
		log.Fatal(err)
	}
	all := AllArticles()
	fmt.Printf("%#v %#v\n", all[0].Id, all[0].Title)
}
```

### Link
https://github.com/swhite24/go-rest-tutorial

