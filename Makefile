

build:
	go build -o goimg2pdf.exe

setup:
	go mod init goimg2pdf
	go mod tidy