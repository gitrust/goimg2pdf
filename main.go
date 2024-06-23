package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
)

const (
	// A4 pdf page width in mm
	PDF_PAGE_WIDTH = 210.0
	// A4 pdf page height in mm
	PDF_PAGE_HEIGHT = 297.0
)

func main() {
	// define output flag
	outputPDF := flag.String("o", "output.pdf", "Output PDF file name")

	// define Usage
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: goimg2pdf [options] jpeg-file [jpeg-file ...]\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		flag.PrintDefaults()
	}

	// parse flags
	flag.Parse()

	// get arguments
	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	generatePdf(outputPDF, args)
}

// generate PDF file for given file arguments
//
// outputPDF - name of output PDF file
// fileArguments - list of image file arguments (wildcard usage possible)
func generatePdf(outputPDF *string, fileArguments []string) {
	// Create a new PDF
	pdf := gofpdf.New("P", "mm", "A4", "")

	// anonymous function to close PDF
	defer func() {
		if err := pdf.OutputFileAndClose(*outputPDF); err != nil {
			log.Fatalf("Error closing PDF file: %v", err)
		}
	}()

	// get all image paths
	var imagePaths []string = argsToFilePaths(fileArguments)

	for _, imagePath := range imagePaths {
		addPdfPage(pdf, imagePath)
	}

	fmt.Println("PDF created:", *outputPDF)
}

// Add an image file as a new page to PDF
func addPdfPage(pdf *gofpdf.Fpdf, imagePath string) error {
	pdf.AddPage()

	img, err := loadImage(imagePath)
	if err != nil {
		log.Printf("Cant load image %s: %v", imagePath, err)
		return err
	}

	// determine image dimensions
	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()

	// calculate scale factor to fit images into PDF pages
	widthRatio := PDF_PAGE_WIDTH / float64(imgWidth)
	heightRatio := PDF_PAGE_HEIGHT / float64(imgHeight)
	ratio := widthRatio
	if heightRatio < widthRatio {
		ratio = heightRatio
	}

	// calculate new image dimensions
	newWidth := float64(imgWidth) * ratio
	newHeight := float64(imgHeight) * ratio

	// calculate X and Y positions to center image
	x := (PDF_PAGE_WIDTH - newWidth) / 2
	y := (PDF_PAGE_HEIGHT - newHeight) / 2

	// add image
	pdf.ImageOptions(imagePath, x, y, newWidth, newHeight, false, gofpdf.ImageOptions{ImageType: "JPG", ReadDpi: true}, 0, "")

	fmt.Println("Added image:", imagePath)

	return nil
}

// Convert all file arguments to existing file paths
func argsToFilePaths(fileArguments []string) []string {
	// get all image paths
	var imagePaths []string
	for _, arg := range fileArguments {
		matches, err := filepath.Glob(arg)
		if err != nil {
			log.Printf("Bad pattern, can't process %s: %v", arg, err)
			continue
		}
		imagePaths = append(imagePaths, matches...)
	}
	return imagePaths
}

// Load image from path
func loadImage(imagePath string) (image.Image, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}
