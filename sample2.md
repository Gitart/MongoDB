## Building a REST Service with Golang - Part 2 (Basic Webserver)
Steven White

### Intro
Golang is pretty hot right now. Over the past few weeks, I've been exploring implementing some of the cloud infrastracture I'd previously built with node in go, partly for fun, partly because go is fast. This led me to believe that setting up a basic REST service would make for a great tutorial.
At the end, we will have a fully functioning REST service for serving basic CRUD on a user resource. Code for the final product can be found here.
This is the second in a series of high level tutorials on setting up a REST service using golang, focused on using public packages and defining a basic webserver.

## Packages
Go allows for external packages to be accessed in your source via the import statement. These packages must be accessible somewhere in your GOPATH or GOROOT environment, depending on whether the packages come from a third party or the Go standard library, respectively. In the previous tutorial, you may have noticed that we imported the fmt package when we tested our setup. The fmt package is part of the standard library. You can find more standard library packages here.
Third party packages can be fetched via the go get command. There are a few special rules for fetching a remote package, described here, but for now we just need to understand the rules for packages publicly available on github. Put simply, any go package on github can fetched via the following.

```
$ go get github.com/USER/PROJECT
```

You'll notice there is no version information in the above command, which is another caveat of go. There are several solutions for this, such as using a makefile to set a temporary GOPATH within the application and checking dependencies into your source control. I've also grown fond of GoDep.

### Router
Since we're trying to build out a webserver, let's find a router package. When searching for a go package to use, GoDoc is always a good place to start. There are many excellect candidates, but let's use the httprouter package. We first need to fetch the package with go get.

```
$ go get github.com/julienschmidt/httprouter
```

To get started, we first need to create a directory to host our application. If you followed along with the last tutorial, the following will get us started. Otherwise, setup a new directory somewhere in your GOPATH.

```
$ mkdir -p ~/src/go/src/github.com/swhite24/go-rest-tutorial
```

Now for some code. Let's add a server.go to bootstrap our app.

```golang
package main

import (
	// Standard library packages
	"fmt"
	"net/http"

	// Third party packages
	"github.com/julienschmidt/httprouter"
)

func main() {
	// Instantiate a new router
	r := httprouter.New()

	// Add a handler on /test
	r.GET("/test", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// Simply write some test data for now
		fmt.Fprint(w, "Welcome!\n")
	})

	// Fire up the server
	http.ListenAndServe("localhost:3000", r)
}

Let's break down the main function.

r := httprouter.New()
Here we are instantiating a new httprouter instance. Documentation for public packages can be found at GoDoc, with the documentation for this package available here.

r.GET("/test", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Simply write some test data for now
	fmt.Fprint(w, "Welcome!\n")
})
```


Here we are adding a new handler to respond to the "/test" route. Notice that the handler itself has a different signature than the standard http.HandlerFunc.
Also, we are able to respond to the request with fmt.Fprint. Strangely, if you look at the signature of fmt.Fprint, it shows the following.

```
func Fprint(w io.Writer, a ...interface{}) (n int, err error)
```

So how are we able to pass our instance of http.ResponseWriter to it as an io.Writer? Because io.Writer is an interface. Interfaces are implemented implicitly in go. In fact, http.ResponseWriter itself is an interface, which also satisfies the io.Writer interface.

And for the last interesting piece.

```
http.ListenAndServe("localhost:3000", r)
```

Here, as you can guess, we are firing up a webserver to listen on localhost:3000, using our router to handle requests. This is another instance where we can see the power of interfaces in go. http.ListenAndServe accepts an interface type http.Handler, which the author of httprouter provides as a convencience.
Let's go ahead and test our example server. First, fire it up.

```
$ go run server.go
```

Then, in another terminal, hit our router with a curl statement.

```
$ curl http://localhost:3000/test
```

### Models

When building webapps in go, I tend to modularize the components as much as possible, both for testing purposes and readability. Models are one example of a piece of functionality I usually break into a separate package.
First, we need to create a directory within our application for the new package. Since I want it to be referred to as models, I'll name the directory the same. This is not required, but I find it to be a good convention.

$ mkdir models
Now let's add a file in the models directory to represent the structure of our user resource. I'll call it something crazy like user.go.

```
package models

type (
	// User represents the structure of our resource
	User struct {
		Name   string
		Gender string
		Age    int
        Id     string
	}
)
```

Simple right? The only thing in this file is the definition of a user. A more complex resource may contain additional methods our references to other types.
Now, let's stub our router with a new handler function to handle retrieving a user. Update server.go to the following.

```golang
package main

import (
	// Standard library packages
	"encoding/json"
	"fmt"
	"net/http"

	// Third party packages
	"github.com/julienschmidt/httprouter"
	"github.com/swhite24/go-rest-tutorial/models"
)

func main() {
	// Instantiate a new router
	r := httprouter.New()

	// Get a user resource
	r.GET("/user/:id", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// Stub an example user
		u := models.User{
			Name:   "Bob Smith",
			Gender: "male",
			Age:    50,
			Id:     p.ByName("id"),
		}

		// Marshal provided interface into JSON structure
		uj, _ := json.Marshal(u)

		// Write content-type, statuscode, payload
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprintf(w, "%s", uj)
	})

	// Fire up the server
	http.ListenAndServe("localhost:3000", r)
}
```

You'll notice two major changes. First, we added an import statement for our newly created models package. Second, we added a new handler to respond to "/user/:id" for retrieving a user. For those unfamiliar with this pattern, :id represents a parameter we can retrieve from the path, which we do to populate the Id field of our user. We also call json.Marshal on our user, which will decode our user into a JSON representation, and deliver that to the client.

### Restart the server and test it out with a curl statement.

```
$ curl http://localhost:3000/user/1
{"Name":"Bob Smith","Gender":"male","Age":50,"Id":"1"}
```

Sweet! We have our example user. But do we really want to deliver the user with capitalized field names? Probably not. We could try to use lower case field names on the user struct, but then the fields will no longer be exported and available in our main package.
The documentation for the json package tells us that we can "alias" field names to be whatever we want using struct tags.

Let's update our user.go file to use lower case names when delivering json.

```golang
package models

type (
	// User represents the structure of our resource
	User struct {
		Name   string `json:"name"`
		Gender string `json:"gender"`
		Age    int    `json:"age"`
		Id     string `json:"id"`
	}
)
```

Now, let's restart the server and test the route again.

```
$ curl http://localhost:3000/user/1
{"name":"Bob Smith","gender":"male","age":50,"id":"1"}
```

### Much better.

Finally, if we stub out the remaining routes for our user resource we will have something like this in server.go.

```golang
r.POST("/user", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Stub an user to be populated from the body
	u := models.User{}

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&u)

	// Add an Id
	u.Id = "foo"

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
})

r.DELETE("/user/:id", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// TODO: only write status for now
	w.WriteHeader(200)
})
```

###  Controllers
One thing you may begin to notice is that our server.go is getting rather bloated with handlers. This is another piece of functionality I typically refactor into a package. Let's do that into a package called controllers.
First, make the directory.

```
$ mkdir controllers
```

Then, add a user.go into the directory.

```golang
package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/swhite24/go-rest-tutorial/models"
)

type (
	// UserController represents the controller for operating on the User resource
	UserController struct{}
)

func NewUserController() *UserController {
	return &UserController{}
}

// GetUser retrieves an individual user resource
func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Stub an example user
	u := models.User{
		Name:   "Bob Smith",
		Gender: "male",
		Age:    50,
		Id:     p.ByName("id"),
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

// CreateUser creates a new user resource
func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Stub an user to be populated from the body
	u := models.User{}

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&u)

	// Add an Id
	u.Id = "foo"

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
}

// RemoveUser removes an existing user resource
func (uc UserController) RemoveUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// TODO: only write status for now
	w.WriteHeader(200)
}
```

You'll notice that we simply made an exact copy of the handlers previously found in server.go and moved them to individual methods of our user controller type.
An updated server.go to use our new controller looks like the following.

```golang
package main

import (
	// Standard library packages
	"net/http"
	// Third party packages
	"github.com/julienschmidt/httprouter"
	"github.com/swhite24/go-rest-tutorial/controllers"
)

func main() {
	// Instantiate a new router
	r := httprouter.New()

	// Get a UserController instance
	uc := controllers.NewUserController()

	// Get a user resource
	r.GET("/user/:id", uc.GetUser)

	r.POST("/user", uc.CreateUser)

	r.DELETE("/user/:id", uc.RemoveUser)

	// Fire up the server
	http.ListenAndServe("localhost:3000", r)
}
```


Much cleaner! Let's run another curl to verify everything is still working.

```

$ curl http://localhost:3000/user/1
{"name":"Bob Smith","gender":"male","age":50,"id":"1"}
```

### Fin
That's the end of the tutorial on setting up a basic webserver with golang. We've created an httprouter instance, a model for our user resource, a controller for our user resource, and finally wired them all together.
Next time, we'll add a backend to get a true webservice feel. Check it out.
