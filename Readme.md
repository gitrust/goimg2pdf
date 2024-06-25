# Goimg2pdf

This small piece of code sparked out of need to send a bunch of JPEG documents
as a single PDF file to my insurance.

Somehow I didnt get IrfanView to accept installed PDF plugin and generate Multipage
images. And I couldn't find any alternative in a hurry without installing new apps.

## File support

Currently JPEG and PNG files are supported only.
File extensions (.jpg, .jpeg, .png).

## Usage

![](goimg2pdf.gif)

```
Usage: goimg2pdf [options] img-file [img-file ...]

Options:
  -o string
        Output PDF file name (default "output.pdf")
```

## Examples

```
goimg2pdf 1.jpeg 2.png

Add image: 1.jpeg
Add image: 2.png
PDF created: output.pdf
```

```
goimg2pdf images/*.jpg

Add image: images/1.jpg
Add image: images/2.jpg
Add image: images/3.jpg
Add image: images/4.jpg
Add image: images/5.jpg
PDF created: output.pdf
```

# Reference

Project uses fpdf library from https://github.com/jung-kurt/gofpdf