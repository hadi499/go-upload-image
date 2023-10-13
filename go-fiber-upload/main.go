package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./", ".html")
	app := fiber.New(fiber.Config{Views: engine})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", nil)
	})

	app.Post("upload", func(c *fiber.Ctx) error {
		var input struct {
			Nama_Gambar string
		}
		if err := c.BodyParser(&input); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		gambar, err := c.FormFile("gambar")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		// mengambil nama gambar
		fmt.Printf("Nama file : %s \n", gambar.Filename)
		// mengambil ukuran file
		fmt.Printf("Ukuran file in byte : %d \n", gambar.Size)

		// mengambil ekstensi file
		splitDots := strings.Split(gambar.Filename, ".")
		ext := splitDots[len(splitDots)-1]
		fmt.Println(ext)
		namaFileBaru := fmt.Sprintf("%s.%s", time.Now().Format("20060102150405"), ext)
		fmt.Println(namaFileBaru)

		// mengambil ukuran gambar
		fileHeader, _ := gambar.Open()
		defer fileHeader.Close()

		imageConfig, _, err := image.DecodeConfig(fileHeader)
		if err != nil {
			log.Print(err)
		}

		width := imageConfig.Width
		height := imageConfig.Height
		fmt.Printf("width %d \n", width)
		fmt.Printf("height %d \n", height)

		// membuat folder upload
		folderUpload := filepath.Join(".", "uploads")
		if err := os.MkdirAll(folderUpload, 0770); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		// menyimpan gambar ke folder
		if err := c.SaveFile(gambar, "./uploads/"+namaFileBaru); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		return c.Status(500).JSON(fiber.Map{
			"title":       input.Nama_Gambar,
			"nama_gambar": namaFileBaru,
			"message":     "gambar berhasil diupload",
		})

	})

	app.Listen(":8080")
}
