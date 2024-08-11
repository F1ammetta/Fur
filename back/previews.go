package back

import (
	// "bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	// "github.com/h2non/bimg"
	"github.com/nfnt/resize"
)

var targetDir string // Replace with your target directory

func createPreviews(sourceDir string) error {
	targetDir = sourceDir
	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.Contains(path, "previews") {
			return nil
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if the file is an image or video
		if isImageFile(path) {
			// Generate preview file path in the target directory
			previewPath := generateImagePreviewPath(path)

			// Check if the preview file already exists
			if _, err := os.Stat(previewPath); os.IsNotExist(err) {
				// If not, create the preview
				err := createImagePreview(path, previewPath)
				if err != nil {
					fmt.Printf("Error creating preview for %s: %v\n", path, err)
				}
			} else {
				fmt.Printf("Preview already exists for %s, skipping.\n", path)
			}
		} else if isVideoFile(path) {
			// Generate video preview file path in the target directory
			videoPreviewPath := generateVideoPreviewPath(path)

			// Check if the video preview file already exists
			if _, err := os.Stat(videoPreviewPath); os.IsNotExist(err) {
				// If not, create the video preview
				err := createVideoPreview(path, videoPreviewPath)
				if err != nil {
					fmt.Printf("Error creating video preview for %s: %v\n", path, err)
				}
			} else {
				fmt.Printf("Video preview already exists for %s, skipping.\n", path)
			}
		}

		return nil
	})

	return err
}

func generateImagePreviewPath(sourcePath string) string {
	// Append "_preview.png" to the base name of the image file preserving subdirectories
	imagePreviewName := strings.TrimSuffix(sourcePath, targetDir) + "_preview.png"

	imagePreviewName = strings.Replace(imagePreviewName, targetDir, "", -1)

	// Generate image preview file path in the target directory
	return filepath.Join(targetDir, "previews", imagePreviewName)
}

func generateVideoPreviewPath(sourcePath string) string {
	// Append "_preview.png" to the base name of the video file preserving subdirectories
	videoPreviewName := strings.TrimSuffix(sourcePath, targetDir) + "_preview.png"

	videoPreviewName = strings.Replace(videoPreviewName, targetDir, "", -1)

	// Generate video preview file path in the target directory
	return filepath.Join(targetDir, "previews", videoPreviewName)
}

func createImagePreview(sourcePath, targetPath string) error {
	// Open the image file
	img, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer img.Close()

	// Decode the original image
	originalImage, format, err := image.Decode(img)
	if err != nil {
		return err
	}

	// Resize the image to 300x300
	resized := resize.Resize(300, 300, originalImage, resize.Lanczos3)

	// Create the target directory if it doesn't exist
	err = os.MkdirAll(filepath.Dir(targetPath), os.ModePerm)
	if err != nil {
		return err
	}

	// Save the resized image as the preview
	err = saveImage(targetPath, format, resized)
	if err != nil {
		return err
	}

	fmt.Printf("Image preview created for %s.\n", sourcePath)

	return nil
}

func saveImage(targetPath string, format string, img image.Image) error {
	out, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer out.Close()

	switch format {
	case "jpeg":
		err = jpeg.Encode(out, img, nil)
	case "gif":
		err = gif.Encode(out, img, nil)
	case "png":
		err = png.Encode(out, img)
	}

	return err
}

func createVideoPreview(sourcePath, targetPath string) error {
	// Use ffmpeg to extract the first frame of the video
	cmd := exec.Command("ffmpeg", "-i", sourcePath, "-ss", "00:00:07", "-vframes", "1", "-q:v", "2", targetPath)
	err := cmd.Run()
	if err != nil {
		return err
	}

	// Resize the extracted frame to 300x300
	resizedImage, err := resizeImage(targetPath)
	if err != nil {
		return err
	}

	// Save the resized image as the video preview
	err = saveImage(targetPath, "png", resizedImage)
	if err != nil {
		return err
	}

	fmt.Printf("Video preview created for %s.\n", sourcePath)

	return nil
}
func isImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}

func isVideoFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".mp4" || ext == ".avi" || ext == ".mkv" || ext == ".mov" || ext == ".wmv" || ext == ".m4v" || ext == ".webm"
}

func resizeImage(sourcePath string) (image.Image, error) {
	// Open the image file
	img, err := os.Open(sourcePath)
	if err != nil {
		return nil, err
	}
	defer img.Close()

	// Decode the original image
	originalImage, _, err := image.Decode(img)
	if err != nil {
		return nil, err
	}

	// Resize the image to 300x300
	resized := resize.Resize(300, 300, originalImage, resize.Lanczos3)

	return resized, nil
}
