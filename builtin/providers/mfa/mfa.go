package mfa

import (
	"bytes"
	"encoding/base32"
	"fmt"
	"sync"

	"crypto/rand"

	"code.google.com/p/rsc/qr"
	log "github.com/Sirupsen/logrus"
	"github.com/hgfischer/go-otp"
	"github.com/savaki/gobot"
)

func Provider() *gobot.Provider {
	return &gobot.Provider{
		Name: "mfa",
		Commands: []gobot.Command{
			{
				Grammar: "mfa register (google)",
				Summary: "register a new mfa device using the specified provider",
				Action:  registerMFA,
			},
			{
				Grammar: `mfa verify (\d+)`,
				Summary: "verify a specific MFA code",
				Action:  verify,
			},
		},
	}
}

func registerMFA(c *gobot.Context) {
	log.Debugf("registering mfa")
	c.Respond(fmt.Sprintf("registering a %s mfa", c.Match(1)))

	data := make([]byte, 10)
	if n, err := rand.Read(data); err != nil {
		c.Fail(err)
		return
	} else if n != len(data) {
		c.Fail(fmt.Errorf("read %d random bytes, wanted %d", n, len(data)))
		return
	}
	secret := base32.StdEncoding.EncodeToString(data)

	saveOtp(c.User, secret)

	code, err := qr.Encode("otpauth://totp/Gobot?secret="+secret, qr.Q)
	if err != nil {
		c.Fail(err)
		return
	}

	c.Upload(gobot.Attachment{
		Title:       "QR Code",
		Filename:    "QR.png",
		Content:     bytes.NewReader(code.PNG()),
		ContentType: "image/png",
	})
}

func verify(c *gobot.Context) {
	log.Debugf("verifying mfa code")

	secret, err := loadOtp(c.User)
	if err != nil {
		c.Fail(err)
		return
	}

	totp := &otp.TOTP{Secret: secret}
	if code := c.Match(1); totp.Now().Verify(code) {
		c.Respond("MFA code valid")
	} else {
		c.Respond(fmt.Sprintf("invalid MFA code, expected %s", totp.Get()))
	}
}

var keys map[string]string = map[string]string{}

var mutex sync.Mutex

func saveOtp(user string, secret string) {
	mutex.Lock()
	defer mutex.Unlock()

	keys[user] = secret
}

func loadOtp(user string) (string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	secret, found := keys[user]
	if !found {
		return "", fmt.Errorf("no secret associated with user, %s", user)
	}

	return secret, nil
}
