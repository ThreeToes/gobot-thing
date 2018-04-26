package nagus

import (
	"image"
	"github.com/golang/freetype"
	"golang.org/x/image/math/fixed"
	"github.com/golang/freetype/truetype"
	"io/ioutil"
	"log"
	"image/color"
)

func GetFont() (*truetype.Font, error){
	f,err := ioutil.ReadFile("./Funhouse.ttf")
	if err != nil {
		return &truetype.Font{}, err
	}
	return truetype.Parse(f)
}

func WriteToImage(phrase string) image.Image {
	rect := image.Rect(0, 0, 450, 450)
	img := image.NewRGBA(rect)
	for x := 0; x < 450; x++ {
		for y := 0; y < 450; y++ {
			img.Set(x, y, color.RGBA{255,255,255,255})
		}
	}
	ctx := freetype.NewContext()
	font, err := GetFont()
	if err != nil {
		log.Fatal(err)
	}
	ctx.SetFont(font)
	ctx.SetFontSize(15)
	ctx.SetDst(img)
	ctx.SetClip(img.Bounds())
	ctx.SetSrc(image.NewUniform(color.RGBA{0,0,0,255}))

	ctx.DrawString(phrase, fixed.Point26_6{
		X: ctx.PointToFixed(15 * 1.5),
		Y: ctx.PointToFixed(15 * 1.5),
	})
	return img
}
