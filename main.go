package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

var validExtensions = []string{".jpg", ".png", ".jpeg"}

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
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: goimg2pdf [options] img-file [img-file ...]\n\n")
		fmt.Println("Accepted image extensions: [*.jpg, *.png, *,jpeg]")
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
		if !isValidExtension(filepath.Ext(imagePath)) {
			log.Println("Unsupported file ", imagePath)
			continue
		}
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

	// calculate new image dimensions to fit into page
	newWidth, newHeight := calcPageDimensions(img)

	// calculate X and Y positions to center image
	x, y := calcCenterPosition(newWidth, newHeight)

	// add image
	pdf.ImageOptions(imagePath, x, y, newWidth, newHeight, false, gofpdf.ImageOptions{ReadDpi: true}, 0, "")

	fmt.Println("Added image:", imagePath)

	return nil
}

// Funktion, um zu prüfen, ob eine Erweiterung in der Liste enthalten ist
func isValidExtension(extension string) bool {
	for _, ext := range validExtensions {
		if strings.EqualFold(ext, extension) {
			return true
		}
	}
	return false
}

// calculate new dimensions to fit image into page
func calcPageDimensions(img image.Image) (width, height float64) {
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
	return newWidth, newHeight
}

// calculate x and y center position of img
func calcCenterPosition(width, height float64) (x, y float64) {
	// calculate X and Y positions to center image
	x = (PDF_PAGE_WIDTH - width) / 2
	y = (PDF_PAGE_HEIGHT - height) / 2

	return x, y
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

	// loads image into memory
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// interface for image operations
type ImageFilter interface {
	Apply(img image.Image) image.Image
}

type GrayscaleFilter struct{}

func (f GrayscaleFilter) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			grayImg.Set(x, y, color.GrayModel.Convert(img.At(x, y)))
		}
	}
	return grayImg
}

type BlackWhiteFilter struct {
	Threshold uint8
}

func (f BlackWhiteFilter) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	bwImg := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			if c.Y > f.Threshold {
				bwImg.Set(x, y, color.White)
			} else {
				bwImg.Set(x, y, color.Black)
			}
		}
	}
	return bwImg
}
