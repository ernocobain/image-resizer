// File: image-resizer/main.go (VERSI REVISI)
package main

import (
	"bytes"
	"image"
	"log"
	"net/http"
	"os"
	"path" // <-- Impor package 'path'
	"strconv"

	webp "github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	app.Get("/cdn", func(c *fiber.Ctx) error {
		src := c.Query("src")
		widthStr := c.Query("w")
		heightStr := c.Query("h")
		format := c.Query("format", "webp")

		if src == "" {
			return c.Status(http.StatusBadRequest).SendString("Missing 'src' parameter")
		}

		// --- PERBAIKAN LOGIKA URL DI SINI ---
		// Ambil hanya nama filenya saja dari path src
		// Contoh: dari "/static/images/logo.png" menjadi "logo.png"
		fileName := path.Base(src)

		// Bangun URL Cloud Storage yang benar
		storageURL := "https://storage.googleapis.com/maunguli-assets/" + fileName
		// ------------------------------------

		resp, err := http.Get(storageURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Printf("Gagal mengambil gambar dari source: %s, error: %v", storageURL, err)
			return c.Status(http.StatusNotFound).SendString("Image not found from source")
		}
		defer resp.Body.Close()

		// ... sisa kode untuk resize dan encode tetap sama ...
		img, _, err := image.Decode(resp.Body)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Failed to decode image")
		}
		width, _ := strconv.Atoi(widthStr)
		height, _ := strconv.Atoi(heightStr)
		if width > 0 || height > 0 {
			img = imaging.Resize(img, width, height, imaging.Lanczos)
		}
		buf := new(bytes.Buffer)
		if format == "webp" {
			err = webp.Encode(buf, img, &webp.Options{Quality: 85})
			c.Set("Content-Type", "image/webp")
		} else {
			err = imaging.Encode(buf, img, imaging.JPEG, imaging.JPEGQuality(85))
			c.Set("Content-Type", "image/jpeg")
		}
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Encode failed")
		}
		c.Set("Cache-Control", "public, max-age=31536000, immutable")
		return c.Send(buf.Bytes())
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Fatal(app.Listen(":" + port))
}
