package main

import (
	"gecko/cmd/server/cmd"
	"log"
	"net/http"
)
import _ "net/http/pprof"

func main() {
	go func() {
		err := http.ListenAndServe("0.0.0.0:7777", nil)
		if err != nil {
			log.Fatalf("Debug error: %v", err)
		}
	}()
	cmd.Execute()
}
