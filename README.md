# Green Light Backend

Following the book, `"Let's Go Further", Alex Edwards`.

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->

**Table of Contents**

- [Green Light Backend](#green-light-backend)
- [Health check API (`System Structuring`, 221118105243-03)](#health-check-api-system-structuring-221118105243-03)
- [Clean code practices (221118104844-03)](#clean-code-practices-221118104844-03)

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

Encapsulation of API routing. We will allocate a file for routing
`/cmd/api/routes.go`. So, `/cmd/api/main.go` only has the essential code to
launch the application.
