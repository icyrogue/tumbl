package main

import (
	"log"

	tumbl "github.com/icyrogue"
)

func main() {
	opts := tumbl.Options{
		Dst: "./repo",
	}
	puller := tumbl.NewPuller("github.com/miguelgfierro/scripts", &opts)
	executor := tumbl.NewExecutor(&opts)

	api := tumbl.NewAPI(executor, puller)
	log.Fatal(api.Start())

}
