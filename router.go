package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func installKillSignalHandler() chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	return sigChan
}

func (api *api) handler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Kalhspera\n")
}

func (api *api) Score(w http.ResponseWriter, r *http.Request) {
	//io.WriteString(w, "To score tou api einai %d \n", api.score)
	fmt.Fprintf(w, "To score tou api einai %d \n", api.score)
}

func (api *api) Router() /*mux.Router*/ {
	//r := mux.NewRouter()
	api.router = mux.NewRouter()
	api.router.HandleFunc("/", api.handler)
	api.router.HandleFunc("/score", api.Score)
	//return r
}

type api struct {
	score  int
	router *mux.Router
}

type Server struct {
	server   *http.Server
	listener net.Listener
}

func (srv Server) serve(listener net.Listener, done chan struct{}) {
	err := srv.server.Serve(listener)
	if err != nil {
		fmt.Println("err: ", err)
	}
	fmt.Println("Returning from server")
	close(done)

}

func main() {
	myApi := &api{
		score: 45,
	}
	myApi.Router()
	done := make(chan struct{})
	//myApi.router = myApi.Router()
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Println(err)
	}
	srv := &http.Server{
		Handler: myApi.router,
		//Addr:    "127.0.0.1:8080",
	}
	sigChan := installKillSignalHandler()

	myServer := &Server{
		server:   srv,
		listener: listener,
	}
	//go srv.ListenAndServe()
	//go srv.Serve(listener)
	anotherChan := make(chan struct{})
	go myServer.serve(listener, done)
	fmt.Println("Server is running...")

	timh := <-sigChan
	myServer.server.Shutdown(context.Background())
	fmt.Println("Server is closing... ", timh)
	go func() {
		<-done
		fmt.Println("Epiasa kai egw to done channel")
		close(anotherChan)
	}()
	<-done
	<-anotherChan
	fmt.Println("Server is safely closed!!")

}
