package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

// Rest - Rest is a rest server object
type Rest struct {
	Version string
	Port    int

	Ready *atomic.Value
	Vault *Vault

	httpServer *http.Server
}

// Run - Run start a rest server
func (s *Rest) Run(ctx context.Context) error {
	router := s.routes()

	s.httpServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", s.Port),
		Handler:           router,
		ReadHeaderTimeout: 60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Printf("rest server started on port %d", s.Port)

	s.Ready.Store(true)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown - Shutdown will stop a rest server
func (s *Rest) Shutdown(ctx context.Context) error {
	log.Print("shutdown rest server")
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			return err
		}
	}
	return nil
}

// routes - routes generating routes for rest server
func (s *Rest) routes() chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.Throttle(1000))
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(middleware.Compress(6, "gzip"))
	router.Use(middleware.GetHead)
	router.Use(middleware.Heartbeat("/ping"))
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Recoverer)

	corsOptions := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE"},
		AllowedHeaders: []string{"Accept", "Content-Type"},
		MaxAge:         300,
	})
	router.Use(corsOptions.Handler)

	router.Use(Healthz)
	router.Use(Readyz(s.Ready))
	router.NotFound(NotFound)

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "dist")
	s.FileServer(router, "/", http.Dir(filesDir))
	router.Get("/index.html", func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile(path.Join(filesDir, "index.html"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		render.HTML(w, r, string(data))
	})

	router.Get("/namespaces", s.Vault.NamespacesList)

	router.Route("/vault", func(vaultRoute chi.Router) {
		vaultRoute.Get("/{namespace}", s.Vault.Get)
		vaultRoute.Get("/{namespace}/{secret}", s.Vault.Find)
		vaultRoute.Post("/{namespace}/{secret}", s.Vault.Add)
		vaultRoute.Delete("/{namespace}/{secret}", s.Vault.Delete)
	})

	return router
}

// FileServer - serving static files
func (s *Rest) FileServer(r chi.Router, path string, root http.FileSystem) {
	origPath := path

	fs := http.StripPrefix(path, http.FileServer(root))
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") && len(r.URL.Path) > 1 && r.URL.Path != (origPath+"/") {
			http.NotFound(w, r)
			return
		}
		fs.ServeHTTP(w, r)
	}))
}

// Healthz - handle health requests, will start returning 200, when server will be up
func Healthz(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && strings.HasSuffix(strings.ToLower(r.URL.Path), "/healthz") {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// Readyz - handle ready requests, will start returning 200, only when server will be ready to server the traffic
func Readyz(ready *atomic.Value) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" && strings.HasSuffix(strings.ToLower(r.URL.Path), "/readyz") {
				if ready == nil || !ready.Load().(bool) {
					w.Header().Set("Content-Type", "text/plain")
					w.WriteHeader(http.StatusServiceUnavailable)
					return
				}
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// NotFound - return a error page for not found
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte("Not found"))
}
