package main

import (
	"context"
	"fmt"
	"k8s.io/client-go/dynamic"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"github.com/gorilla/mux"
	"gopkg.in/natefinch/lumberjack.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	oofsGVR = schema.GroupVersionResource{
		Group:    "build.dev",
		Version:  "v1alpha1",
		Resource: "buildruns",
	}
)

func handler(w http.ResponseWriter, r *http.Request) {
	config, _ := rest.InClusterConfig()
	dynClient, errClient := dynamic.NewForConfig(config)

	w.Write([]byte(fmt.Sprintf("dynClient is, %s\n", dynClient)))

	if errClient != nil {
		w.Write([]byte(fmt.Sprintf("errClient is, %s\n", errClient)))
	}
	w.Write([]byte(fmt.Sprintf("Hello, %s\n", config)))

	crdClient := dynClient.Resource(oofsGVR)
	w.Write([]byte(fmt.Sprintf("crdClient, %s\n", crdClient)))
	crd, errCrd := crdClient.Namespace("a8058b2f-8a6d").Get("kaniko-golang-buildrun-liya-03", metav1.GetOptions{})
	if errCrd != nil {
		w.Write([]byte(fmt.Sprintf("errCrd, %s\n", errCrd)))
		//http.Error(w, errCrd.Error(), http.StatusBadRequest)
	}
	w.Write([]byte(fmt.Sprintf("crdClient, %s\n", crd)))

	query := r.URL.Query()
	name := query.Get("name")
	if name == "" {
		name = "Guest"
	}
	w.Write([]byte(fmt.Sprintf("Hello, %s\n", name)))
}

func main() {
	// Create Server and Route Handlers
	r := mux.NewRouter()

	r.HandleFunc("/", handler)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Configure Logging
	LOG_FILE_LOCATION := os.Getenv("LOG_FILE_LOCATION")
	if LOG_FILE_LOCATION != "" {
		log.SetOutput(&lumberjack.Logger{
			Filename:   LOG_FILE_LOCATION,
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		})
	}

	// Start Server
	go func() {
		log.Println("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful Shutdown
	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}
