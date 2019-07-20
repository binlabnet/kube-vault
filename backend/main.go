package main

import (
	"context"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

type args struct {
	Namespace  string `long:"n" env:"NAMESPACE" default:"default" description:"default namespace"`
	Kubeconfig string `long:"kubeconfig" env:"HOME" default:"" description:"(optional) absolute path to the kubeconfig file"`
	Version    string `long:"version" env:"VERSION" default:"unknown" description:"version number"`
	Port       int    `long:"port" env:"PORT" default:"9090" description:"rest server port"`
	Debug      bool   `long:"debug" env:"DEBUG" description:"debug"`
}

type app struct {
	args  args
	srv   *Rest
	vlt   *Vault
	ready *atomic.Value
}

func main() {
	var args args
	p := flags.NewParser(&args, flags.Default)
	if _, err := p.Parse(); err != nil {
		log.Panicf("error when parsing arguments: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Print("interrupt signal")
		cancel()
	}()

	app, err := New(args)
	if err != nil {
		log.Printf("error when setup application: %+v", err)
	}

	err = app.run(ctx, cancel)
	if err != nil {
		log.Printf("error when start application: %+v", err)
	}

	log.Print("application terminated")
}

func New(args args) (*app, error) {
	ready := &atomic.Value{}
	ready.Store(false)

	vlt := &Vault{}
	err := vlt.Init(args.Kubeconfig, args.Debug)
	if err != nil {
		return nil, err
	}

	srv := &Rest{
		Version: args.Version,
		Port:    args.Port,
		Ready:   ready,
		Vault:   vlt,
	}

	return &app{
		args:  args,
		srv:   srv,
		vlt:   vlt,
		ready: ready,
	}, nil
}

func (a *app) run(ctx context.Context, cancel context.CancelFunc) error {
	go func() {
		<-ctx.Done()
		a.ready.Store(false)

		err := a.srv.Shutdown(ctx)
		if err != nil {
			log.Printf("rest shutdown error, %+v", err)
		}
	}()

	return a.srv.Run(ctx)
}
