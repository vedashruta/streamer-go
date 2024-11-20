package main

import (
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		resp := map[string]any{
			"ok":    false,
			"error": "Access Denied",
		}
		c.JSON(200, resp)
	})

	router.GET("/list", func(c *gin.Context) {
		wd, err := os.Getwd()
		if err != nil {
			c.String(http.StatusNotFound, "File not found")
			return
		}
		path := wd + "/media/"
		dirs, err := os.ReadDir(path)
		if err != nil {
			c.String(http.StatusNotFound, "File not found")
			return
		}
		type fileData struct {
			Name string  `json:"name"`
			Size float64 `json:"size"` // Size in bytes
		}
		fileDetails := []fileData{}
		for _, file := range dirs {
			if !file.IsDir() {
				filePath := filepath.Join(path, file.Name())
				fileInfo, err := os.Stat(filePath)
				if err != nil {
					return

				}
				fileDetails = append(fileDetails, fileData{
					Name: file.Name(),
					Size: float64(fileInfo.Size()) / (1024 * 1024),
				})
			}
		}
		c.JSON(200, fileDetails)
	})

	router.GET("/stream/:file", func(c *gin.Context) {
		wd, err := os.Getwd()
		if err != nil {
			c.String(http.StatusNotFound, "File not found")
			return
		}
		path := wd + "/media/"
		fileName := c.Param("file")
		file, err := os.Open(path + fileName)
		if err != nil {
			c.String(http.StatusNotFound, "File not found")
			return
		}
		defer file.Close()
		ext := filepath.Ext(fileName)
		contentType := mime.TypeByExtension(ext)
		c.Header("Content-Type", contentType)
		buffer := make([]byte, 4096*4096)
		io.CopyBuffer(c.Writer, file, buffer)
	})
	router.Run(":9944")
}
