package yandex

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"

	"github.com/liyue201/goqr"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

func decodeSVG(svg []byte) (image.Image, error) {
	icon, err := oksvg.ReadIconStream(bytes.NewReader(svg))
	if err != nil {
		return nil, fmt.Errorf("Unable to read icon stream: %w", err)
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
		return url, fmt.Errorf("Unable to decode SVG: %w", err)
	}
	qrCodes, err := goqr.Recognize(img)
	if err != nil {
		return url, fmt.Errorf("Unable to recognize QR: %w", err)
	}
	if len(qrCodes) != 1 {
		return url, fmt.Errorf("Invalid QRCode")
	}
	url = string(qrCodes[0].Payload)
	return url, nil
}
