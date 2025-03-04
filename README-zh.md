# HostModifier
HostModifier 使用 Golang

## API 端點

### 顯示所有註冊的主機
- **端點**: `[GET] /api/hosts`
- **描述**: 顯示 `/etc/hosts` 文件中所有註冊的主機列表。
- **範例請求**:
  ```sh
  curl -X GET http://localhost:8080/api/hosts
  ```
- **範例回應**:
  ```json
  {
    "hosts": [
      "127.0.0.1 example.com",
      "10.10.10.100 beta.example.com"
    ]
  }
  ```

### 創建新主機
- **端點**: `[POST] /api/hosts`
- **描述**: 在 `/etc/hosts` 文件中創建一個新的主機條目。
- **範例請求**:
  ```sh
  curl -X POST http://localhost:8080/api/hosts -d '{"ip": "10.10.10.100", "hostname": "beta.example.com"}' -H "Content-Type: application/json"
  ```
- **範例回應**:
  ```json
  {
    "message": "Host entry created successfully."
  }
  ```

### 修改現有主機
- **端點**: `[PUT] /api/hosts/{host}`
- **描述**: 修改 `/etc/hosts` 文件中的現有主機條目。
- **範例請求**:
  ```sh
  curl -X PUT http://localhost:8080/api/hosts/beta.example.com -d '{"ip": "10.10.10.101"}' -H "Content-Type: application/json"
  ```
- **範例回應**:
  ```json
  {
    "message": "Host entry modified successfully."
  }
  ```

### 確認主機的 IP
- **端點**: `[GET] /api/hosts/{host}`
- **描述**: 確認 `/etc/hosts` 文件中特定主機的 IP 地址。
- **範例請求**:
  ```sh
  curl -X GET http://localhost:8080/api/hosts/beta.example.com
  ```
- **範例回應**:
  ```json
  {
    "ip": "10.10.10.100"
  }
  ```

### 刪除主機
- **端點**: `[DELETE] /api/hosts/{host}`
- **描述**: 從 `/etc/hosts` 文件中刪除主機條目。
- **範例請求**:
  ```sh
  curl -X DELETE http://localhost:8080/api/hosts/beta.example.com
  ```
- **範例回應**:
  ```json
  {
    "message": "Host entry deleted successfully."
  }
  ```

## 權限
- 該應用程序需要修改 `/etc/hosts` 文件的權限。確保應用程序具有讀寫此文件的必要權限。

## setupRoutes 函數
- **文件**: `main.go`
- **描述**: 設置路由並啟動服務器。
- **代碼**:
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

## HostEntry 結構體
- **文件**: `routes.go`
- **描述**: 表示從 hosts 文件中解析的條目。
- **代碼**:
  ```go
  type HostEntry struct {
      IP       string `json:"ip,omitempty"`
      Hostname string `json:"hostname,omitempty"`
      Status   string `json:"status"`        // "active", "deleted"
      Raw      string `json:"raw,omitempty"` // Raw line from host file (only for deleted)
  }
  ```

## parseHostEntry 函數
- **文件**: `routes.go`
- **描述**: 分析 hosts 文件中的單行並返回 HostEntry。
- **代碼**:
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

## getHostsHandler 函數
- **文件**: `routes.go`
- **描述**: 處理獲取所有主機的請求。
- **代碼**:
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

## createHostHandler 函數
- **文件**: `routes.go`
- **描述**: 處理創建新主機的請求。
- **代碼**:
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

## modifyHostHandler 函數
- **文件**: `routes.go`
- **描述**: 處理修改現有主機的請求。
- **代碼**:
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

## confirmHostIPHandler 函數
- **文件**: `routes.go`
- **描述**: 處理確認主機 IP 的請求。
- **代碼**:
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

## deleteHostHandler 函數
- **文件**: `routes.go`
- **描述**: 處理刪除主機的請求。
- **代碼**:
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

## checkPermissions 函數
- **文件**: `main.go`
- **描述**: 檢查應用程序是否有權修改 `/etc/hosts` 文件。
- **代碼**:
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

## getHostsFilePath 函數
- **文件**: `main.go`
- **描述**: 根據操作系統返回 hosts 文件的路徑。
- **代碼**:
  ```go
  func getHostsFilePath() string {
      if runtime.GOOS == "windows" {
          return "C:\\Windows\\System32\\drivers\\etc\\hosts"
      }
      return "/etc/hosts"
  }
  ```

## readHostsFile 函數
- **文件**: `hosts.go`
- **描述**: 讀取 hosts 文件的內容並返回行。
- **代碼**:
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

## writeHostsFile 函數
- **文件**: `hosts.go`
- **描述**: 將行寫入 hosts 文件。
- **代碼**:
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

## addHost 函數
- **文件**: `hosts.go`
- **描述**: 向 hosts 文件中添加新主機條目。
- **代碼**:
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

## modifyHost 函數
- **文件**: `hosts.go`
- **描述**: 修改 hosts 文件中的現有主機條目。
- **代碼**:
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

## confirmHostIP 函數
- **文件**: `hosts.go`
- **描述**: 確認 hosts 文件中特定主機的 IP 地址。
- **代碼**:
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

## deleteHost 函數
- **文件**: `hosts.go`
- **描述**: 從 hosts 文件中刪除主機條目。
- **代碼**:
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
- **文件**: `Makefile`
- **描述**: 包含項目的構建、運行和依賴管理目標。
- **代碼**:
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
- **文件**: `go.mod`
- **描述**: 包含項目的模塊路徑和依賴項。
- **代碼**:
  ```go
  module github.com/sam33339999/HostModifier

  go 1.16

  require github.com/gorilla/mux v1.8.1
  ```
