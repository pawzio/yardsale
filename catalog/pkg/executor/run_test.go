package executor

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestRun_CtxCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		Run(ctx,
			func(ctx context.Context) error {
				fmt.Println("starting 1")
				go func() {
					fmt.Println("started 1")
					time.Sleep(time.Hour)
				}()
				<-ctx.Done()
				fmt.Println("completed 1")
				return nil
			},
			func(ctx context.Context) error {
				fmt.Println("starting 2")
				go func() {
					fmt.Println("started 2")
					time.Sleep(time.Hour)
				}()
				<-ctx.Done()
				fmt.Println("completed 2")
				return nil
			},
			func(ctx context.Context) error {
				fmt.Println("starting 3")
				go func() {
					fmt.Println("started 3")
					time.Sleep(time.Hour)
				}()
				<-ctx.Done()
				fmt.Println("completed 3")
				return nil
			},
		)
	}()

	time.Sleep(5 * time.Second)

	cancel()

	wg.Wait()
}

func TestRun_ServiceErr(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		Run(ctx,
			func(ctx context.Context) error {
				fmt.Println("starting 1")
				go func() {
					fmt.Println("started 1")
					time.Sleep(time.Hour)
				}()
				<-ctx.Done()
				fmt.Println("completed 1")
				return nil
			},
			func(ctx context.Context) error {
				fmt.Println("starting 2")
				go func() {
					fmt.Println("started 2")
					time.Sleep(time.Hour)
				}()
				<-ctx.Done()
				fmt.Println("completed 2")
				return nil
			},
			func(ctx context.Context) error {
				fmt.Println("starting 3")

				ch := make(chan error, 1)
				go func() {
					fmt.Println("started 3")
					time.Sleep(5 * time.Second)
					ch <- errors.New("some error")
				}()
				select {
				case err := <-ch:
					return err
				case <-ctx.Done():
				}
				fmt.Println("completed 3")
				return nil
			},
		)
	}()

	wg.Wait()
}
