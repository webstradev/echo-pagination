# echo-pagination
[![Run Tests](https://github.com/webstradev/echo-pagination/actions/workflows/test.yml/badge.svg)](https://github.com/webstradev/echo-pagination/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/webstradev/echo-pagination/branch/master/graph/badge.svg?token=C2D4QHYHI4)](https://codecov.io/gh/webstradev/echo-pagination)
[![Go Reference](https://pkg.go.dev/badge/github.com/webstradev/echo-pagination.svg)](https://pkg.go.dev/github.com/webstradev/echo-pagination)
[![Go Report Card](https://goreportcard.com/badge/github.com/webstradev/echo-pagination)](https://goreportcard.com/report/github.com/webstradev/echo-pagination)
[![CodeQL](https://github.com/webstradev/echo-pagination/actions/workflows/codeql.yml/badge.svg)](https://github.com/webstradev/echo-pagination/actions/workflows/codeql.yml)

Simple pagination middleware for the echo framework. Allows for the usage of url parameters like `?page=1&size=25` to paginate data on your API.

# Installation
```bash
$ go get github.com/webstradev/echo-pagination
```

# Default Usage
This package comes with various default options that are configurable using functional options. By default the paginator will use the query parameters `page` and `size` with values of `1` and `10` and a maximum page size of `100`.

### Using the middleware on a router will apply it to all requests on that router:
```go
package main

import (
  "net/http"

  "github.com/labstack/echo/v4"
  "github.com/webstradev/echo-pagination/pkg/pagination"
)

func main(){
  e := echo.New()
  e.Use(pagination.New())

  e.GET("/hello", func(c echo.Context) error {
    return c.Status(http.StatusOK)  
  })
  
	e.Logger.Fatal(e.Start(":1323"))
}
```

#### Using the middleware on a single route will only apply it to that route:
```go
package main

import (
  "net/http"

  "github.com/labstack/echo/v4"
  "github.com/webstradev/echo-pagination/pkg/pagination"
)

func main(){
  e := echo.New()
  
  e.GET("/hello", func(c echo.Context){
    page := c.GetInt("page")
  
		return c.String(http.StatusOK, "Hello, World!")
  }, pagination.New())
  
	e.Logger.Fatal(e.Start(":1323"))
}
```
The `page` and `size` are now available in the echo context of a request and can be used to paginate your data (for example in an SQL query).

## Custom Usage
To create a pagination middleware with custom parameters the New() function supports various custom options provided as functions that overwrite the default value.
All the options can be seen in the example below.
```go
package main

import (
  "net/http"
  
  "github.com/labstack/echo/v4"
  "github.com/webstradev/echo-pagination/pkg/pagination"
)

func main(){
  e := echo.New()
  
  paginator := pagination.New(
    pagination.WithPageText("page"), 
    pagination.WithSizeText("rowsPerPage"),
    pagination.WithDefaultPage(1),
    pagination.WithDefaultPageSize(15),
    pagination.WithMinPageSize(5),
    pagination.WithMaxPageSize(15),
  )
  
  e.GET("/hello", func(c echo.Context){
    return c.String(http.StatusOK, "Hello, World!")
  }, paginator)
  
	e.Logger.Fatal(e.Start(":1323"))
}
```

The custom middleware can also be used on an entire router object similarly to the first example fo the Default Usage.
