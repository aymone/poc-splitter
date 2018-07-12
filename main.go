package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	router.Static("/", "./web")

	router.POST("/upload", func(c *gin.Context) {
		// Source
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}

		dst := fmt.Sprintf("./data/dir/%s", file.Filename)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		c.String(http.StatusOK, "File %s uploaded successfully", file.Filename)
	})

	router.Run(":8080")

	// filename := "portabilidade.pdf"
	// err := chunk.Join("./data/dir", filename)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
}
