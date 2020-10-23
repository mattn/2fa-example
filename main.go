package main

import (
	"encoding/base32"
	"net/http"
	"net/url"

	dgoogauth "github.com/dgryski/dgoogauth"
	"github.com/labstack/echo"
	qr "rsc.io/qr"
)

func main() {
	// TODO This secret must be really secret!
	secret := []byte{'H', 'e', 'l', 'l', 'o', '!', 0xDE, 0xAD, 0xBE, 0xEF}
	secretBase32 := base32.StdEncoding.EncodeToString(secret)

	otpc := &dgoogauth.OTPConfig{
		Secret:      secretBase32,
		WindowSize:  3,
		HotpCounter: 0,
		// UTC:         true,
	}

	e := echo.New()
	e.GET("/qrcode", func(c echo.Context) error {
		uri, _ := url.Parse("otpauth://totp")
		desc := "すごいサービス"
		account := "user@example.com"
		uri.Path += "/" + url.PathEscape(account)
		params := url.Values{}
		params.Add("secret", secretBase32)
		params.Add("issuer", desc)
		uri.RawQuery = params.Encode()
		code, err := qr.Encode(uri.String(), qr.Q)
		if err != nil {
			panic(err)
		}
		return c.Blob(http.StatusOK, "image/png", code.PNG())
	})
	e.POST("/auth", func(c echo.Context) error {
		token := c.FormValue("code")
		result, err := otpc.Authenticate(token)
		if err != nil || !result {
			return c.Redirect(http.StatusFound, "/?result=bad")
		}
		return c.Redirect(http.StatusFound, "/ok")
	})
	e.GET("/ok", func(c echo.Context) error {
		return c.String(http.StatusOK, "すごい")
	})
	e.Static("/", "assets")
	e.Logger.Fatal(e.Start(":8989"))
}
