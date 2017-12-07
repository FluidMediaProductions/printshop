package main

import (
	"os"
	_ "image/png"
	_ "image/jpeg"
	"image"
	"image/draw"
	"image/png"
    "github.com/disintegration/imaging"
)

const base = "hoodie"
var topLeftBound = image.Pt(164, 107)
var bottomRightBound = image.Pt(387, 315)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadBaseImages(name string) (image.Image, image.Image) {
	backFile := "imgs/"+name+"-b.png"
	frontFile := "imgs/"+name+"-f.png"
	frontExists := true

	if _, err := os.Stat(backFile); os.IsNotExist(err) {
		panic("Base file does not exist")
	}
	if _, err := os.Stat(frontFile); os.IsNotExist(err) {
		frontExists = false
	}

	back, err := os.Open(backFile)
	defer back.Close()
	check(err)
	backImage, _, err := image.Decode(back)
	check(err)
	var frontImage image.Image
	if frontExists {
		front, err := os.Open(frontFile)
		defer front.Close()
		check(err)
		frontImage, _, err = image.Decode(front)
		check(err)
		if frontImage.Bounds() != backImage.Bounds() {
			panic("Front and back images are not the same size")
		}
	} else {
		frontImage = image.NewNRGBA(backImage.Bounds())
	}

	return backImage, frontImage
}

func main() {
	sourceFile := "source.png"
	backImage, frontImage := loadBaseImages(base)

	size := backImage.Bounds()

	finalImage := image.NewNRGBA(size)

	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		panic("Source file does not exist")
	}

	source, err := os.Open(sourceFile)
	defer source.Close()
	check(err)
	sourceImage, _, err := image.Decode(source)
	check(err)

	output, err := os.OpenFile("out.png", os.O_CREATE|os.O_WRONLY, 0644)
	check(err)

	boundsSize := image.Pt(bottomRightBound.X-topLeftBound.X, bottomRightBound.Y-topLeftBound.Y)

	sourceImage = imaging.Fit(sourceImage, boundsSize.X, boundsSize.Y, imaging.Lanczos)
	sourceImageSize := sourceImage.Bounds()

	pos := image.Pt((topLeftBound.X+boundsSize.X/2)-(sourceImageSize.Max.X/2),
		(topLeftBound.Y+boundsSize.Y/2)-(sourceImageSize.Max.Y/2))
	bounds := image.Rect(pos.X, pos.Y, pos.X+boundsSize.X, pos.Y+boundsSize.Y)

	draw.Draw(finalImage, finalImage.Bounds(), backImage, image.Pt(0, 0), draw.Over)
	draw.Draw(finalImage, bounds, sourceImage, image.Pt(0,0), draw.Over)
	draw.Draw(finalImage, finalImage.Bounds(), frontImage, image.Pt(0, 0), draw.Over)

	png.Encode(output, finalImage)
}
