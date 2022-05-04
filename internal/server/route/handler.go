package route

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

func NewHandler(schema *jsonschema.Schema, defaults, values interface{}, assetsHandler http.Handler) (*chi.Mux, error) {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	// router.Use(middleware.Logger)

	router.Get("/", createRenderFormHandlerFunc(schema, defaults, values))
	router.Post("/", createHandleFormHandlerFunc(schema, defaults, values))

	router.Handle("/assets/*", assetsHandler)

	return router, nil
}
