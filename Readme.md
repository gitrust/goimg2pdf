# Goimg2pdf

This small piece of code sparked out of need to send a bunch of JPEG documents
as a single PDF file via email to my insurance.

Somehow I didnt get IrfanView to accept installed PDF plugin and generate Multipage
images. And I couldn't find any alternative in a hurry without installing new apps.

## Usage

```
Usage: goimg2pdf [options] jpeg-file

Options:
  -o string
        Output PDF file name (default "output.pdf")
```

## Examples

```
goimg2pdf 1.jpg 2.jpg

Add image: 1.jpg
Add image: 2.jpg
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

