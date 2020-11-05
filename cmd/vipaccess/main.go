package main

import (
	"flag"
	"fmt"
	"github.com/asticode/go-astikit"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/asticode/go-astilectron"
	"github.com/athomason/vipaccess-go/vipaccess"
)

func main() {
	var (
		png         = flag.String("png", "", "png path to write QR code")
		accountName = flag.String("account-name", "", "Account name (default: VIP Access)")
		issuer      = flag.String("issuer", "", "Specify the issuer name to use (default: Symantec)")
		secret      = flag.String("show", "", "Generating access codes using an existing credential")
	)
	flag.Parse()

	if flag.NFlag() > 0 {
		var credential = generate(accountName, issuer, secret)
		if credential != nil {
			generateQr(credential, png)
		}
	} else {
		var l = log.New(log.Writer(), log.Prefix(), log.Flags())

		var app, err = astilectron.New(l, astilectron.Options{
			AppName:           "Symantec VIP access",
			BaseDirectoryPath: ".",
		})
		if err != nil {
			l.Fatal(fmt.Errorf("main: creating astilectron failed: %w", err))
		}
		defer app.Close()

		app.HandleSignals()

		if err = app.Start(); err != nil {
			l.Fatal(fmt.Errorf("main: starting astilectron failed: %w", err))
		}

		var w *astilectron.Window
		if w, err = app.NewWindow("resources/index.html", &astilectron.WindowOptions{
			AlwaysOnTop:    astikit.BoolPtr(true),
			Height:         astikit.IntPtr(200),
			Width:          astikit.IntPtr(200),
			UseContentSize: astikit.BoolPtr(true),
		}); err != nil {
			l.Fatal(fmt.Errorf("main: new window failed: %w", err))
		}

		if err = w.Create(); err != nil {
			l.Fatal(fmt.Errorf("main: creating window failed: %w", err))
		}

		w.OnMessage(func(message *astilectron.EventMessage) interface{} {
			var text string
			message.Unmarshal(&text)

			if text == "ping" {
				var _, errStat = os.Stat("key.txt")
				if !os.IsNotExist(errStat) {
					var line, errFile = ioutil.ReadFile("key.txt")
					if (errFile == nil) {
						var accountName = ""
						var issuer = ""
						var secret = strings.Split(string(line), ",")[1]
						var code = generate(&accountName, &issuer, &secret)
						return strings.Split(string(line), ",")[0] + "," + code.SixDigit
					}
				} else {
					return ""
				}
			} else if text == "generate" {
				var accountName = ""
				var issuer = ""
				var secret = ""
				var credential = generate(&accountName, &issuer, &secret)
				var csv = credential.ID + "," + vipaccess.B32(credential.Key)
				ioutil.WriteFile("key.txt", []byte(csv), 0644)

				secret = vipaccess.B32(credential.Key)
				var code = generate(&accountName, &issuer, &secret)
				return string(credential.ID) + "," + code.SixDigit
			}
			return nil
		})

		app.Wait()
	}
}

func generate(accountName *string, issuer *string, secret *string) *vipaccess.Credential {
	p := vipaccess.GenerateRandomParameters()

	if *accountName != "" {
		p.AccountName = *accountName
	}
	if *issuer != "" {
		p.Issuer = *issuer
	}

	if *secret == "" {
		c, err := vipaccess.GenerateCredential(p)
		if err != nil {
			log.Fatal(err)
		}

		if err := c.Validate(); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("OTP credential: %s\nExpires: %s (%s)\n",
			c.URI(), c.Expires, -time.Since(c.Expires))

		return c
	} else {
		var code = vipaccess.GenerateTOTPCode(vipaccess.StringToB32(*secret), time.Now())
		fmt.Printf("Security code valid for 30 seconds:\n%s",
			code)

		var credential = &vipaccess.Credential{SixDigit: code}
		return credential
	}
}

func generateQr(c *vipaccess.Credential, png *string) {
	if *png != "" {
		f, err := os.Create(*png)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := f.Write(c.QRCodePNG()); err != nil {
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("wrote QR code to %s\n", *png)
	}
}