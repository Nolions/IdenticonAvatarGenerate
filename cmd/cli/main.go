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
		name = flag.String("name", "", "Identicon Name")
		path = flag.String("path", "/", "Generate Image path on Disk")
	)
	flag.Parse()

	if *name == "" {
		flag.Usage()
		os.Exit(0)
	}

	data := []byte(*name) // utf8 to byte(10)

	i := identicon.Generate(data)

	f, err := os.Create(fmt.Sprintf("%s/%s.png", *path, *name))
	if err != nil {
		fmt.Printf("error:Output Identicon Image Fail, err:  %v", err)
		return
	}
	defer f.Close()

	if err := i.WriteImage(f); err != nil {
		log.Fatalln(err)
	}
}
