package qrlogo

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"

	qr "github.com/skip2/go-qrcode"
)

// Encoder defines settings for QR/Overlay encoder.
type Encoder struct {
	AlphaThreshold int
	GreyThreshold  int
	QRLevel        qr.RecoveryLevel
}

// DefaultEncoder is the encoder with default settings.
var DefaultEncoder = Encoder{
	AlphaThreshold: 2000,       // FIXME: don't remember where this came from
	GreyThreshold:  30,         // in percent
	QRLevel:        qr.Highest, // recommended, as logo steals some redundant space
}

// Encode encodes QR image, adds logo overlay and renders result as PNG.
func Encode(str string, logo image.Image, size int) (*bytes.Buffer, error) {
	return DefaultEncoder.Encode(str, logo, size)
}

// Encode encodes QR image, adds logo overlay and renders result as PNG.
func (e Encoder) Encode(str string, logo image.Image, size int) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	code, err := qr.New(str, e.QRLevel)
	if err != nil {
		return nil, err
	}

	qrCode := code.Image(size)

	offsetX := qrCode.Bounds().Max.X - logo.Bounds().Max.X
	offsetY := qrCode.Bounds().Max.Y - logo.Bounds().Max.Y

    offset := image.Pt(offsetX, offsetY)
    qrCodeBounds := qrCode.Bounds()

	output := image.NewRGBA(qrCodeBounds)
    draw.Draw(output, qrCodeBounds, qrCode, image.ZP, draw.Src)
    draw.Draw(output, logo.Bounds().Add(offset), logo, image.ZP, draw.Over) 

	err = png.Encode(&buf, output)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
