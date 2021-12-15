package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
)

type pixel struct {
	x     int
	y     int
	green int
	red   int
	blue  int
}

func init() {
	// damn important or else At(), Bounds() functions will
	// caused memory pointer error!!
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
}

var imageByPixel [256]pixel
var iter int
var red uint8
var green uint8
var blue uint8
var alpha uint8

var imagePath string = "./synthwave.jpeg"

func main() {
	iter = 0
	imgfile, err := os.Open(imagePath)

	if err != nil {
		fmt.Println("img.jpg file not found!")
		os.Exit(1)
	}

	defer imgfile.Close()

	// get image height and width with image/jpeg
	// change accordinly if file is png or gif

	imgCfg, _, err := image.DecodeConfig(imgfile)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	width := imgCfg.Width
	height := imgCfg.Height

	fmt.Println("Width : ", width)
	fmt.Println("Height : ", height)

	// we need to reset the io.Reader again for image.Decode() function below to work
	// otherwise we will  - panic: runtime error: invalid memory address or nil pointer dereference
	// there is no build in rewind for io.Reader, use Seek(0,0)
	imgfile.Seek(0, 0)

	// get the image
	img, _, err := image.Decode(imgfile)

	// Create an widght by height image
	imgNew := image.NewRGBA(image.Rect(0, 0, width, height))

	fmt.Println(img.At(10, 10).RGBA())

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r, g, b, a := img.At(x, y).RGBA()

			// create array for ESP32 image decoding
			imageByPixel[iter].x = x
			imageByPixel[iter].y = y
			imageByPixel[iter].green = int((((1 - a) * g) + (a * g)) / 255) //green decoding into the RGB from RGBA
			imageByPixel[iter].red = int((((1 - a) * r) + (a * r)) / 255)   //red decoding into the RGB from RGBA
			imageByPixel[iter].blue = int((((1 - a) * b) + (a * b)) / 255)  //blue decoding into the RGB from RGBA

			imgNew.Set(x, y, color.RGBA{uint8(imageByPixel[iter].red), uint8(imageByPixel[iter].green), uint8(imageByPixel[iter].blue), uint8(a)}) //do reverse image reconstruction

			fmt.Printf("{%v, %v, %v, %v, %v}, \n", x, y, imageByPixel[iter].green, imageByPixel[iter].red, imageByPixel[iter].blue) //esp32 ready
			iter++
		}
	}

	// Save to out.png
	f, _ := os.OpenFile("out.png", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	png.Encode(f, imgNew)

}
