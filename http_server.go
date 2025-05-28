package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Healthmessage struct {
	Status string `json:"status"`
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func healthcheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Healthmessage{Status: "OK"})
}

func error(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "Internal Server Error\n")
}

func hostname(w http.ResponseWriter, req *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown hostname: %s\n", err)
	} else {
		fmt.Fprintf(w, "%s\n", hostname)
	}

}

func gracefulShutdown(sig chan os.Signal) {
	fmt.Printf("Signal received %s", <-sig)
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	go gracefulShutdown(c)

	http.HandleFunc("/headers", headers)
	http.HandleFunc("/healthcheck", healthcheck)
	http.HandleFunc("/ready", healthcheck)
	http.HandleFunc("/error", error)
	http.HandleFunc("/hostname", hostname)

	httpPort, present := os.LookupEnv("HTTP_PORT")
	if !present {
		httpPort = "4191"
	}
	fmt.Printf("Listen and serve on port %s\n", httpPort)
	http.ListenAndServe(fmt.Sprintf(":%s", httpPort), nil)
}
