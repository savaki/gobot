package mfa

import (
	"bytes"
	"fmt"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/gokyle/hotp"

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

	otp, err := hotp.GenerateHOTP(6, true)
	if err != nil {
		c.Fail(err)
		return
	}

	qrCode, err := otp.QR("Gobot")
	if err != nil {
		c.Fail(err)
		return
	}

	saveOtp(c.User, otp)

	c.Upload(gobot.Attachment{
		Title:       "QR Code",
		Filename:    "QR.png",
		Content:     bytes.NewReader(qrCode),
		ContentType: "image/png",
	})
}

func verify(c *gobot.Context) {
	log.Debugf("verifying mfa code")

	otp, err := loadOtp(c.User)
	if err != nil {
		c.Fail(err)
		return
	}

	code := c.Match(1)
	if otp.Check(code) {
		c.Respond("MFA code valid")
	} else {
		c.Respond("invalid MFA code")
	}
}

var keys map[string]*hotp.HOTP = map[string]*hotp.HOTP{}

var mutex sync.Mutex

func saveOtp(user string, otp *hotp.HOTP) {
	mutex.Lock()
	defer mutex.Unlock()

	keys[user] = otp
}

func loadOtp(user string) (*hotp.HOTP, error) {
	mutex.Lock()
	defer mutex.Unlock()

	otp, found := keys[user]
	if !found {
		return nil, fmt.Errorf("no otp associated with user, %s", user)
	}

	return otp, nil
}
