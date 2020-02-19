package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	mux "github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	levels      = make(map[string]int)
	mem         = []int64{}
	ready       = false
	callCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "hello",
			Name:      "call_counter",
			Help:      "Number of calls made to all routes (including /healthz but not /metrics)",
		})
	memGague = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "hello",
			Name:      "mem_gague",
			Help:      "Amount of application memory currently allocated",
		})
)

func levelHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Println("Responding to request for level: " + id)
	//Taper mem alloc benefit after 10 allocations
	mSize := len(mem)
	if mSize > (5 * 1024 * 1024) {
		mSize = (5 * 1024 * 1024) + (mSize-(5*1024*1024))/5
	}
	time.Sleep(time.Duration((11*1024*1024)/(mSize+(1024*1024))) * time.Second)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"id": "` + id + `","level": ` + strconv.Itoa(levels[id]) + `}`))
	callCounter.Add(1)
}

func memHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Responding to request for mem")
	buf := make([]int64, 1024*1024)
	for index := range buf {
		buf[index] = 0
	}
	mem = append(mem, buf...)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "mem allocated!"}`))
	memGague.Add(float64(8 * len(buf)))
	callCounter.Add(1)
}

func livenessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "alive!"}`))
	callCounter.Add(1)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if ready {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "ready!"}`))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"message": "not ready!"}`))
	}
	callCounter.Add(1)
}

func seedLevels() {
	levels["10"] = 80
	levels["15"] = 81
	levels["20"] = 82
	levels["25"] = 83
	levels["30"] = 84
}

func startUp() {
	time.Sleep(30 * time.Second)
	ready = true
}

func main() {
	seedLevels()
	prometheus.MustRegister(callCounter)
	prometheus.MustRegister(memGague)

	r := mux.NewRouter()
	r.HandleFunc("/cans/{id}", levelHandler).Methods("GET")
	r.HandleFunc("/mem", memHandler).Methods("GET")
	r.HandleFunc("/healthz", livenessHandler).Methods("GET")
	r.HandleFunc("/readiz", readinessHandler).Methods("GET")
	r.Handle("/metrics", promhttp.Handler()).Methods("GET")

	port := "8080"
	log.Println("Trash Levels Server listening on port: " + port)
	go startUp()
	go log.Fatal(http.ListenAndServe(":"+port, r))
	select {}
}
