# Go CRM Backend

Implementation of a basic Go web server emulating the backend of a simple CRM application written as a final project for the [Udacity Golang Course](https://www.udacity.com/course/golang--cd11970).

## Features

The project consists of a web server written in the **Go** language which after started serves a REST API to which users can make HTTP requests using utilities like [Postman](https://www.postman.com/) or [cURL](https://curl.se/). The application provides the following five endpoints to which users can make CRUD operations to the CRM application:

| Operation                             | Endpoint          |
| ------------------------------------- | ----------------- |
| Creating a customer                   | `/customers`      |
| Getting the data of a single customer | `/customers/{id}` |
| Getting all customers                 | `/customers`      |
| Updating a customer                   | `/customers/{id}` |
| Deleting a customer                   | `/customers/{id}` |

These CRUD operations interact with the CRM application in order to create, retrieve, update and delete customers from the CRM, and each of these endpoints is assigned to its appropriate HTTP method (GET, POST, PUT, DELETE) by leveraging a custom HTTP method-based router.

The application runs *"in-memory"* and does not implement any mechanism for persisting the customers data after the server is stopped. Internally, the application stores the customers data in a Go `map`, which serves as the database of the CRM.

As a note, the application just serves a backend API and does not have any frontend other than a static welcome page located at the home route of the server.

## Requirements to run the application

In order to run the application, the Go language needs to be installed. The application was written in the version **1.19.3** of Go, but it should work with any recent version of the language. You can download Go [here](https://go.dev/doc/install).

Also, a code editor from which to compile and run the application is required. [Visual Studio Code](https://code.visualstudio.com/) together with the [Go extension](https://code.visualstudio.com/docs/languages/go) is a good option, but feel free to use your favorite editor. And in order to interact with the application as a user, you will need any utility software that allows you to make requests to a REST API, [Postman](https://www.postman.com/) or [cURL](https://curl.se/) are good options.

In order for the application to compile, the two following third-party Go packages need to be installed:

- [**gorilla/mux**](https://pkg.go.dev/github.com/gorilla/mux) package, which provides a custom HTTP router (or mux) able to perform HTTP method-based routing.
- Google's [**uuid**](https://pkg.go.dev/github.com/google/uuid#section-readme) package, required to provide each of the CRM customers with an UUID.

To install them, open a terminal, either from the OS you are using or Visual Studio Code's terminal (if you are using this IDE) and navigate to the root directory of the CRM backend project, then write the `go get` command followed by the URL of the package you want to install and hit enter. This will install the packages and make them available to be used in the application code:

```
go get github.com/gorilla/mux
go get github.com/google/uuid
```

## Starting and using the application

To start the application, open a terminal and navigate to the root directory of the CRM project, where the `main.go` file is located. Then, run the `go run main.go` command. This will start the server and make it start to listen for requests at port `3000`, which is the default (hardcoded) port for the server, although you can can change it in the `main.go` file to any other port if required.  

After the server is started, the application can be accessed through http://localhost:3000 and REST requests can start to be made to its different API endpoints. As stated before, the home route of the server will serve an static HTML welcome page. In this page you can view the documentation of the REST API, which specifies the format of the JSON payloads that the API expects to receive and the way in which to pass URL variables for each of the API endpoints.

## Unit Tests

Located at the root directory of the CRM project is the `main_test.go` file. This is a Go file containing basic unit tests for the CRM application, which just test that correct HTTP status codes and Content-Types are returned by the different API endpoints of the application in the server responses. This file is not of my authorship and was downloaded from the *Resources* section of the Udacity Golang Course.