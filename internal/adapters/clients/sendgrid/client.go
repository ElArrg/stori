package sendgrid

type Client interface {
	SendEmailV3(toEmail string, toName string, templateId string, payload map[string]any) error
}
