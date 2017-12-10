package main

import (
	"image"
	"os"
	"io/ioutil"
	"encoding/json"
)

type baseImage struct {
	backImage image.Image
	frontImage image.Image
	topLeftBound image.Point
	bottomRightBound image.Point
}

func loadImageConfig(name string) (baseImage, error) {
	jsonText, err := ioutil.ReadFile("config/"+name+".json")
	if err != nil {
		return baseImage{}, err
	}

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
	}, nil
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
