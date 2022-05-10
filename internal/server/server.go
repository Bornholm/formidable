package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"forge.cadoles.com/wpetit/formidable/internal/server/template"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

type Server struct {
	host     string
	port     uint
	schema   *jsonschema.Schema
	defaults interface{}
	values   interface{}
	onUpdate OnUpdateFunc
}

type OnUpdateFunc func(values interface{}) error

func (s *Server) Start(ctx context.Context) (<-chan net.Addr, <-chan error) {
	errs := make(chan error)
	addrs := make(chan net.Addr)

	go s.run(ctx, addrs, errs)

	return addrs, errs
}

func (s *Server) run(parentCtx context.Context, addrs chan net.Addr, errs chan error) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		errs <- errors.WithStack(err)

		return
	}

	addrs <- listener.Addr()

	defer func() {
		if err := listener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			errs <- errors.WithStack(err)
		}

		close(errs)
		close(addrs)
	}()

	go func() {
		<-ctx.Done()

		if err := listener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			log.Printf("%+v", errors.WithStack(err))
		}
	}()

	templates := getEmbeddedTemplates()

	if err := template.Load(templates, "template"); err != nil {
		errs <- errors.WithStack(err)

		return
	}

	assets := getEmbeddedAssets()
	assetsHandler := http.FileServer(http.FS(assets))

	router := chi.NewRouter()

	router.Get("/", s.serveFormReq)
	router.Post("/", s.handleFormReq)
	router.Handle("/assets/*", assetsHandler)

	log.Println("http server listening")

	if err := http.Serve(listener, router); err != nil && !errors.Is(err, net.ErrClosed) {
		errs <- errors.WithStack(err)
	}

	log.Println("http server exiting")
}

func New(funcs ...OptionFunc) *Server {
	opt := defaultOption()
	for _, fn := range funcs {
		fn(opt)
	}

	return &Server{
		host:     opt.Host,
		port:     opt.Port,
		schema:   opt.Schema,
		defaults: opt.Defaults,
		values:   opt.Values,
		onUpdate: opt.OnUpdate,
	}
}
