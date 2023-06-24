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

3. Starting server
    #### for service provider
    cd project_directory/service-provider
    go mod tidy
    go run main.go

    #### for user
    cd project_directory/user
    go mod tidy
    go run main.go
