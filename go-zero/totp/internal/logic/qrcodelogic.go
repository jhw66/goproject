package logic

import (
	"bytes"
	"errors"
	"image/png"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

// OTPAuthURLQRCodePNG 将 otpauth://... 内容编码为 PNG 二维码图片（供网页 <img> 或 App 展示）。
func OTPAuthURLQRCodePNG(otpauthURL string, size int) ([]byte, error) {
	u := strings.TrimSpace(otpauthURL)
	if u == "" {
		return nil, errors.New("otpauthUrl is required")
	}
	if size <= 0 {
		size = 1024
	}

	code, err := qr.Encode(u, qr.M, qr.Auto)
	if err != nil {
		return nil, err
	}
	code, err = barcode.Scale(code, size, size)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, code); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
