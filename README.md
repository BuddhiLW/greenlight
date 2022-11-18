# Green Light Backend

Following the book, `"Let's Go Further", Alex Edwards`.

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->

**Table of Contents**

- [Green Light Backend](#green-light-backend)
- [Health check API (`System Structuring`, 221118105243-03)](#health-check-api-system-structuring-221118105243-03)
- [Clean code practices (221118104844-03)](#clean-code-practices-221118104844-03)
- [Movies, `GET`, `POST`, `PUT`, `DELETE` (`CRUD`, 221118105749-03)](#movies-get-post-put-delete-crud-221118105749-03)
- [Clean Code Practices (221118111017-03)](#clean-code-practices-221118111017-03)
- [Encapsulation](#encapsulation)
- [**DRY** principle (221118123928-03)](#dry-principle-221118123928-03)
- [`JSON` responses](#json-responses)

<!-- markdown-toc end -->

## Health check API (`System Structuring`, 221118105243-03)

The **Architectural Design** of the `Process` for the _API_ endpoint goes as follows:

| URL Pattern       | Handler            | Action                         |
| ----------------- | ------------------ | ------------------------------ |
| /v1/healthcheck   | healthcheckHandler | Show application information.  |
| ----------------- | ------------------ | ------------------------------ |

### Clean code practices (221118104844-03)

- The `healthcheckHandler` is implemented as a _Method_ to the _application_
  type. This way, the handler can access values instantiated inside the
  application, in `main()` (at main.go).

```go
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environmnet: %s\n", app.config.env)
	fmt.Fprintf(w, "version: %s\n", version)
}
```

Note that we are making use of `app.config.env` and `version`, which are values
inside our `application`. Thus, we avoid `global variables` and `closures`.

## Movies, `GET`, `POST`, `PUT`, `DELETE` (`CRUD`, 221118105749-03)

| Method     | URL Pattern       | Handler                | Action                          |
| ---------- | ----------------- | ---------------------- | ------------------------------- |
| `GET`      | /v1/healthcheck   | `healthcheckHandler`   | Show app info                   |
| `GET`      | /v1/movies        | `listMoviesHandler`    | Show a list of movies infos     |
| `POST`     | /v1/movies        | `createMovieHandler`   | Create a new movie              |
| `GET`      | /v1/movies/:id    | `showMovieHandler`     | Show info of a specific movie   |
| `PUT`      | /v1/movies/:id    | `editMovieHandler`     | Update info about movie         |
| `DELETE`   | /v1/movies/:id    | `deleteMovieHandler`   | Delete a specific movie         |
| ---------- | ----------------- | ---------------------- | ------------------------------- |

### Clean Code Practices (221118111017-03)

#### Encapsulation

Encapsulation of API routing. We will allocate a file for routing
`/cmd/api/routes.go`. So, `/cmd/api/main.go` only has the essential code to
launch the application.

At `/cmd/api/routes.go`,

```go
func (app *application) routes() *httprouter.Router {

	r := httprouter.New()

	r.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	r.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
	r.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)

	return r
}
```

At `/cmd/api/main.go`, we call the created handler for our server.

```go
srv := &http.Server{
    Addr:         fmt.Sprintf(":%d", cfg.port),
    Handler:      app.routes(),
    IdleTimeout:  time.Minute,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 30 * time.Second,
}
```

#### **DRY** principle (221118123928-03)

`DRY: Don't Repeat Yourself `

We will create a _helper function_ to read the parameter, from the URL, from an
incoming request.

```go
func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("Invalid id parameter. Must be an integer, greater than zero")
	}

	return id, nil
}
```

This way, `/cmd/api/movies.go` can be written as:

```go
package main

import (
	"fmt"
	"net/http"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "show the details of movie %d\n", id)
}
```

## `JSON` responses (221118144752-03)

- `JSON`: `JavaScript Object Notation`, is a pure-text structured data. It's a map
  of key-value pairs. Each object is delimited, in between brackets `{}`. The
  values of each key can be numbers, strings, objects or vectors.

- _JSON_ responses can be written with `Sprintf` or any method that outputs a
  `string` type.
- Generally, a powerful and simple tool to convert **Go** _data-types_ into
  _JSON_ is the `func json.Mashal(v interface{}) ([]byte, error)` which output a
  _JSON_-formated string.

OBS: `interface{}` means any `type`, including `constructs`.

### Clean Code, `DRY` (221118145846-03)

It's extremely common to _Marshal_ and _Unmarshal_ Go data-structure and JSON.
So, to insure we `DRY`, we will write helper functions.

```go
func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
```
