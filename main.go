package main

import (
	"fmt"
	"log"

	"github.com/Cantabar/friend-quest-ark-controller/core"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	instances := core.GetInstancesAndActivePlayers()
	fmt.Println(instances)
}
