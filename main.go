package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/MikeCodeSun/go-mux-auth/route"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("PORT")
	r := route.UserRoute()
	fmt.Println("go server on Port", port)
	log.Fatal(http.ListenAndServe(":"+ port, r))
}