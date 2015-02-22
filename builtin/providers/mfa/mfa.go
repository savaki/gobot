package mfa

import (
	"fmt"

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
		},
	}
}

func registerMFA(c *gobot.Context) {
	log.Debugf("registering mfa")
	c.Respond(fmt.Sprintf("registering a %s mfa", c.Match(1)))

	otp, err := hotp.GenerateHOTP(6, false)
	if err != nil {
		c.Fail(err)
		return
	}

	otp.QR(c.User)
}
