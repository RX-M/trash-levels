package main

import (
	mux "github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type CanLevel struct {
	Id    string
	Level int
}

var levels = make(map[string]int)

func levelHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Println("Responding to request for level: " + id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"id": "` + id + `","level": ` + strconv.Itoa(levels[id]) + `}`))
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

	r := mux.NewRouter()
	r.HandleFunc("/cans/{id}", levelHandler).Methods("GET")

	port := "8080"
	log.Println("Listening on: " + port)
	go log.Fatal(http.ListenAndServe(":"+port, r))
	select {}
}
