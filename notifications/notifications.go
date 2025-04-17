package notifications

import (
	"fmt"
	"payments_service/models"
)

type NotificationsStruct struct {
}

func (n *NotificationsStruct) SendNotification(payment models.Payment) {
	fmt.Printf("Платеж ID %d провален после %d попыток\n", payment.ID, payment.Attempts)
	// Здесь можно добавить отправку email или webhook
}
