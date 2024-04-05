package notifications

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/elarrg/stori/ledger/internal/models"
	"github.com/elarrg/stori/ledger/internal/repository"
	"github.com/elarrg/stori/ledger/internal/service/notifications/dispatchers"
)

type Option func(*DefaultService)

func WithEmailDispatcher(dispatcher dispatchers.Dispatcher) Option {
	return func(service *DefaultService) {
		service.registerDispatcher(models.EmailChannel, dispatcher)
	}
}

type DefaultService struct {
	notificationsRepo repository.Notifications
	accountRepo       repository.Accounts

	channelDispatcher map[models.Channel]dispatchers.Dispatcher
}

func NewDefaultService(nr repository.Notifications, ar repository.Accounts, options ...Option) *DefaultService {
	ds := &DefaultService{
		notificationsRepo: nr,
		accountRepo:       ar,
		channelDispatcher: make(map[models.Channel]dispatchers.Dispatcher),
	}

	for _, opt := range options {
		opt(ds)
	}

	return ds
}

func (d *DefaultService) SendNotification(ctx context.Context, accountID string, operationName string, payload map[string]any) []error {
	var errs []error

	account, err := d.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return append(errs, errors.New("account not found"))
		}
	}

	activeChannels, err := d.notificationsRepo.GetEnabledChannelsByAccountID(ctx, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return append(errs, errors.New("error getting active channels from account settings"))
	}

	templates, err := d.notificationsRepo.GetActiveTemplatesByOperationAndChannels(ctx, operationName, activeChannels)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return append(errs, errors.New("error getting the templates for the notification channels"))
	}

	for _, tmp := range templates {
		if dispatcher, ok := d.channelDispatcher[tmp.Channel]; ok {
			err := dispatcher.Dispatch(ctx, *account, tmp, payload)
			if err != nil {
				errs = append(errs, err)
			}
		} else {
			errs = append(errs, fmt.Errorf("error processing template for channel %v", tmp.Channel))
		}
	}

	return errs
}

func (d *DefaultService) registerDispatcher(channel models.Channel, dispatcher dispatchers.Dispatcher) {
	d.channelDispatcher[channel] = dispatcher
}
