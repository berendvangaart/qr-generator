package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"qr-generator/services"
	"qr-generator/utils"
)

// QRHandler handles HTTP requests for QR code generation
type QRHandler struct {
	Service *services.QRService
}

// NewQRHandler creates a new instance of QRHandler with dependency injection
func NewQRHandler(service *services.QRService) *QRHandler {
	return &QRHandler{
		Service: service,
	}
}

// GenerateQRCode is the HTTP handler for /generate endpoint
func (h *QRHandler) GenerateQRCode(writer http.ResponseWriter, request *http.Request) {
	request.ParseMultipartForm(10 << 20)
	var size, content string = request.FormValue("size"), request.FormValue("content")
	var codeData []byte

	if content == "" {
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(
			"Could not determine the desired QR code content.",
		)
		return
	}

	qrCodeSize, err := strconv.Atoi(size)
	if err != nil || size == "" {
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(
			"Could not determine the desired QR code size.",
		)
		return
	}

	watermarkFile, _, err := request.FormFile("watermark")
	if err != nil && errors.Is(err, http.ErrMissingFile) {
		codeData, err = h.Service.GenerateSimpleQR(content, qrCodeSize)
		if err != nil {
			writer.WriteHeader(400)
			json.NewEncoder(writer).Encode(
				fmt.Sprintf("Could not generate QR code. %v", err),
			)
			return
		}
		writer.Header().Add("Content-Type", "image/png")
		writer.Write(codeData)
		return
	}

	watermark, err := utils.UploadFile(watermarkFile)
	if err != nil {
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(
			fmt.Sprint("Could not upload the watermark image.", err),
		)
		return
	}

	contentType := http.DetectContentType(watermark)
	if err != nil {
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(
			fmt.Sprintf(
				"Provided watermark image is a %s not a PNG. %v.", err, contentType,
			),
		)
		return
	}

	codeData, err = h.Service.GenerateQRWithWatermark(content, qrCodeSize, watermark)
	if err != nil {
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(
			fmt.Sprintf(
				"Could not generate QR code with the watermark image. %v", err,
			),
		)
		return
	}

	writer.Header().Set("Content-Type", "image/png")
	writer.Write(codeData)
}
