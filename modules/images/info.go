package images

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"regexp"

	"userstyles.world/modules/config"
	"userstyles.world/modules/log"
)

func fixRawURL(url string) string {
	re := regexp.MustCompile(`(?mi)^(http.*)/(raw|src|blob)/(.*.(png|jpe?g|avif|webp))(\?.*)*$`)
	return re.ReplaceAllString(url, "${1}/raw/${3}")
}

func GenerateImagesForStyle(id, preview string, isOriginalLocal bool) error {
	template := path.Join(config.ImageDir, id)
	original := template + ".original"
	jpeg := template + ".jpeg"
	webp := template + ".webp"

	// Is the preview image not a local file?
	// Let's download it.
	if !isOriginalLocal {
		preview = fixRawURL(preview)

		req, err := http.Get(preview)
		if err != nil {
			log.Warn.Printf("Failed to fetch image URL for %v: %v\n", id, err)
			return err
		}
		defer req.Body.Close()

		// Make sure to received the full body before processing it.
		data, _ := io.ReadAll(req.Body)
		err = os.WriteFile(original, data, 0o600)
		if err != nil {
			log.Warn.Printf("Failed to process image for %v: %v\n", id, err)
			return err
		}
	}

	if err := decodeImage(original, jpeg, ImageTypeJPEG); err != nil {
		log.Warn.Printf("Failed to process image for %v: %v\n", id, err)
		return err
	}

	if err := decodeImage(original, webp, ImageTypeWEBP); err != nil {
		log.Warn.Printf("Failed to process image for %v: %v\n", id, err)
		return err
	}

	return nil
}

func Generate(img multipart.File, id, preview string) error {
	if img != nil {
		data, _ := io.ReadAll(img)
		original := path.Join(config.ImageDir, id+".original")
		if err := os.WriteFile(original, data, 0o600); err != nil {
			return err
		}

		if err := GenerateImagesForStyle(id, preview, true); err != nil {
			return err
		}
	} else if preview != "" {
		if err := GenerateImagesForStyle(id, preview, false); err != nil {
			return err
		}
	}

	return nil
}
