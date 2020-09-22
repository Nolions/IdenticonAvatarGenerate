package main

import (
	"flag"
	"fmt"
	"identiconAvatar"
	"log"
	"os"
)

func main() {
	var (
		name = flag.String("name", "", "Set the name where you want to generate an Identicon for")
	)
	flag.Parse()

	if *name == "" {
		flag.Usage()
		os.Exit(0)
	}

	data := []byte(*name) // utf8 to byte(10)

	fmt.Println(*name)
	fmt.Println(data[:])
	i := identicon.Generate(data)

	if err := i.WriteImage(); err != nil {
		log.Fatalln(err)
	}
}
