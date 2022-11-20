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
- [`JSON` responses/encoding (221118144752-03)](#json-responsesencoding-221118144752-03)
- [Clean Code, `DRY` (221118145846-03)](#clean-code-dry-221118145846-03)
- [Change encoding of `structs` (221118162000-03)](#change-encoding-of-structs-221118162000-03)
- [Optional directives (221118162809-03)](#optional-directives-221118162809-03)
- [`Enveloping` - Make every response a `higher-level object` (221118165655-03)](#enveloping---make-every-response-a-higher-level-object-221118165655-03)
- [TODO Format the `JSON` arbitrarily](#todo-format-the-json-arbitrarily)
- [`JSON` requests/decoding (221119075356-03)](#json-requestsdecoding-221119075356-03)
- [`createMovieHandler` - business logic (221119075724-03)](#createmoviehandler---business-logic-221119075724-03)
- [`POST` example using `curl` (221119075625-03)](#post-example-using-curl-221119075625-03)
- [Error handling - `readJSON` method (221119092520-03)](#error-handling---readjson-method-221119092520-03)
- [Custom JSON Decoding - Chap 4.4 (221119093555-03)](#custom-json-decoding---chap-44-221119093555-03)
- [`json.Unmarshaler` interface -- The solution (221120114402-03)](#jsonunmarshaler-interface----the-solution-221120114402-03)
- [`UnmarshalJSON` method implementation for `Runtime` (221120120058-03)](#unmarshaljson-method-implementation-for-runtime-221120120058-03)
- [`CURL` example (221120122039-03)](#curl-example-221120122039-03)
- [Resources](#resources)

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
requests are only `decoded`, by using the same format.

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

## Validating `JSON` Input (221120140913-03)

We will perform validation of the _JSON_ `POST` request, following certain
business-rules restrictions:

> - The movie title provided by the client is not empty and is not more than 500 bytes long.
> - The movie year is not empty and is between 1888 and the current year.
> - The movie runtime is not empty and is a positive integer.
> - The movie has between one and five (unique) genres.

> If any of those checks fail, we want to send the client a `422 Unprocessable Entity ` response along with error messages which clearly describe the
> validation failures.
> (page 96)

### Custom validator package (221120150759-03)

In order to perform a list of restrictions on the JSON POST request, we could create a variable that consists of a list of key-value pair of errors. Because, a request could have none, one or many errors, related to our restrictive rules.

We will implement these functions, and other `regex`-related help functions that help us accomplish this _validation_.

In `internal/validator/validator.go`,

```go
package validator

import "regexp"

// from https://html.spec.whatwg.org/#valid-e-mail-address; adding '\' for each isolated "\".
var (
	EailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type Validator struct {
	Errors map[string]string
}

// Instantiate an empty Validator ("error list")
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// If there is no errors, then it's a valid request
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// if v.Errors[key] doesn't exist yet, add message-value to the v.Errors[key]
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
```

#### Implement new error message, related to `validator` (221120153026-03)

In `cmd/api/errors.go`, add `422 Unprocessable Entity` type of error, for our application.

```go
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
```

In `cmd/api/movies.go`, add the following checks `createMovieHandler`,

```go
v := validator.New()
// Check if title is not empty. If it is empty (check is False), add title-type error
v.Check(!(input.Title == ""), "title", "Must be provided")
// Check if title is not less then 500 bytes. If it is greater (check is False), add title-type error
v.Check(!(len(input.Title) <= 500), "title", "Must not be more than 500 bytes long")

v.Check(!(input.Runtime < 0), "runtime", "Must be a positive integer")
v.Check(!(input.Runtime == 0), "runtime", "Must be provided")

// Check if year value is not empty
v.Check(!(input.Year == 0), "year", "Must be provided")
// Check if year is not lesser than 1888
v.Check(!(input.Year < 1888), "year", "Must be greater than 1888")
// Check if movie-year is not a value greater than current year it's being added
v.Check(!(input.Year >= int32(time.Now().Year())), "year", "Must not be in the future")

v.Check(!(input.Genres == nil), "genres", "Must be provided")
v.Check(!(len(input.Genres) <= 1), "genres", "Must contain at lest one genre")
v.Check(!(len(input.Genres) >= 5), "genres", "Must not contain more than 5 genres")
v.Check(validator.Unique(input.Genres), "genres", "Must not contain duplicates")

if !v.Valid() {
	app.failedValidationResponse(w, r, v.Errors)
	return
}
```

# Resources

- https://html.spec.whatwg.org/#valid-e-mail-address (spec for e-mail validation
  regex)
- https://github.com/xescugc/marshaler (Custom-type URL marsheler implementation)
