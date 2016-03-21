package main

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
	"net/http"
	"os"
)

var db *geoip2.Reader

func main() {
	var mmdb, port string
	flag.StringVar(&mmdb, "mmdb", "", "Path to MaxMind database file")
	flag.StringVar(&port, "port", "9000", "Port number for HTTP server")
	flag.Parse()

	if mmdb == "" {
		log.Fatal("Please specify the path to MaxMind database file")
	}

	_db, err := geoip2.Open(mmdb)
	if err != nil {
		log.Fatal(err)
	}
	db = _db
	defer db.Close()

	log.Println("Starting go-geoip2-server...")
	log.Printf("mmdb=%s", mmdb)
	log.Printf("port=%s", port)

	r := mux.NewRouter()
	r.HandleFunc("/{ip}", handler)
	http.Handle("/",
		handlers.CombinedLoggingHandler(
			os.Stdout,
			handlers.ProxyHeaders(
				handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(r))))
	http.ListenAndServe(":"+port, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ip := net.ParseIP(vars["ip"])
	if ip == nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	record, err := db.City(ip)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(record)
}
