package utils

import (
	"image"

	"github.com/pquerna/otp/totp"
	"github.com/spf13/viper"
)

type totpHelper struct{}

type TOTPData struct {
	Secret string      `json:"-"`
	URL    string      `json:"url"`
	Image  image.Image `json:"-"`
}

func (totpHelper) GenerateTOTP(userId string) (TOTPData, error) {
	appName := viper.GetString("app.name")
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      appName,
		AccountName: userId,
	})
	if err != nil {
		return TOTPData{}, err
	}
	data := TOTPData{
		Secret: key.Secret(),
		URL:    key.URL(),
	}

	image, err := key.Image(200, 200)
	if err != nil {
		return TOTPData{}, err
	}
	data.Image = image
	return data, nil
}

func (totpHelper) ValidateTOTP(passcode, secret string) bool {
	return totp.Validate(passcode, secret)
}

var TOTP = totpHelper{}
