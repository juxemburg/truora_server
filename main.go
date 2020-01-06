package main

import (
	"context"
	"fmt"
	"github.com/juxemburg/truora_server/controllers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

const (
	portArgPosition = iota + 1
)

var (
	srv http.Server
)

func main() {
	port, err := strconv.ParseInt(os.Args[portArgPosition], 0, 32)
	if err != nil {
		log.Println("Error while parsing input parameters")
		log.Print(err.Error())
		os.Exit(1)
	}
	startup(int(port))
}

func startup(port int) {
	log.Print("Starting server on port:", port, " ...")
	srv.Addr = fmt.Sprint(":", port)
	srv.Handler = controllers.GetRouteConfig()

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)

		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
	log.Println("Server shutdown...")
}
