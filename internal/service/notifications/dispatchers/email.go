package dispatchers

import (
	"context"
	"errors"
	"fmt"

	"github.com/elarrg/stori/ledger/internal/adapters/clients/sendgrid"

	"github.com/elarrg/stori/ledger/internal/models"
)

type EmailService struct {
	sendgridClient sendgrid.Client
}

func NewEmailProcessor(sc sendgrid.Client) *EmailService {
	return &EmailService{
		sendgridClient: sc,
	}
}

func (e *EmailService) Dispatch(ctx context.Context, account models.Account, template models.Template, payload map[string]any) error {
	if template.Channel != models.EmailChannel {
		return fmt.Errorf("email: cannot process the template channel %v", template.Channel)
	}

	if !template.Active {
		return errors.New("email: template is not active")
	}

	var err error
	switch template.Operation {
	case AccountSummaryOp:
		err = e.accountSummaryHandler(ctx, account, template, payload)
	default:
		err = errors.New("email: operation not supported")
	}

	if err != nil {
		return err
	}

	return nil
}

func (e *EmailService) accountSummaryHandler(ctx context.Context, account models.Account, template models.Template, payload map[string]any) error {
	payload["name"] = account.Firstname

	var err error
	switch template.SourceType {
	case "sendgrid":
		err = e.sendgridClient.SendEmailV3(account.Email, account.Firstname, template.Source, payload)
	}

	if err != nil {
		// todo log
		return errors.New("email: error while sending account summary notification")
	}

	return nil
}
