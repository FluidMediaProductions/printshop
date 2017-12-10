package main

import (
	"os"
	_ "image/png"
	_ "image/jpeg"
	"image"
	"image/draw"
	"image/png"
    "github.com/disintegration/imaging"
	"encoding/json"
	"io/ioutil"
	"fmt"
)

const baseName = "hoodie"
const sourceFile = "source.png"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type baseImage struct {
	backImage image.Image
	frontImage image.Image
	topLeftBound image.Point
	bottomRightBound image.Point
}

func loadImageConfig(name string) baseImage {
	jsonText, err := ioutil.ReadFile("config/"+name+".json")
	check(err)

	var dat map[string]interface{}

	if err := json.Unmarshal([]byte(jsonText), &dat); err != nil {
		panic(err)
	}

	backImage, frontImage := loadBaseImages(dat["back"].(string), dat["front"].(string))

	topLeftBound := image.Pt(int(dat["topLeft"].(map[string]interface{})["X"].(float64)),
		int(dat["topLeft"].(map[string]interface{})["Y"].(float64)))
	bottomRightBound := image.Pt(int(dat["bottomRight"].(map[string]interface{})["X"].(float64)),
		int(dat["bottomRight"].(map[string]interface{})["Y"].(float64)))

	return baseImage{
		backImage: backImage,
		frontImage: frontImage,
		topLeftBound: topLeftBound,
		bottomRightBound: bottomRightBound,
	}
}

func loadBaseImages(backName string, frontName string) (image.Image, image.Image) {
	backFile := "imgs/"+backName
	frontFile := "imgs/"+frontName
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
	base := loadImageConfig(baseName)

	size := base.backImage.Bounds()

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

	boundsSize := image.Pt(base.bottomRightBound.X-base.topLeftBound.X,
		base.bottomRightBound.Y-base.topLeftBound.Y)

	sourceImage = imaging.Fit(sourceImage, boundsSize.X, boundsSize.Y, imaging.Lanczos)
	sourceImageSize := sourceImage.Bounds()

	pos := image.Pt((base.topLeftBound.X+boundsSize.X/2)-(sourceImageSize.Max.X/2),
		(base.topLeftBound.Y+boundsSize.Y/2)-(sourceImageSize.Max.Y/2))
	bounds := image.Rect(pos.X, pos.Y, pos.X+boundsSize.X, pos.Y+boundsSize.Y)

	draw.Draw(finalImage, finalImage.Bounds(), base.backImage, image.Pt(0, 0), draw.Over)
	draw.Draw(finalImage, bounds, sourceImage, image.Pt(0,0), draw.Over)
	draw.Draw(finalImage, finalImage.Bounds(), base.frontImage, image.Pt(0, 0), draw.Over)

	png.Encode(output, finalImage)
}
