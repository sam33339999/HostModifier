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
