package main

import (
	"context"
	"flag"
//	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	u "github.com/pprisn/test_ajax_form/utils"
)

var loging string

func DemoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "templates/demo.html")
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "static/index.html")
	}
}


// получаем измененные данные и сохраняем их в БД
func FillHandler(w http.ResponseWriter, r *http.Request) {
	resp := u.Message(true, "success")
        resp["firstname"] = "Mortadelo"
	resp["lastname"] = "Filemon"
	resp["address_street"] = "Rua del Percebe 13"
	resp["address_city"] = "Madrid"
	resp["address_zip"] = "28010"
	resp["[name='emails[0]"] = "superintendencia@cia.es"
	u.Respond(w, resp)
//	http.Redirect(w, r, "/", 301)
}



func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print(r.RemoteAddr, "\t", r.Method, "\t", r.URL)
		h.ServeHTTP(w, r)
	})
}

func main() {
	var dir string

	flag.StringVar(&dir, "dir", "./static/", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	router := mux.NewRouter()
	router.HandleFunc("/", IndexHandler)
//        router.HandleFunc("/fill", FillHandler)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	

	router.Use(Middleware)


	srv := &http.Server{
		Handler: router,
		Addr:    "localhost:3001",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}

	//log.Fatal(srv.ListenAndServe())

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
