# Building a REST Service with Golang - Part 3 (Adding a Backend)

### Intro
Golang is pretty hot right now. Over the past few weeks, I've been exploring implementing some of the cloud infrastracture I'd previously built with node in go, partly for fun, partly because go is fast. This led me to believe that setting up a basic REST service would make for a great tutorial.
At the end, we will have a fully functioning REST service for serving basic CRUD on a user resource. Code for the final product can be found here.
This is the third and final entry in a series of high level tutorials on setting up a REST service using golang, focused on wiring up our previous webserver to a mongo backend.

### Setup mgo 
At the end of the last tutorial, we had a functioning, albeit useless, webserver offering some CRUD operations on a user resource. In this entry, we will be tying that to a mongodb backend for persistent data. If you don't have mongo installed on your system, you can find instructions here or install via brew.
The most popular mongo driver for golang is easily mgo. We'll be using two packages provided by mgo, mgo itself and bson. We need to grab each package to get going.

```
$ go get gopkg.in/mgo.v2
$ go get gopkg.in/mgo.v2/bson
```

We also need to establish a connection to mongo to provide to our controllers. Add the following function to our server.go.

### Session
func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://localhost")

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}
	return s
}
```


Now, our controller is going to need a mongo session to use in the CRUD methods. 
Let's change how we get a UserController to the following.

```golang
// Get a UserController instance
uc := controllers.NewUserController(getSession())
Update the UserController struct
The first thing we need to do to our UserController is to extend the struct to contain a reference to a *mgo.Session so that our controller methods can access mongo. Update the UserController definition to match the following.

UserController struct {
	session *mgo.Session
}
```

Next, we need to update our NewUserController function to receive a *mgo.Session and instantiate the controller with it. Update NewUserController to match the following.

```golang
func NewUserController(s *mgo.Session) *UserController {
	return &UserController{s}
}
```

### Update the User Model

Before we start updating the controller methods to use mongo, we need to update our model to 
integrate with mgo. Similar to how we updated the user model to use a json struct tag for 
outputting JSON data, we need to add a struct tag to tell mgo how to store the user information.

### Update user.go in the models directory to contain the following.

```golang
package models

import "gopkg.in/mgo.v2/bson"

type (
	// User represents the structure of our resource
	User struct {
		Id     bson.ObjectId `json:"id" bson:"_id"`
		Name   string        `json:"name" bson:"name"`
		Gender string        `json:"gender" bson:"gender"`
		Age    int           `json:"age" bson:"age"`
	}
)
```

You'll notice that we are importing only the bson package from mgo. The bson package provides an implementation of the bson specification for go. We use this to change the type of our Id field to a bson.ObjectId.
We also extend each field with a bson struct tag to describe the property name of the mongo document containing the user.

### Update the POST Controller
Now that our model has been updated and our controller has a reference to an active mongo session, let's start integrating mongo into our service by updating the CreateUser method of our UserController. Update the method to the following.


```golang
// CreateUser creates a new user resource
func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Stub an user to be populated from the body
	u := models.User{}

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&u)

	// Add an Id
	u.Id = bson.NewObjectId()

	// Write the user to mongo
	uc.session.DB("go_rest_tutorial").C("users").Insert(u)

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
}
```

There are two significant changes here. First, we ask the bson package for a new ObjectId to store on our user. Second, we write the user to the users collection in the go_rest_tutorial database.

### Let's test adding a user.

```
$ curl -XPOST -H 'Content-Type: application/json' -d '{"name": "Bob Smith", "gender": "male", "age": 50}' http://localhost:3000/user
{"id":"5497246c380a967ff1000003","name":"Bob Smith","gender":"male","age":50}
```

Awesome!   
Our user was stored in mongo and we have successfully generated an ObjectId.


### Update the GET Controller
Now that we are able to create new users, it would be nice to be able to fetch them as well. To do this, let's update the GetUser method of our UserController to the following.

```golang
// GetUser retrieves an individual user resource
func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Stub user
	u := models.User{}

	// Fetch user
	if err := uc.session.DB("go_rest_tutorial").C("users").FindId(oid).One(&u); err != nil {
		w.WriteHeader(404)
		return
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}
```

Again, there a couple of significant changes. First, we are converting our id path parameter to a bson.ObjectId to be used to find the user. In the event that the id cannot be converted, we bail with a 404 Not Found. Second, we use our mongo session to find a user matching the provided id from the users collection. We also added an error check when we find the user to deliver a 404 Not Found.

###  Let's test it out by grabbing the user we created previously.

```
$ curl http://localhost:3000/user/5497246c380a967ff1000003
{"id":"5497246c380a967ff1000003","name":"Bob Smith","gender":"male","age":50}
```

Sweet! We got the same user that we had posted above. Now if we change the id we should get a 404.

```
$ curl -i http://localhost:3000/user/5497246c380a967ff1000004
HTTP/1.1 404 Not Found
Date: Sun, 21 Dec 2014 19:58:34 GMT
Content-Length: 0
Content-Type: text/plain; charset=utf-8
```

### Update the DELETE Controller
So we are able to create new users and retrieve them, now let's add the ability to remove them. To accomplish this, let's update the RemoveUser method of our controller to the following.

```golang
// RemoveUser removes an existing user resource
func (uc UserController) RemoveUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Remove user
	if err := uc.session.DB("go_rest_tutorial").C("users").RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}

	// Write status
	w.WriteHeader(200)
}
```

As you can see, we are again getting an ObjectId out of the id path parameter, and bailing with a 404 if the id cannot by converted to an ObjectId. We then remove the document matching this id and deliver a 200.

### We can test this by first deleting the user we created, then trying to fetch the user.

```
$ curl -XDELETE http://localhost:3000/user/5497246c380a967ff1000003
$ curl -i http://localhost:3000/user/5497246c380a967ff1000003
HTTP/1.1 404 Not Found
Date: Sun, 21 Dec 2014 20:04:28 GMT
Content-Length: 0
Content-Type: text/plain; charset=utf-8
```

Great! We removed the user and confirmed that the user no longer exists.
Fin
This is end of the tutorial on adding a backend to our REST service, and also the end of the tutorial series in general. While there most certainly are many improvements to be made, such as adding some security and much better error handling, hopefully this high-level overview has been beneficial.

Thanks to the author :
https://stevenwhite.com/building-a-rest-service-with-golang-3/
