package yandex

import (
	"bytes"
	"image"
	"image/draw"

	"github.com/liyue201/goqr"
	"github.com/pkg/errors"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

func decodeSVG(svg []byte) (image.Image, error) {
	icon, err := oksvg.ReadIconStream(bytes.NewReader(svg))
	if err != nil {
		return nil, err
	}
	w := 265
	h := 265
	icon.SetTarget(0, 0, float64(w), float64(h))
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(rgba, rgba.Bounds(), image.White, image.ZP, draw.Src)
	icon.Draw(rasterx.NewDasher(w, h, rasterx.NewScannerGV(w, h, rgba, rgba.Bounds())), 1)

	return rgba, nil
}

func decodeQR(svg []byte) (url string, err error) {
	img, err := decodeSVG(svg)
	if err != nil {
		return url, err
	}
	qrCodes, err := goqr.Recognize(img)
	if err != nil {
		return url, err
	}
	if len(qrCodes) != 1 {
		return url, errors.New("YandexAPI.decodeQR: Invalid QRCode")
	}
	url = string(qrCodes[0].Payload)
	return url, nil
}
