package cmd

import (
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile)
	e := NewRouter()
	e.Logger.Fatal(e.Start("localhost:5050"))
}
