# GoWorkoutTracker
### Description
---
A Frontend Masters project resembling the backend of a workout tracker application written in go

### Arhitecture
---
The entire application will be running using an app object that is parsed and attributed  during the main function

main.go
- configuration and lifecycle of the application
	-  we instantiate the app object
	- instantiate the server using the router as the handler
	- we send messages through the logger to the terminal
	- we listen for errors sent in the handlers  and break on found
	- we end the database connection at the end of execution

app.go
- main object of the application which will contain 
	- model handlers and stores
	- logging of the application
	- database connection
routes.go
- router of the application which instantiates the handlers as routes on the app object using the "chi" package
tokens.go
- handling of creating tokens and handling token scopes
middleware.go
- handling of authentification middleware by binding and checking for the token in the context of requests
fs.go
- file that will tell the compiler the order in which to run the migrations

### Functionality
---
logging 
- is done using the log package
response formatting 
- is done using the utils functions
	-  set headers for status and content type
	- format data as proper json using json.MarshalIndent
error logging
- is done by sending messages through the logger to the console and sending discrete error reponses because the client does not need to know what happens in the backend 
authentification
- statefull token auth using middlewares
- password encryption
	-  done through specific methods using the bcrypt and crypto/sha256 packages
- tokens
	- token is created by allocating random memory bytes, applying, base 32 encoding, and generating a sha256 checksum hash
	- entries are created in the database containing userd id, expiration time, and scope and specific token hash checksum
testing
- table tests of database manipulation methods using stretchr/testify  package
### File structure
---
- go.mod
- go.sum
- main.go
- internal
	- api
		- handlers(methods called on routes) of the application 
	- app
		- app.go
	- routes
		- routes.go
	- store
		- methods for comunicating with the database
	- utils
		- methods and environment variables used throughout the system
	- middleware
		- handling of authentification middleware 
- migrations
	- migrations used to create the database
- public
	- a list of example requests can be used to test the functionality of the project

### Database
---
- postgres database
- we use pgx as the driver to create the postgres database
- goose for database migrations
- using embedded migrations

### Improvements:
---
- constants.go file for defining error messages and print statements that are going to be reused
- abstract token creation to pass the scope in the request body
- avoiding code duplication in error checking in handler sections
- add functionality to remove old tokens in a cron based functionality
