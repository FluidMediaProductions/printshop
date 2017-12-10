package main

import (
	"image"
	"image/draw"
	"image/png"
	_ "image/jpeg"
	"github.com/disintegration/imaging"
	"io"
)

func processImage(base baseImage, source io.Reader, out io.Writer) error {

	size := base.backImage.Bounds()

	finalImage := image.NewNRGBA(size)

	sourceImage, _, err := image.Decode(source)
	if err != nil {
		return err
	}

	boundsSize := image.Pt(base.bottomRightBound.X-base.topLeftBound.X,
		base.bottomRightBound.Y-base.topLeftBound.Y)

	sourceImage = imaging.Fit(sourceImage, boundsSize.X, boundsSize.Y, imaging.Lanczos)
	sourceImageSize := sourceImage.Bounds()

	pos := image.Pt((base.topLeftBound.X+boundsSize.X/2)-(sourceImageSize.Max.X/2),
		(base.topLeftBound.Y+boundsSize.Y/2)-(sourceImageSize.Max.Y/2))
	bounds := image.Rect(pos.X, pos.Y, pos.X+boundsSize.X, pos.Y+boundsSize.Y)

	draw.Draw(finalImage, finalImage.Bounds(), base.backImage, image.Pt(0, 0), draw.Over)
	draw.Draw(finalImage, bounds, sourceImage, image.Pt(0, 0), draw.Over)
	draw.Draw(finalImage, finalImage.Bounds(), base.frontImage, image.Pt(0, 0), draw.Over)

	png.Encode(out, finalImage)
	return nil
}
