package background_service

import (
	"context"
	"payments_service/notifications"
	"payments_service/storage"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type BackgroundService struct {
	storage storage.StorageActions
	notify  notifications.NotificationsStruct
}

func NewBackgroundService(storage storage.StorageActions) *BackgroundService {
	return &BackgroundService{
		storage: storage,
		notify:  notifications.NotificationsStruct{},
	}
}

func (b *BackgroundService) CheckStatuses(ctx context.Context) error {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		case <-ticker.C:
		}

		log.Info().Msg("start check")
		storagePayments, err := b.storage.GetPendingPayments()
		if err != nil {
			return err
		}
		log.Info().Msgf("%d Objects found", len(storagePayments))
		for _, p := range storagePayments {
			b.PaymentsStateMachine(p)
			//fmt.Println(p)
		}
	}

	return nil
}

func (b *BackgroundService) PaymentsStateMachine(payment storage.Payment) {
	now := time.Now()

	switch strings.ToLower(payment.State) {
	case "new":
		log.Info().Msgf("Start 'New' with payment %d", payment.ID)
		if b.validation1(payment) {
			payment.State = "Waiting_for_validation2"
			err := b.storage.UpdatePayment(payment)
			if err != nil {
				log.Error().Err(err)
			}
		} else {
			payment.State = "Faild"
			b.storage.UpdatePayment(payment)
		}

	case "waiting_for_validation2":
		log.Info().Msgf("Start 'Waiting_for_validation2' with payment %d", payment.ID)
		if now.Sub(payment.CreatedAt) > 60*time.Minute {
			b.storage.DeletePayment(payment.ID)
		} else if payment.Attempts >= 3 {
			payment.State = "Faild"
			b.storage.UpdatePayment(payment)
			b.notify.SendNotification(payment)
		} else if b.validation2(payment) {
			payment.State = "ready_for_closure"
			b.storage.UpdatePayment(payment)
		} else {
			payment.Attempts += 1
			b.storage.UpdatePayment(payment)
		}
	}
}

func (b *BackgroundService) validation1(payment storage.Payment) bool {

	//items len validation
	if len(payment.Items) == 0 {
		log.Error().Msg("Error. Validation1. There are not items")
		return false
	}

	//item validation
	var calculatedTotal float32
	for _, item := range payment.Items {
		calculatedTotal += item.Price * float32(item.Quantity)
	}

	if payment.Amount != calculatedTotal {
		log.Error().Msgf("Error. Validation1. payment amount %.2f not equal items sum %.2f\n", payment.Amount, calculatedTotal)
		return false
	}

	return true
}

func (b *BackgroundService) validation2(payment storage.Payment) bool {
	if payment.Amount >= 88005553535 {
		log.Error().Msg("Error. Validation2")
		return false
	}

	return true
}
