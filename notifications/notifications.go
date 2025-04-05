package notifications

import (
	"fmt"
	"payments_service/storage"
)

type NotificationsStruct struct {
}

func (n *NotificationsStruct) SendNotification(payment storage.Payment) {
	fmt.Printf("Платеж ID %d провален после %d попыток\n", payment.ID, payment.Attempts)
	// Здесь можно добавить отправку email или webhook
}
