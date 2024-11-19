// package main

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"server/config"
// 	"strings"

// 	"github.com/gofiber/fiber/v2"
// )

// func main() {
// 	err := config.Init()
// 	if err != nil {
// 		return
// 	}
// 	app := fiber.New(
// 		fiber.Config{
// 			CaseSensitive: true,
// 			Concurrency:   config.Conf.Concurrency,
// 		},
// 	)

// 	// Routes
// 	type FileInfo struct {
// 		FileName string  `json:"file_name"`
// 		FileSize float64 `json:"file_size"` // Size in bytes
// 	}
// 	app.Get("/list", func(c *fiber.Ctx) (err error) {
// 		wd, err := os.Getwd()
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).SendString("Could not read directory")
// 		}
// 		path := wd + "/media"
// 		dirs, err := os.ReadDir(path)
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).SendString("Could not read directory")
// 		}
// 		var fileDetails []FileInfo
// 		for _, file := range dirs {
// 			if !file.IsDir() {
// 				// Get the file's stat to retrieve the size
// 				filePath := filepath.Join(path, file.Name())
// 				fileInfo, err := os.Stat(filePath)
// 				if err != nil {
// 					return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Could not stat file %s", file.Name()))
// 				}

// 				// Append the file name and size to the response slice
// 				fileDetails = append(fileDetails, FileInfo{
// 					FileName: file.Name(),
// 					FileSize: float64(fileInfo.Size()) / (1024 * 1024),
// 				})
// 			}
// 		}
// 		return c.JSON(fileDetails)
// 	})
// 	app.Get("/stream/:file", func(c *fiber.Ctx) (err error) {
// 		wd, err := os.Getwd()
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).SendString("Could not read directory")
// 		}
// 		path := wd + "/media/"
// 		fileName := c.Params("file")
// 		// Open the file
// 		file, err := os.Open(path + fileName)
// 		if err != nil {
// 			return c.Status(fiber.StatusNotFound).SendString("File not found")
// 		}
// 		defer file.Close()
// 		ext := strings.ToLower(filepath.Ext(fileName))
// 		var contentType string

// 		switch ext {
// 		case ".jpg", ".jpeg":
// 			contentType = "image/jpeg" // For JPEG images
// 		case ".png":
// 			contentType = "image/png" // For PNG images
// 		case ".gif":
// 			contentType = "image/gif" // For GIF images
// 		default:
// 			contentType = "application/octet-stream" // Default for unknown formats
// 		}

// 		// Set the correct Content-Type for the image file
// 		c.Set("Content-Type", contentType)
// 		// Ensure the file is displayed inline (i.e., not as an attachment)
// 		c.Set("Content-Disposition", "inline; filename="+fileName)
// 		// Enable chunked transfer encoding
// 		c.Set("Transfer-Encoding", "chunked")

// 		// Get the file size
// 		fileInfo, err := file.Stat()
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).SendString("Unable to get file info")
// 		}
// 		fileSize := fileInfo.Size()

// 		// Define the chunk size (e.g., 1MB chunks)
// 		chunkSize := int64(1024 * 1024) // 1MB chunk size
// 		var start int64 = 0

// 		// Loop to send the image in chunks
// 		for start < fileSize {
// 			// Determine the end byte of the current chunk
// 			end := start + chunkSize
// 			if end > fileSize {
// 				end = fileSize
// 			}

// 			// Seek to the start byte for this chunk
// 			_, err := file.Seek(start, 0)
// 			if err != nil {
// 				log.Printf("Error seeking to byte %d: %v", start, err)
// 				return c.Status(fiber.StatusInternalServerError).SendString("Error seeking in file")
// 			}

// 			// Read the chunk into a buffer
// 			buf := make([]byte, end-start)
// 			_, err = file.Read(buf)
// 			if err != nil {
// 				return c.Status(fiber.StatusInternalServerError).SendString("Error reading file")
// 			}

// 			// Send the chunk to the client
// 			err = c.Send(buf)
// 			if err != nil {
// 				return c.Status(fiber.StatusInternalServerError).SendString("Error sending file chunk")
// 			}

// 			// Update the start byte for the next chunk
// 			start = end

//				// Wait for 5 seconds before sending the next chunk
//				// time.Sleep(1 * time.Second)
//			}
//			return
//		})
//		err = app.Listen(":" + config.Conf.Port)
//		if err != nil {
//			panic(err)
//		}
//	}
package main

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to our video streaming platform!")
	})

	router.GET("/stream/:file", func(c *gin.Context) {
		wd, err := os.Getwd()
		if err != nil {
			c.String(http.StatusOK, "Welcome to our video streaming platform!")
			return
		}
		path := wd + "/media/"
		fileName := c.Param("file")
		// Open the file
		file, err := os.Open(path + fileName)
		if err != nil {
			c.String(http.StatusNotFound, "File not found")
			return
		}
		defer file.Close()

		c.Header("Content-Type", "video/mp4")
		buffer := make([]byte, 4096*4096)
		io.CopyBuffer(c.Writer, file, buffer)
	})
	router.Run(":9944")
}
