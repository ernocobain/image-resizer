// File: image-resizer/main.go (VERSI FINAL & BENAR)
package main

import (
	"bytes"
	"image"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

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
		// Hapus garis miring di awal path jika ada
		cleanSrc := strings.TrimPrefix(src, "/")

		// Bangun URL Cloud Storage yang benar dengan path lengkap
		storageURL := "https://storage.googleapis.com/maunguli-assets/" + cleanSrc
		// ------------------------------------

		resp, err := http.Get(storageURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Printf("Gagal mengambil gambar dari source: %s, error: %v", storageURL, err)
			return c.Status(http.StatusNotFound).SendString("Image not found from source")
		}
		defer resp.Body.Close()

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
		c.Set("Cache-Control", "public, max-age=31536000, immutable")

		if format == "webp" {
			err = webp.Encode(buf, img, &webp.Options{Quality: 80})
			c.Set("Content-Type", "image/webp")
		} else { // Default ke JPEG
			err = imaging.Encode(buf, img, imaging.JPEG, imaging.JPEGQuality(80))
			c.Set("Content-Type", "image/jpeg")
		}

		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Encode failed")
		}

		return c.Send(buf.Bytes())
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Fatal(app.Listen(":" + port))
}
