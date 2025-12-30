package services

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"

	"qr-generator/utils"

	qrcode "github.com/skip2/go-qrcode"
)

// QRService handles QR code generation business logic
type QRService struct {
	// Empty for now, allows future dependency injection
}

// NewQRService creates a new instance of QRService
func NewQRService() *QRService {
	return &QRService{}
}

// GenerateSimpleQR generates a basic QR code without watermark
func (s *QRService) GenerateSimpleQR(content string, size int) ([]byte, error) {
	qrCode, err := qrcode.Encode(content, qrcode.Medium, size)
	if err != nil {
		return nil, fmt.Errorf("could not generate a QR code: %v", err)
	}
	return qrCode, nil
}

// GenerateQRWithWatermark generates a QR code with a centered watermark
func (s *QRService) GenerateQRWithWatermark(content string, size int, watermark []byte) ([]byte, error) {
	qrCode, err := s.GenerateSimpleQR(content, size)
	if err != nil {
		return nil, err
	}

	qrCode, err = s.addWatermark(qrCode, watermark)
	if err != nil {
		return nil, fmt.Errorf("could not add watermark to QR code: %v", err)
	}

	return qrCode, nil
}

// addWatermark composites a watermark image onto a QR code
func (s *QRService) addWatermark(qrCode []byte, watermarkData []byte) ([]byte, error) {
	qrCodeData, err := png.Decode(bytes.NewBuffer(qrCode))
	if err != nil {
		return nil, fmt.Errorf("could not decode QR code: %v", err)
	}

	watermarkWidth := uint(float64(qrCodeData.Bounds().Dx()) * 0.25)
	watermark, err := utils.ResizeWatermark(bytes.NewBuffer(watermarkData), watermarkWidth)
	if err != nil {
		return nil, fmt.Errorf("Could not resize the watermark image.", err)
	}

	watermarkImage, err := png.Decode(bytes.NewBuffer(watermark))
	if err != nil {
		return nil, fmt.Errorf("could not decode watermark: %v", err)
	}

	var halfQrCodeWidth, halfWatermarkWidth int = qrCodeData.Bounds().Dx() / 2, watermarkImage.Bounds().Dx() / 2
	offset := image.Pt(
		halfQrCodeWidth-halfWatermarkWidth,
		halfQrCodeWidth-halfWatermarkWidth,
	)

	watermarkImageBounds := qrCodeData.Bounds()
	m := image.NewRGBA(watermarkImageBounds)

	draw.Draw(m, watermarkImageBounds, qrCodeData, image.Point{}, draw.Src)
	draw.Draw(
		m,
		watermarkImage.Bounds().Add(offset),
		watermarkImage,
		image.Point{},
		draw.Over,
	)

	watermarkedQRCode := bytes.NewBuffer(nil)
	png.Encode(watermarkedQRCode, m)

	return watermarkedQRCode.Bytes(), nil
}
