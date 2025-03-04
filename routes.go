package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// HostEntry represents a parsed entry from the hosts file.
type HostEntry struct {
	IP       string `json:"ip,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	Status   string `json:"status"`        // "active", "deleted"
	Raw      string `json:"raw,omitempty"` // Raw line from host file (only for deleted)
}

func setupRoutes() {
	router := mux.NewRouter()

	router.HandleFunc("/api/hosts", getHostsHandler).Methods("GET")
	router.HandleFunc("/api/hosts", createHostHandler).Methods("POST")
	router.HandleFunc("/api/hosts/{host}", modifyHostHandler).Methods("PUT")
	router.HandleFunc("/api/hosts/{host}", confirmHostIPHandler).Methods("GET")
	router.HandleFunc("/api/hosts/{host}", deleteHostHandler).Methods("DELETE")

	http.Handle("/", router)

	port := ":8080"
	fmt.Printf("Starting Host Modifier API server on port %s\n", port)
	fmt.Println("Available routes:")
	fmt.Println("  GET /api/hosts - Get all hosts")
	fmt.Println("  POST /api/hosts - Create a new host entry (JSON body: {ip: \"127.0.0.1\", hostname: \"example.com\"})")
	fmt.Println("  PUT /api/hosts/{host} - Modify a host entry (JSON body: {ip: \"192.168.1.1\"})")
	fmt.Println("  GET /api/hosts/{host} - Confirm the IP of a host")
	fmt.Println("  DELETE /api/hosts/{host} - Delete a host entry")
	fmt.Println("Server is ready to accept requests.")

	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

func getHostsHandler(w http.ResponseWriter, r *http.Request) {
	lines, err := readHostsFile()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var hostEntries []HostEntry
	for _, line := range lines {
		entry := parseHostEntry(line)
		//skip the comment line
		if entry.Status != "" {
			hostEntries = append(hostEntries, entry)
		}
	}

	response := map[string][]HostEntry{"hosts": hostEntries}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func createHostHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		IP       string `json:"ip"`
		Hostname string `json:"hostname"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}
	if requestData.IP == "" || requestData.Hostname == "" {
		http.Error(w, "Missing 'ip' or 'hostname' in request body.", http.StatusBadRequest)
		return
	}
	if err := addHost(requestData.IP, requestData.Hostname); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Host entry created successfully."}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func modifyHostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hostname := vars["host"]

	var requestData struct {
		IP string `json:"ip"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}
	if requestData.IP == "" {
		http.Error(w, "Missing 'ip' in request body.", http.StatusBadRequest)
		return
	}
	if err := modifyHost(hostname, requestData.IP); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Host entry modified successfully."}
	w.Header().Set("Content-Type", "application/json")
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
	w.Header().Set("Content-Type", "application/json")
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// parseHostEntry analyzes a single line from the hosts file and returns a HostEntry.
func parseHostEntry(line string) HostEntry {
	line = strings.TrimSpace(line)

	if strings.HasPrefix(line, "# DELETE AT") {
		return HostEntry{
			Status: "deleted",
			Raw:    line,
		}
	}

	if strings.HasPrefix(line, "#") {
		//skip the comment line
		return HostEntry{}
	}
	parts := strings.Fields(line)
	if len(parts) >= 2 {
		return HostEntry{
			IP:       parts[0],
			Hostname: parts[1],
			Status:   "active",
		}
	}

	// Unrecognizable format, treat as a comment, and skip this line
	return HostEntry{}
}
