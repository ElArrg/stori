package sendgrid

import (
	"net/http"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type ClientConfigs struct {
	SenderEmail string `koanf:"sender-email"`
	SenderName  string `koanf:"sender-name"`
	Key         string `koanf:"key"`
	Host        string `koanf:"host"`
	SandboxMode bool   `koanf:"sandbox_mode"`
}

type DefaultClient struct {
	configs *ClientConfigs
}

func NewDefaultClient(configs *ClientConfigs) *DefaultClient {
	return &DefaultClient{
		configs: configs,
	}
}

func (d *DefaultClient) SendEmailV3(toEmail string, toName string, templateId string, payload map[string]any) error {
	request := sendgrid.GetRequest(d.configs.Key, "/v3/mail/send", d.configs.Host)
	request.Method = http.MethodPost

	// Create email
	m := mail.NewV3Mail()
	m.SetFrom(mail.NewEmail(d.configs.SenderName, d.configs.SenderEmail))

	m.SetTemplateID(templateId)

	p := mail.NewPersonalization()
	p.AddTos(mail.NewEmail(toName, toEmail))
	p.DynamicTemplateData = payload

	m.AddPersonalizations(p)

	if d.configs.SandboxMode {
		// Set sandbox mode on
		settings := mail.NewMailSettings()
		settings.SetSandboxMode(mail.NewSetting(true))
		m.SetMailSettings(settings)
	}

	request.Body = mail.GetRequestBody(m)
	response, err := sendgrid.API(request)

	if err != nil || response.StatusCode >= http.StatusBadRequest {
		return err
	}

	return nil
}
