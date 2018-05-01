package nagus

import (
	"image"
	"github.com/golang/freetype"
	"golang.org/x/image/math/fixed"
	"github.com/golang/freetype/truetype"
	"io/ioutil"
	"log"
	"image/color"
	"math/rand"
	"os"
	_ "image/png"
	_ "image/jpeg"
	_ "image/gif"
	"fmt"
	"bytes"
	"image/draw"
)

// Defines a box that we can write text to
type TemplateBox struct {
	Bounds image.Rectangle `json:"bounds"`
	Group string `json:"group"`
}

// Defines a template for images
type MacroTemplate struct {
	Boxes []TemplateBox `json:"boxes"`
	SourceImage image.Image `json:"source_image"`
	SourceMetadata image.Config `json:"source_metadata"`
}

type UnknownType struct {
	Path string
}

func (err *UnknownType) Error() string {
	return fmt.Sprintf("Could not load path %s! Unknown filetype", err.Path)
}

type ImageManager struct {
	ConfigFolder string
}

func NewImageManager(configPath string) *ImageManager {
	return &ImageManager{
		ConfigFolder: configPath,
	}
}

func GetFont() (*truetype.Font, error){
	f,err := ioutil.ReadFile("./Funhouse.ttf")
	if err != nil {
		return &truetype.Font{}, err
	}
	return truetype.Parse(f)
}

func (svc *ImageManager) GetRandomImage() (MacroTemplate, error) {
	imageDirectory := bytes.NewBufferString(svc.ConfigFolder)
	imageDirectory.WriteRune(os.PathSeparator)
	imageDirectory.WriteString("images")
	files, err := ioutil.ReadDir(imageDirectory.String())
	if err != nil {
		return MacroTemplate{}, err
	}
	fileCount := len(files)
	pick := rand.Int()
	if pick < 0 {
		pick = pick * -1
	}
	pick = pick % fileCount
	pickedFile := files[pick]
	pathBuffer := bytes.NewBufferString(imageDirectory.String())
	pathBuffer.WriteRune(os.PathSeparator)
	pathBuffer.WriteString(pickedFile.Name())
	img, err := OpenImage(pathBuffer.String())
	if err != nil {
		return MacroTemplate{}, err
	}
	imgConf, err := GetConfig(pathBuffer.String())
	return MacroTemplate{
		SourceImage: img,
		Boxes: nil,
		SourceMetadata: imgConf,
	}, nil
}

func GetConfig(path string) (image.Config,error) {
	f, err := os.Open(path)
	if err != nil {
		return image.Config{}, err
	}
	img, str, err := image.DecodeConfig(f)
	log.Println(str)
	defer f.Close()
	return img, err
}

func OpenImage(path string) (image.Image, error){
	f, err := os.Open(path)
	if err != nil {
		return image.NewUniform(color.Black), err
	}
	img, str, err := image.Decode(f)
	log.Println(str)
	defer f.Close()
	return img, err
}

func (svc *ImageManager) WriteToImage(phrase string) image.Image {
	img, err := svc.GetRandomImage()
	if err != nil {
		log.Panic(err)
	}
	ctx := freetype.NewContext()
	font, err := GetFont()
	if err != nil {
		log.Fatal(err)
	}
	rect := image.Rect(0,0, img.SourceMetadata.Width, img.SourceMetadata.Height)
	drawImg := image.NewRGBA(rect)
	draw.Draw(drawImg, rect, img.SourceImage, rect.Min, draw.Src)
	ctx.SetFont(font)
	ctx.SetFontSize(15)
	ctx.SetDst(drawImg)
	ctx.SetClip(rect)
	ctx.SetSrc(image.NewUniform(color.RGBA{0,0,0,255}))

	ctx.DrawString(phrase, fixed.Point26_6{
		X: ctx.PointToFixed(15 * 1.5),
		Y: ctx.PointToFixed(15 * 1.5),
	})
	return drawImg
}
