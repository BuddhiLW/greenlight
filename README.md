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
- [`JSON` responses (221118144752-03)](#json-responses-221118144752-03)
- [Clean Code, `DRY` (221118145846-03)](#clean-code-dry-221118145846-03)
- [Change encoding of `structs` (221118162000-03)](#change-encoding-of-structs-221118162000-03)
- [Optional directives (221118162809-03)](#optional-directives-221118162809-03)
- [`Enveloping` - Make every response a `higher-level object` (221118165655-03)](#enveloping---make-every-response-a-higher-level-object-221118165655-03)
- [TODO Format the `JSON` arbitrarily](#todo-format-the-json-arbitrarily)

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

## `JSON` responses/encoding (221118144752-03)

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

Also, this helper function let us add headers to your JSON encoding process.

At `helpers.go`,

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

Rewriting `healthcheck.go`,

```go
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server couldn't complete the request", http.StatusInternalServerError)
	}
}
```

OBS: `json.MarshalIndent(v interface{}, "", "\t")` will _pretty-print_ the
response. **BUT**, at the cost that it will take 65% more time to run, and 30%
more memory, than `json.Marshal()`.

### Change encoding of `structs` (221118162000-03)

Encoding `structs` is pretty straight forward. Furthermore, we can change the
way the resulting encoding appears in `JSON`. For example, if you want to
conform to _snake-case_.

At, `internal/data/movies.go`,

```go
type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Year      int32     `json:"year"`
	Runtime   int32     `json:"runtime"`
	Genres    []string  `json:"genres"`
	Version   int32     `json:"version"`
}
```

At `cmd/api/movies.go`,

```go
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, movie, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server couldn't proceed with the request", http.StatusInternalServerError)
	}

}
```

#### Optional directives (221118162809-03)

One can specify if the field should only be encoded if it's _non-empty_, or if
it shouldn't be encoded at all.

- `-`: use `json:"-"` to never encode field.
- `omitempty`: use `json:"field_name, omitempty"` to encode only when _non-empty_.
- `string`: use `json:"field_name, string"` to enforce encoding as _string_.

OBS:

> Hint: If you want to use omitempty and not change the key name then you can leave it
> blank in the struct tag — like this: json:",omitempty" . Notice that the leading comma
> is still required. (page 51)

### `Enveloping` - Make every response a `higher-level object` (221118165655-03)

We will have to change `writeJSON` at `/cmd/api/helpers.go`; `showMovieHandler`,
at `/cmd/api/movies.go` and `healthcheckHandler` at `/cmd/api/healthcheck.go`.
They all must conform to the new `whiteJSON` specification.

First, we create a `envelope` _type_.

At `helpers.go`,

```go
type envelope map[string]interface{}
```

And, change `writeJSON` to receive a `envelope` type, (`data envelope`).

```go
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
```

Then, change `showMovieHandler`, because `writeJSON` must receive an `envelope` _type_.

```go
envel := envelope{"movie": movie}
err = app.writeJSON(w, http.StatusOK, envel, nil)
```

Finally, change `healthcheckHandler`,

```go
envel := envelope{
	"status": "available",
	"system_info": map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	},
}

err := app.writeJSON(w, http.StatusOK, envel, nil)
```

### TODO Format the `JSON` arbitrarily

## `JSON` requests/decoding (221119075356-03)

### `createMovieHandler` - business logic (221119075724-03)

Using `Decode` onto `json.NewDecoder`.

`r.Body` (**JSON**) gets _decoded_ into a Go data-structure (**input-struct**), and stored inside the `input` variable.

```go
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}
```

- What does it mean to use `%+v` as the printing format? (221119080509-03)

  > Printing. The verbs: General: `%v` the value in a default
  > format when printing structs, the plus flag (`%+v`) adds field
  > names `%#v` a Go-syntax representation of the value `%T`

But, as we want to error-handle, the function will become,

```go
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}
```

`readJSON` not only will decode, but will also error-handle the application, regarding all expected production-errors possible.

### `POST` example using `curl` (221119075625-03)

```sh
$ curl -d '{"title":"Moana","runtime":107, "genres":["animation","adventure"]}' http://localhost:4000/v1/movies
```

```sh
{Title:Moana Year:0 Runtime:107 Genres:[animation adventure]}
```

### Error handling - `readJSON` method (221119092520-03)

This will ensure `resiliance` and `uniformity` to our `POST` request errors.

The list of **Errors** we will friendly-handle:

- Maximum size of 1 MB (security related: `DoS`).
- Tell the user there are unknown-fields being sent to us, which we refuse to
  handle.
- Syntax Error.
- Field-type Error.
- Empty `POST` request.
- Any other invalid requests will return panic.
- Any other kind of error, if not the above, will send the default
  Go-implemented error for that case.

```go
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {

	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formated JSON (at character %d).",
				syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("Body contains badly-formated JSON.")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("Body contains incorrect JSON type for field %q.",
					unmarshalTypeError.Field)
			}
			return fmt.Errorf("Body contains incorrect JSON type (at chracter %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("Body must not be empty.")

		// regex that the error is of type "json: unkown field"
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			// remove the "json: unkown field" sentence, from the error - remains the field name.
			fieldName := strings.TrimPrefix(err.Error(), "json: unkown field")
			return fmt.Errorf("Body contains unkown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("Body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}
```

### Custom JSON Decoding - Chap 4.4 (221119093555-03)

Conversely to `JSON encoding` _runtime_ as `"runtime": "<runtime> mins"` - which
we already did at the `Format the JSON arbitrarily` topic - we can enforce that
requests are only accepted using the same format.

Currently, we would get:

```sh
curl -d '{"title": "Moana", "runtime": "107 mins"}' localhost:4000/v1/movies
```

```sh
{
        "error": "Body contains incorrect JSON type for field \"runtime\"."
}
```

> To make this work, what we need to do is _intercept the decoding process_ and
> manually convert the `"<runtime> mins"` JSON string into an `int32` instead.
> (Alex Edwards, page 92)

### `json.Unmarshaler` interface -- The solution (221120114402-03)

> type Unmarshaler ¶

```go
 type Unmarshaler interface {
	UnmarshalJSON([]byte) error
 }
```

> Unmarshaler is the interface implemented by types that can unmarshal
> a JSON description of themselves. The input can be assumed to be a
> valid encoding of a JSON value. UnmarshalJSON must copy the JSON
> data if it wishes to retain the data after returning.
>
> By convention, to approximate the behavior of Unmarshal itself,
> Unmarshalers implement UnmarshalJSON([]byte("null")) as a no-op.
> [https://pkg.go.dev/encoding/json#Unmarshaler](https://pkg.go.dev/encoding/json#Unmarshaler)

If the a decoding `destination type` implements a `UnmarshalJSON` method, then
when the decoding of this type will use this method.

So, all we got to do is implement the `UnmarshalJSON` method for the `Runtime` type.

#### `UnmarshalJSON` method implementation for `Runtime` (221120120058-03)

NOTE that once we wrote, at `internal/data/runtime.go`,

```go
type Runtime int32
```

The `Runtime` custom type now behaves as the underlying `int32` type, except
it's hard-coded methods - namely, `UnmarshalJSON` - can be re-implemented.

So, we go to `cmd/api/movies.go`, as make the `input` struct have a `Runtime`
field of type `data.Runtime`.

```go
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"` // Make this field a data.Runtime type.
		Genres  []string     `json:"genres"`
	}

	...
}
```

And, we re-implement the `UnmarshalJSON` method for `Runtime`, at
`internal/data/runtime.go`.

**GOAL:** transform _<runtime> mins_ `string`-type into a `Runtime` type and
ship it to our `r` `Runtime` type variable; _<runtime> mins_ comes as a
`[]byte`, which will be our `jsonValue` variable argument.

`UnmarshalJSON: ("<runtime> mins" []byte) -> (<runtime> Runtime)`

```go
var ErrInvalidRuntimeFormat = errors.New("Invalid 'runtime' format, should be '<runtime> mins'")

func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {

	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	parts := strings.Split(unquotedJSONValue, " ")
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	*r = Runtime(i)

	return nil
}
```

#### `CURL` example (221120122039-03)

> You should see that the request completes successfully, and the number is
> extracted from the string and assigned the `Runtime` field of our `input`
> struct. (page 94-95)

```sh
curl -d '{"title": "Moana", "runtime": "107 mins"}' localhost:4000/v1/movies
```

```sh
{Title:Moana Year:0 Runtime:107 Genres:[]}
```

```sh
curl -d '{"title": "Moana", "runtime": 107}' localhost:4000/v1/movies
```

```sh
{
        "error": "Invalid 'runtime' format, should be '_runtime_ mins'"
}
```

```sh
curl -d '{"title": "Moana", "runtime": "107 minutes"}' localhost:4000/v1/movies
```

```sh
{
        "error": "Invalid 'runtime' format, should be '_runtime_ mins'"
}
```
