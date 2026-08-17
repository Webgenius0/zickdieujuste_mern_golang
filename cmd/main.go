package main

import (
	"fmt"
	"log"

	"gotickets/internal/config"
	"gotickets/internal/server"
	"gotickets/internal/upload"
)

func main() {
	cfg := config.LoadEnv()
	db := config.ConnectDatabase(cfg)

	// Initialise the Cloudinary uploader.
	// To swap to another provider (e.g. AWS S3) later, implement upload.Uploader in a new file
	// and change only this wiring — domain code remains untouched.
	var uploader upload.Uploader
	if cfg.CloudinaryCloudName != "" && cfg.CloudinaryApiKey != "" && cfg.CloudinaryApiSecret != "" {
		u, err := upload.NewCloudinaryUploader(cfg.CloudinaryCloudName, cfg.CloudinaryApiKey, cfg.CloudinaryApiSecret)
		if err != nil {
			log.Printf("Warning: Cloudinary initialisation failed (%v) — avatar upload will be unavailable\n", err)
		} else {
			uploader = u
		}
	} else {
		fmt.Println("Warning: CLOUDINARY_CLOUD_NAME/API_KEY/API_SECRET not set — avatar upload disabled")
	}

	server.Start(db, cfg, uploader)
}
