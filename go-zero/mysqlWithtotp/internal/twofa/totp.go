package twofa

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image/png"
	"strings"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type Enrollment struct {
	Secret        string
	AccountName   string
	Issuer        string
	OTPAuthURL    string
	QRCodeDataURL string
}

func GenerateEnrollment(issuer, accountName string, qrCodeSize int) (*Enrollment, error) {
	issuer = strings.TrimSpace(issuer)
	accountName = strings.TrimSpace(accountName)
	if issuer == "" || accountName == "" {
		return nil, errors.New("issuer and accountName are required")
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
	})
	if err != nil {
		return nil, err
	}

	dataURL, err := qrCodeDataURL(key.URL(), qrCodeSize)
	if err != nil {
		return nil, err
	}

	return &Enrollment{
		Secret:        key.Secret(),
		AccountName:   accountName,
		Issuer:        issuer,
		OTPAuthURL:    key.URL(),
		QRCodeDataURL: dataURL,
	}, nil
}

func ValidateCode(secret, code string) bool {
	secret = strings.TrimSpace(secret)
	code = strings.TrimSpace(code)
	if secret == "" || code == "" {
		return false
	}

	valid, err := totp.ValidateCustom(code, secret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	return err == nil && valid
}

func qrCodeDataURL(otpauthURL string, qrCodeSize int) (string, error) {
	if qrCodeSize <= 0 {
		qrCodeSize = 256
	}

	code, err := qr.Encode(otpauthURL, qr.M, qr.Auto)
	if err != nil {
		return "", err
	}
	code, err = barcode.Scale(code, qrCodeSize, qrCodeSize)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, code); err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
