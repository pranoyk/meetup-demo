# Demo project Golang meetup July 2023

### About
It is demo project for mircoservice based architecture in Go

### Installation


1. Clone the repo
   ```sh
   git clone https://https://github.com/avinilcodes/meetup-demo
   ```
2. Install MongoDB compass from their website(UI)
    ##### start server using default port

3. Install go from https://go.dev/doc/install

4. Starting server
    #### for service provider
    ```sh
    cd meetup-demo/service-provider
    go mod tidy
    go run main.go
    ```

    #### for user
    ```sh
    cd meetup-demo/user
    go mod tidy
    go run main.go
    ```
    
5. Use postman collection for APIs

#### gRPC

Install gRPC plugins

```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

export PATH="$PATH:$(go env GOPATH)/bin"
```

Proto files are located under common folder

Use the following command in order to generate proto files

```sh
cd common

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/user.proto
```
