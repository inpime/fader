package standard

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/flosch/pongo2"
	"github.com/mailgun/mailgun-go"
)

func SendEmail(to, subject, file string, ctx interface{}) (string, error) {
	opt := MainSettings().Config.M("mailer")
	provider := opt.String("provider")
	domain := opt.String("domain")
	privateKey := opt.String("privatekey")
	publicKey := opt.String("publickey")
	sender := opt.String("sender")

	logrus.WithFields(logrus.Fields{
		"_api": addonName,
		"_opt": opt,
	}).Info("options mailer")

	if provider != "mailgun" {
		logrus.WithFields(logrus.Fields{
			"_api": addonName,
		}).Error("supported only 'mailgun' provider, got %q", provider)

		return "", fmt.Errorf("not supported mail provider")
	}

	tpl, err := pongo2.FromFile(file)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"_api": addonName,
			"file": file,
		}).WithError(err).Error("get tempalte file")

		return "", err
	}

	html, err := tpl.Execute(pongo2.Context{
		"ctx": ctx,
	})

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"_api": addonName,
			"file": file,
		}).WithError(err).Error("execute template")
		return "", err
	}

	gun := mailgun.NewMailgun(domain,
		privateKey,
		publicKey)
	m := mailgun.NewMessage(sender, subject, "Message Body (see html)", to)
	m.SetTracking(true)
	m.SetHtml(html)

	_, id, err := gun.Send(m)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"_api": addonName,
		}).WithError(err).Error("Send email")

		return "", err
	}

	logrus.WithFields(logrus.Fields{
		"_api":     addonName,
		"_trackid": id,
		"_to":      to,
	}).Debug("Send email")

	return id, nil
}
