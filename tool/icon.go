package tool

import (
	"code.google.com/p/graphics-go/graphics"
	"image"
	"image/draw"
	"image/png"
	"os"
)

// func test() {
// 	if err := GenerateImage("src.png", "0.png", "dest.png", 100, 100); err != nil {
// 		log.Fatal(err)
// 	}
// 	return
// }

func MergeImage(backPath, overPath string) (rgba *image.RGBA, err error) {
	var img image.Image
	if img, err = LoadImage(backPath); err != nil {
		return
	}
	rgba = image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Over)

	if img, err = LoadImage(overPath); err != nil {
		return
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Over)
	return
}

func ScaleImageFile(srcImagePath, destImagePath string, newWidth, newHeight int) (err error) {
	var srcImage image.Image
	if srcImage, err = LoadImage(srcImagePath); err != nil {
		return
	}

	destRGBA := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	if err = graphics.Scale(destRGBA, srcImage); err != nil {
		return
	}

	if err = SaveImage(destRGBA, destImagePath); err != nil {
		return
	}

	return
}

func ScaleImage(src *image.RGBA, newWidth, newHeight int) (dest *image.RGBA, err error) {
	dest = image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	if err = graphics.Scale(dest, src); err != nil {
		return
	}
	return
}

func GenerateImage(backImagePath, overImagePath, destImagePath string, newWidth, newHeight int) (err error) {
	var srcRGBA, destRGBA *image.RGBA
	if srcRGBA, err = MergeImage(backImagePath, overImagePath); err != nil {
		return
	}

	if destRGBA, err = ScaleImage(srcRGBA, newWidth, newHeight); err != nil {
		return
	}

	if err = SaveImage(destRGBA, destImagePath); err != nil {
		return
	}
	return
}

func SaveImage(rgba *image.RGBA, path string) (err error) {
	var file *os.File
	if file, err = os.Create(path); err != nil {
		return
	}
	defer file.Close()
	return png.Encode(file, rgba)
}

func LoadImage(path string) (img image.Image, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	return
}
