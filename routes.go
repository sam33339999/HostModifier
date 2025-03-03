package main

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
)

func setupRoutes() {
	router := mux.NewRouter()

	router.HandleFunc("/api/hosts", getHostsHandler).Methods("GET")
	router.HandleFunc("/api/hosts", createHostHandler).Methods("POST")
	router.HandleFunc("/api/hosts/{host}", modifyHostHandler).Methods("PUT")
	router.HandleFunc("/api/hosts/{host}", confirmHostIPHandler).Methods("GET")
	router.HandleFunc("/api/hosts/{host}", deleteHostHandler).Methods("DELETE")

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func getHostsHandler(w http.ResponseWriter, r *http.Request) {
	hosts, err := readHostsFile()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string][]string{"hosts": hosts}
	json.NewEncoder(w).Encode(response)
}

func createHostHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		IP       string `json:"ip"`
		Hostname string `json:"hostname"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := addHost(requestData.IP, requestData.Hostname); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Host entry created successfully."}
	json.NewEncoder(w).Encode(response)
}

func modifyHostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hostname := vars["host"]

	var requestData struct {
		IP string `json:"ip"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := modifyHost(hostname, requestData.IP); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Host entry modified successfully."}
	json.NewEncoder(w).Encode(response)
}

func confirmHostIPHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hostname := vars["host"]

	ip, err := confirmHostIP(hostname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"ip": ip}
	json.NewEncoder(w).Encode(response)
}

func deleteHostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hostname := vars["host"]

	if err := deleteHost(hostname); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Host entry deleted successfully."}
	json.NewEncoder(w).Encode(response)
}
