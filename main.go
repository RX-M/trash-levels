package main

import (
	mux "github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strconv"
)

type CanLevel struct {
	Id    string
	Level int
}

var levels = make(map[string]int)

var callCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Namespace: "hello",
		Name:      "call_counter",
		Help:      "Number of calls made to all routes (including /healthz but not /metrics)",
	})

func levelHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Println("Responding to request for level: " + id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"id": "` + id + `","level": ` + strconv.Itoa(levels[id]) + `}`))
	callCounter.Add(1)
}

func livenessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "alive!"}`))
	callCounter.Add(1)
}

func seedLevels() {
	levels["10"] = 80
	levels["15"] = 81
	levels["20"] = 82
	levels["25"] = 83
	levels["30"] = 84
}

func main() {
	seedLevels()
	prometheus.MustRegister(callCounter)

	r := mux.NewRouter()
	r.HandleFunc("/cans/{id}", levelHandler).Methods("GET")
	r.HandleFunc("/healthz", livenessHandler).Methods("GET")
	r.Handle("/metrics", promhttp.Handler())

	port := "8080"
	log.Println("Listening on: " + port)
	go log.Fatal(http.ListenAndServe(":"+port, r))
	select {}
}
