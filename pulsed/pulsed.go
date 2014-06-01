package main

import (
	"fmt"
	"github.com/cobratbq/flagtag"
	"github.com/cobratbq/pulse"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type configuration struct {
	DatabaseHost string `flag:"dbhost,localhost,Database host."`
	DatabaseName string `flag:"dbname,pulse,Database name."`
	Collection   string `flag:"collection,pulses,Collection for recording pulses."`
	Port         uint   `flag:"port,8000,Port number for pulse server."`
}

func main() {
	// initialize with specified program arguments
	var config configuration
	flagtag.MustConfigureAndParse(&config)

	// initialize connection to mongodb
	var connectString = fmt.Sprintf("mongodb://%s", config.DatabaseHost)
	conn, err := pulse.Dial(connectString, config.DatabaseName, config.Collection)
	if err != nil {
		log.Printf("Failed to connect to pulse database '%s': %s.", connectString, err.Error())
		return
	}
	defer conn.Close()
	log.Println("Connection to mongoDB established.")

	// initialize http pulse server
	r := mux.NewRouter()
	r.HandleFunc("/", info)
	r.HandleFunc("/show/{namespace}", createPulseShowHandler(conn))
	r.HandleFunc("/pulse/{namespace}", createPulseRecordHandler(conn))
	http.Handle("/", r)

	// start http pulse server
	log.Printf("Starting http server on :%d ...", config.Port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}

func info(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	resp.Write([]byte("Usage:\n"))
	resp.Write([]byte("* Record a pulse to your own namespace, e.g. '/pulse/test' for namespace 'test'.\n"))
	resp.Write([]byte("* Show previously recorded pules for your own namespace, e.g. '/show/test' for namespace 'test'.\n"))
}

func createPulseShowHandler(conn *pulse.Connection) func(http.ResponseWriter, *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		vars := mux.Vars(req)
		ns, ok := vars["namespace"]
		if !ok {
			log.Println("Failed to find namespace path variable.")
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		pulses, err := conn.Get(ns)
		if err != nil {
			log.Println(err.Error())
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, pulse := range pulses {
			fmt.Fprintf(resp, "%s,%s\n", pulse.Namespace, pulse.Time.String())
		}
	}
}

func createPulseRecordHandler(conn *pulse.Connection) func(http.ResponseWriter, *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		vars := mux.Vars(req)
		ns, ok := vars["namespace"]
		if !ok {
			log.Println("Failed to find namespace path variable. Not recording this pulse.")
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("PULSE '%s'", ns)
		conn.Record(ns)
		resp.Write([]byte("OK"))
	}
}
