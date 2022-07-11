package executor

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Run runs multiple services in parallel.
// It will keep running forever unless:-
// - The parent context gets cancelled
// - A SIGTERM or OS Interrupt is received
// - Any one of the services encounters an error
func Run(ctx context.Context, svc ...Service) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(svc))

	for i := range svc {
		go func(s Service) {
			defer wg.Done()
			if err := s(ctx); err != nil { // Wait
				log.Println(fmt.Errorf("svc exec encountered err: %w", err))
				cancel()
			}
		}(svc[i])
	}

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGTERM, os.Interrupt)

	select {
	case s := <-termChan:
		log.Println(fmt.Sprintf("Received signal: [%s]. Shutting down all services...", s.String()))
		cancel()
	case <-ctx.Done():
		log.Println("Context cancelled. Shutting down all services...")
	}

	wg.Wait()

	log.Println("All services shutdown successfully")
}
