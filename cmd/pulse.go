package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/pulse", pulseInfo)
	r.HandleFunc("/pulse/{namespace}", pulseRecord)
	http.Handle("/", r)
	log.Println(http.ListenAndServe(":8080", nil))
}

func home(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	resp.Write([]byte("Hello world!"))
}

func pulseInfo(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	resp.Write([]byte("Pulse to your own namespace, e.g. '/pulse/my-namespace'"))
}

func pulseRecord(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	vars := mux.Vars(req)
	ns, ok := vars["namespace"]
	if !ok {
		log.Println("Failed to find namespace path variable. Not recording this pulse.")
	}
	log.Printf("PULSE '%s'", ns)
}
