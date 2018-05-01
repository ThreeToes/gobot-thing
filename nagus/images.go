package nagus

import (
	"image"
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
	"path/filepath"
	"encoding/json"
	"strings"
	"github.com/fogleman/gg"
)

type TemplateBounds struct {
	X int `json:"x"`
	Y int `json:"y"`
	Width int `json:"width"`
	Height int `json:"height"`
}

// Defines a box that we can write text to
type TemplateBox struct {
	Bounds TemplateBounds `json:"bounds"`
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
	for ;strings.HasSuffix(files[pick].Name(), "json"); {
		pick = (pick + 1) % fileCount
	}
	pickedFile := files[pick]
	pathBuffer := bytes.NewBufferString(imageDirectory.String())
	pathBuffer.WriteRune(os.PathSeparator)
	pathBuffer.WriteString(pickedFile.Name())
	img, err := OpenImage(pathBuffer.String())
	if err != nil {
		return MacroTemplate{}, err
	}
	imgConf, err := GetConfig(pathBuffer.String())
	if err != nil {
		return MacroTemplate{}, err
	}
	boxConf, err := svc.GetImageConfig(pathBuffer.String())
	return MacroTemplate{
		SourceImage: img,
		Boxes: boxConf,
		SourceMetadata: imgConf,
	}, nil
}

func (svc *ImageManager) GetImageConfig(path string) ([]TemplateBox, error) {
	ext := filepath.Ext(path)
	fileBuf := bytes.NewBufferString(path[0: len(path) - len(ext)])
	fileBuf.WriteString(".json")
	contents, err := ioutil.ReadFile(fileBuf.String())
	var template []TemplateBox
	if err != nil {
		return template, err
	}
	err = json.Unmarshal(contents, &template)
	return template, err
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
	sess := gg.NewContext(img.SourceMetadata.Width, img.SourceMetadata.Height)
	sess.DrawImage(img.SourceImage, 0, 0)
	xPos := float64(img.Boxes[0].Bounds.X)
	yPos := float64(img.Boxes[0].Bounds.Y)
	sess.SetColor(color.Black)
	sess.DrawStringWrapped(phrase, xPos, yPos, 0.0, 0.0, float64(img.Boxes[0].Bounds.Width), 1, gg.AlignLeft)

	return sess.Image()
}
