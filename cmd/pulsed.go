package main

import (
	"fmt"
	"github.com/cobratbq/pulse"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var p *pulse.Connection

func main() {

	// initialize connection to mongodb
	var connectString = "mongodb://localhost/pulse"
	conn, err := pulse.Dial(connectString, "pulse", "pulses")
	if err != nil {
		log.Printf("Failed to connect to pulse database '%s': %s.", connectString, err.Error())
		return
	}
	defer conn.Close()
	p = conn

	// initialize http pulse server
	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/pulse", pulseInfo)
	r.HandleFunc("/show/{namespace}", pulseShow)
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

func pulseShow(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	vars := mux.Vars(req)
	ns, ok := vars["namespace"]
	if !ok {
		log.Println("Failed to find namespace path variable. Not recording this pulse.")
		return
	}
	pulses, err := p.Get(ns)
	if err != nil {
		log.Println(err.Error())
		return
	}
	fmt.Fprintln(resp, "Pulses:")
	for _, pulse := range pulses {
		fmt.Fprintf(resp, "%s: %s\n", pulse.Namespace, pulse.Time.String())
	}
}

func pulseRecord(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	vars := mux.Vars(req)
	ns, ok := vars["namespace"]
	if !ok {
		log.Println("Failed to find namespace path variable. Not recording this pulse.")
		return
	}
	log.Printf("PULSE '%s'", ns)
	p.Record(ns)
}
