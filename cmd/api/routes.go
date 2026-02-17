package main

func (app *application) routes() {
	_ = app.router.Group("/v1")
}
