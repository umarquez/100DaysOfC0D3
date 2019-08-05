package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"net/http"
)

const rndImgApi = "https://source.unsplash.com"

func decodeImageFromReader(reader io.Reader) (image.Image, error) {
	img, err := jpeg.Decode(reader)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func RandomImage(size image.Point) (image.Image, error) {
	url := fmt.Sprintf("%v/%vx%v", rndImgApi, size.X, size.Y)
	resp, err := http.Get(url)
	if err != nil {
		err = fmt.Errorf("error fetching random image from the API, %v", err)
		return nil, err
	}

	img, err := decodeImageFromReader(resp.Body)
	if err != nil {
		err = fmt.Errorf("error decoding random image, %v", err)
		return nil, err
	}

	return img, err
}

func LocalImage(imgPath string) (image.Image, error) {
	// Leemos el contenido del archivo.
	content, err := ioutil.ReadFile(imgPath)
	if err != nil {
		err = fmt.Errorf("error reading local image (%v), %v", imgPath, err)
		return nil, err
	}

	// Creamos un buffer con el contenido del archivo para decodificarlo como JPG
	img, err := decodeImageFromReader(bytes.NewBuffer(content))
	if err != nil {
		err = fmt.Errorf("error decoding local image (%v), %v", imgPath, err)
		return nil, err
	}

	return img, nil
}
