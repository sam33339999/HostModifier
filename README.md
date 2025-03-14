# HostModifier
HostModifier With Golang

## API Endpoints

### Display all registered hosts
- **Endpoint**: `[GET] /api/hosts`
- **Description**: Displays the list of all registered hosts in the `/etc/hosts` file.
- **Example Request**:
  ```sh
  curl -X GET http://localhost:8080/api/hosts
  ```
- **Example Response**:
  ```json
  {
    "hosts": [
      "127.0.0.1 example.com",
      "10.10.10.100 beta.example.com"
    ]
  }
  ```

### Create a new host
- **Endpoint**: `[POST] /api/hosts`
- **Description**: Creates a new host entry in the `/etc/hosts` file.
- **Example Request**:
  ```sh
  curl -X POST http://localhost:8080/api/hosts -d '{"ip": "10.10.10.100", "hostname": "beta.example.com"}' -H "Content-Type: application/json"
  ```
- **Example Response**:
  ```json
  {
    "message": "Host entry created successfully."
  }
  ```

### Modify an existing host
- **Endpoint**: `[PUT] /api/hosts/{host}`
- **Description**: Modifies an existing host entry in the `/etc/hosts` file.
- **Example Request**:
  ```sh
  curl -X PUT http://localhost:8080/api/hosts/beta.example.com -d '{"ip": "10.10.10.101"}' -H "Content-Type: application/json"
  ```
- **Example Response**:
  ```json
  {
    "message": "Host entry modified successfully."
  }
  ```

### Confirm the IP of a host
- **Endpoint**: `[GET] /api/hosts/{host}`
- **Description**: Confirms the IP address of a specific host in the `/etc/hosts` file.
- **Example Request**:
  ```sh
  curl -X GET http://localhost:8080/api/hosts/beta.example.com
  ```
- **Example Response**:
  ```json
  {
    "ip": "10.10.10.100"
  }
  ```

### Delete a host
- **Endpoint**: `[DELETE] /api/hosts/{host}`
- **Description**: Deletes a host entry from the `/etc/hosts` file.
- **Example Request**:
  ```sh
  curl -X DELETE http://localhost:8080/api/hosts/beta.example.com
  ```
- **Example Response**:
  ```json
  {
    "message": "Host entry deleted successfully."
  }
  ```

## Permissions
- The application requires permission to modify the `/etc/hosts` file. Ensure that the application has the necessary permissions to read and write to this file.

## setupRoutes function
- **File**: `main.go`
- **Description**: Sets up the routes and starts the server.
- **Code**:
  ```go
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
  ```

## HostEntry struct
- **File**: `routes.go`
- **Description**: Represents a parsed entry from the hosts file.
- **Code**:
  ```go
  type HostEntry struct {
      IP       string `json:"ip,omitempty"`
      Hostname string `json:"hostname,omitempty"`
      Status   string `json:"status"`        // "active", "deleted"
      Raw      string `json:"raw,omitempty"` // Raw line from host file (only for deleted)
  }
  ```

## parseHostEntry function
- **File**: `routes.go`
- **Description**: Analyzes a single line from the hosts file and returns a HostEntry.
- **Code**:
  ```go
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
  ```

## getHostsHandler function
- **File**: `routes.go`
- **Description**: Handles the request to get all hosts.
- **Code**:
  ```go
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
  ```

## createHostHandler function
- **File**: `routes.go`
- **Description**: Handles the request to create a new host.
- **Code**:
  ```go
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
  ```

## modifyHostHandler function
- **File**: `routes.go`
- **Description**: Handles the request to modify an existing host.
- **Code**:
  ```go
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
  ```

## confirmHostIPHandler function
- **File**: `routes.go`
- **Description**: Handles the request to confirm the IP of a host.
- **Code**:
  ```go
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
  ```

## deleteHostHandler function
- **File**: `routes.go`
- **Description**: Handles the request to delete a host.
- **Code**:
  ```go
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
  ```

## checkPermissions function
- **File**: `main.go`
- **Description**: Checks if the application has permission to modify the `/etc/hosts` file.
- **Code**:
  ```go
  func checkPermissions() bool {
      hostsFilePath := getHostsFilePath()
      file, err := os.OpenFile(hostsFilePath, os.O_WRONLY, 0644)
      if err != nil {
          return false
      }
      defer file.Close()
      return true
  }
  ```

## getHostsFilePath function
- **File**: `main.go`
- **Description**: Returns the path to the hosts file based on the operating system.
- **Code**:
  ```go
  func getHostsFilePath() string {
      if runtime.GOOS == "windows" {
          return "C:\\Windows\\System32\\drivers\\etc\\hosts"
      }
      return "/etc/hosts"
  }
  ```

## readHostsFile function
- **File**: `hosts.go`
- **Description**: Reads the contents of the hosts file and returns the lines.
- **Code**:
  ```go
  func readHostsFile() ([]string, error) {
      hostsFilePath := getHostsFilePath()
      file, err := os.Open(hostsFilePath)
      if err != nil {
          return nil, err
      }
      defer file.Close()

      var lines []string
      scanner := bufio.NewScanner(file)
      for scanner.Scan() {
          lines = append(lines, scanner.Text())
      }
      return lines, scanner.Err()
  }
  ```

## writeHostsFile function
- **File**: `hosts.go`
- **Description**: Writes the lines to the hosts file.
- **Code**:
  ```go
  func writeHostsFile(lines []string) error {
      hostsFilePath := getHostsFilePath()
      file, err := os.OpenFile(hostsFilePath, os.O_WRONLY|os.O_TRUNC, 0644)
      if err != nil {
          return err
      }
      defer file.Close()

      writer := bufio.NewWriter(file)
      for _, line := range lines {
          _, err := writer.WriteString(line + "\n")
          if err != nil {
              return err
          }
      }
      return writer.Flush()
  }
  ```

## addHost function
- **File**: `hosts.go`
- **Description**: Adds a new host entry to the hosts file.
- **Code**:
  ```go
  func addHost(ip, hostname string) error {
      lines, err := readHostsFile()
      if err != nil {
          return err
      }

      newEntry := fmt.Sprintf("%s %s   %s (%s)", ip, hostname, commentPrefix, time.Now().Format("2006/01/02 15:04:05"))
      lines = append(lines, newEntry)

      return writeHostsFile(lines)
  }
  ```

## modifyHost function
- **File**: `hosts.go`
- **Description**: Modifies an existing host entry in the hosts file.
- **Code**:
  ```go
  func modifyHost(hostname, newIP string) error {
      lines, err := readHostsFile()
      if err != nil {
          return err
      }

      for i, line := range lines {
          if strings.Contains(line, hostname) {
              lines[i] = fmt.Sprintf("%s %s   %s (%s)", newIP, hostname, commentPrefix, time.Now().Format("2006/01/02 15:04:05"))
              break
          }
      }

      return writeHostsFile(lines)
  }
  ```

## confirmHostIP function
- **File**: `hosts.go`
- **Description**: Confirms the IP address of a specific host in the hosts file.
- **Code**:
  ```go
  func confirmHostIP(hostname string) (string, error) {
      lines, err := readHostsFile()
      if err != nil {
          return "", err
      }

      for _, line := range lines {
          if strings.Contains(line, hostname) {
              parts := strings.Fields(line)
              if len(parts) > 0 {
                  return parts[0], nil
              }
          }
      }

      return "", fmt.Errorf("host not found")
  }
  ```

## deleteHost function
- **File**: `hosts.go`
- **Description**: Deletes a host entry from the hosts file.
- **Code**:
  ```go
  func deleteHost(hostname string) error {
      lines, err := readHostsFile()
      if err != nil {
          return err
      }

      for i, line := range lines {
          if strings.Contains(line, hostname) {
              lines[i] = fmt.Sprintf("# DELETE AT(%s) %s", time.Now().Format("2006/01/02 15:04:05"), line)
              break
          }
      }

      return writeHostsFile(lines)
  }
  ```

## Makefile
- **File**: `Makefile`
- **Description**: Contains the build, run, and dependency management targets for the project.
- **Code**:
  ```makefile
  # Makefile for HostModifier

  # Define the Go module path
  MODULE_PATH := github.com/sam33339999/HostModifier

  # Define the binary name
  BINARY_NAME := hostmodifier

  # Define the build target
  .PHONY: build
  build:
      go build -o $(BINARY_NAME) .

  # Define the run target
  .PHONY: run
  run: build
      ./$(BINARY_NAME)

  # Define the deps target
  .PHONY: deps
  deps:
      go mod tidy
  ```

## go.mod
- **File**: `go.mod`
- **Description**: Contains the module path and dependencies for the project.
- **Code**:
  ```go
  module github.com/sam33339999/HostModifier

  go 1.16

  require github.com/gorilla/mux v1.8.1
  ```
