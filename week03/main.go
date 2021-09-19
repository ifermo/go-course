package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	group, ctx := errgroup.WithContext(context.Background())
	group.Go(func() error {
		return signalHandle(ctx)
	})
	group.Go(func() error {
		return httpHandle(ctx, ":80")
	})
	group.Go(func() error {
		return httpHandle(ctx, ":8080")
	})
	if err := group.Wait(); err != nil {
		log.Fatalf("err: %v", err)
	}
}

func signalHandle(ctx context.Context) error {
	quitCh := make(chan os.Signal)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case sig := <-quitCh:
		return fmt.Errorf("signal exit: %s", sig)
	case <-ctx.Done():
		return nil
	}
}

func httpHandle(ctx context.Context, addr string) error {
	mux := http.NewServeMux()
	done := make(chan int)
	mux.HandleFunc("/done", func(w http.ResponseWriter, req *http.Request) {
		_, _ = fmt.Fprintf(w, "exit!")
		done <- 1
	})
	server := &http.Server{Addr: addr, Handler: mux}
	go func() {
		select {
		case <-ctx.Done():
			_ = server.Shutdown(context.Background())
		case <-done:
			log.Printf("server[%s] done", addr)
			_ = server.Shutdown(context.Background())
		}
	}()
	return server.ListenAndServe()
}
