package main

import (
	"os"
	_ "image/png"
	_ "image/jpeg"
	"image"
	"image/draw"
	"image/png"
)

const base string = "hoodie"

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
		frontImage, _, err := image.Decode(front)
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
	frontImage, backImage := loadBaseImages(base)

	finalImage := image.NewNRGBA(backImage.Bounds())

	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		panic("Source file does not exist")
	}

	source, err := os.Open(sourceFile)
	defer source.Close()
	check(err)
	sourceImage, _, err := image.Decode(source)
	check(err)

	draw.Draw(finalImage, finalImage.Bounds(), backImage, image.Pt(0, 0), draw.Over)
	draw.Draw(finalImage, finalImage.Bounds(), sourceImage, image.Pt(0, 0), draw.Over)
	draw.Draw(finalImage, finalImage.Bounds(), frontImage, image.Pt(0, 0), draw.Over)

	png.Encode()
}
