package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"server/config"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal(err)
		return
	}
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
		// Get the current working directory
		wd, err := os.Getwd()
		if err != nil {
			c.String(http.StatusNotFound, "File not found")
			return
		}
		// Set the path for media files
		path := wd + "/media/"
		fileName := c.Param("file")

		// Open the file
		file, err := os.Open(path + fileName)
		if err != nil {
			c.String(http.StatusNotFound, "File not found")
			return
		}
		defer file.Close()

		// Get file extension for content type
		ext := filepath.Ext(fileName)
		contentType := mime.TypeByExtension(ext)
		c.Header("Content-Type", contentType)

		// Get file info to determine the size
		fileInfo, err := file.Stat()
		if err != nil {
			c.String(http.StatusInternalServerError, "Unable to get file info")
			return
		}
		fileSize := fileInfo.Size()

		// Handle Range requests
		rangeHeader := c.GetHeader("Range")
		if rangeHeader != "" {
			rangeParts := strings.Split(rangeHeader, "=")
			if len(rangeParts) == 2 && strings.HasPrefix(rangeParts[1], "bytes") {
				rangeBytes := strings.TrimPrefix(rangeParts[1], "bytes=")
				rangeValues := strings.Split(rangeBytes, "-")
				if len(rangeValues) == 2 {
					start, err := strconv.ParseInt(rangeValues[0], 10, 64)
					if err != nil {
						c.String(http.StatusBadRequest, "Invalid range start")
						return
					}
					end := start
					if rangeValues[1] != "" {
						end, err = strconv.ParseInt(rangeValues[1], 10, 64)
						if err != nil {
							c.String(http.StatusBadRequest, "Invalid range end")
							return
						}
					} else {
						end = fileSize - 1
					}
					c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
					c.Header("Content-Length", fmt.Sprintf("%d", end-start+1))
					c.Status(http.StatusPartialContent)
					_, err = file.Seek(start, io.SeekStart)
					if err != nil {
						c.String(http.StatusInternalServerError, "Failed to seek file")
						return
					}
					buffer := make([]byte, config.Conf.BufferSize)
					bytesRemaining := end - start + 1
					for bytesRemaining > 0 {
						chunkSize := int64(len(buffer))
						if bytesRemaining < chunkSize {
							chunkSize = bytesRemaining
						}
						n, err := file.Read(buffer[:chunkSize])
						if err != nil && err != io.EOF {
							c.String(http.StatusInternalServerError, "Error reading file")
							return
						}
						_, err = c.Writer.Write(buffer[:n])
						if err != nil {
							c.String(http.StatusInternalServerError, "Error streaming file")
							return
						}

						bytesRemaining -= int64(n)
					}
					return
				}
			}
		}

		// Default behavior: Stream the entire file if no Range header
		c.Status(http.StatusOK)
		_, err = io.Copy(c.Writer, file)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error streaming file")
		}
	})

	router.Run(config.Conf.Port)
}
