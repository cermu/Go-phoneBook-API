package main

import (
	"context"
	"github.com/cermu/Go-phoneBook-API/models"
	"github.com/cermu/Go-phoneBook-API/routers"
	utl "github.com/cermu/Go-phoneBook-API/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// get file name and line number when the code crashes
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	models.InitDB()    // Initialize a database connection
	models.MigrateDB() //perform database migrations
	// close database and redis connection after use
	defer func() {
		if dbErr := models.DBConnection.Close(); dbErr != nil {
			log.Printf("WARNING | Database connection failed to close with message: %v\n", dbErr.Error())
		}

		if redisErr := utl.RedisClient().Close(); redisErr != nil {
			log.Printf("WARNING | Redis connection failed to close with message: %v\n", redisErr.Error())
		}
	}()

	apiPort := utl.ReadConfigs().GetInt("APP.PORT")
	apiENV := utl.ReadConfigs().GetString("APP.ENV")

	// API server
	apiServer := &http.Server{
		Addr:    utl.ReadConfigs().GetString("APP.ADDRESS"),
		Handler: routers.NewRouter(),
	}

	// start the server in a go routine
	go func() {
		log.Printf("INFO | Starting API server on port: %v with %v configs\n", apiPort, apiENV)
		err := apiServer.ListenAndServe()

		if err != nil {
			if err.Error() != "http: Server closed" {
				log.Fatalf("ERROR | Failed to start API server: %v\n", err)
			}
		}
	}()

	// shut down the server
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// block until a signal is received
	receivedSignal := <-ch

	log.Printf("WARNING | Shutting down API server %v signal received\n", receivedSignal)
	err := apiServer.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("ERROR | Failed to shut down API server: %v\n", err)
	}
	log.Println("INFO | API server has shut down")
}
