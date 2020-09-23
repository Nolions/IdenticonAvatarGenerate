package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	identicon "identiconAvatar"
	"log"
)

func main() {
	var (
		port = flag.String("port", "4321", "web service point")
	)

	r := gin.Default()
	r.GET("/:name", func(c *gin.Context) {
		name := c.Param("name")

		log.Printf("Identicon: %s", name)

		w := c.Writer
		w.Header().Set("Content-Type", "image/png")
		if err := identicon.Generate([]byte(name)).WriteImage(w); err != nil {
			log.Fatalf("Output Identicon Image to  web server Fail, err: %v", err)
		}
	})

	_ = r.Run(fmt.Sprintf("0.0.0.0: %s", *port))
}
