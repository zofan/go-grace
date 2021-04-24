package grace

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Grace struct {
	signal chan os.Signal
}

func New() *Grace {
	g := &Grace{
		signal: make(chan os.Signal, 1),
	}

	signal.Notify(g.signal, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM)

	return g
}

func (g *Grace) WaitShutdown() {
	log.Println(`Wait shutdown...`)
	<-g.signal
	log.Println(`Shutdown...`)
}

func (g *Grace) ForceShutdown() {
	log.Println(`Force shutdown`)
	g.signal <- syscall.SIGTERM
}

func (g *Grace) Shutdown(timeout time.Duration, fn func(ctx context.Context) error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := fn(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
