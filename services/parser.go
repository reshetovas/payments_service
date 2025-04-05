package services

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"payments_service/storage"
	"time"
)

type Item struct {
	Name     string
	Price    float32
	Quantity int
}

type PaymentData struct {
	Date         string
	Shop_ID      int
	Address      string
	Total_Amount float32
	Items        []Item
}

type PaymentFile struct {
	Payment PaymentData
}

type ParseService struct {
	storage storage.StorageActions
}

func NewParseService(storage storage.StorageActions) *ParseService {
	return &ParseService{
		storage: storage,
	}
}

func (p ParseService) ParsePaymentsFile(fileName string) error {

	fmt.Println("start parsing")

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// var payments []PaymentData

	for i := 1; scanner.Scan(); i++ {

		var raw PaymentFile

		line := scanner.Text()
		json.Unmarshal([]byte(line), &raw)

		// printStructInfo(raw.Payment)

		err := p.savePaymentToDB(raw.Payment)
		if err != nil {
			return err
		}

		//payments = append(payments, raw.Payment)
		//fmt.Println(raw)
	}
	//	fmt.Println("Payments", payments)
	return err
}

func (p ParseService) savePaymentToDB(payment PaymentData) error {

	createdAt := time.Now()

	paymentInput := storage.Payment{
		Amount:      payment.Total_Amount,
		Description: "Imported from JSON",
		CreatedAt:   createdAt,
		ShopID:      payment.Shop_ID,
		Address:     payment.Address,
	}

	// fmt.Println(paymentInput)

	// Создаём платёж в БД
	paymentID, err := p.storage.CreatePayment(paymentInput)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// fmt.Println("start creating items")
	// Добавляем товары в БД
	for _, item := range payment.Items {
		err = p.storage.CreateItem(storage.Item{
			PaymentID: paymentID,
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  item.Quantity,
		})
		if err != nil {
			return err
		}
		// fmt.Println("item round finish")
	}
	return nil
}

// func printStructInfo(p PaymentData) {
// 	v := reflect.ValueOf(p)
// 	t := v.Type()

// 	fmt.Println("Struct Payment:")
// 	for i := 0; i < v.NumField(); i++ {
// 		field := t.Field(i)
// 		value := v.Field(i)

// 		fmt.Printf("- %s (%s) = %v\n", field.Name, field.Type, value)
// 	}
// }

// package services

// import (
// 	"encoding/json"
// 	"os"
// 	"payments_service/storage"
// 	"time"
// )

// type Item struct {
// 	Name     string
// 	Price    float32
// 	Quantity int
// }

// type PaymentData struct {
// 	Date        string
// 	ShopID      int
// 	Address     string
// 	TotalAmount float32
// 	Items       []Item
// }

// // main struct for JSON-file
// type PaymentsFile struct {
// 	Payment PaymentData
// }

// func ParsePaymentsFile(filename string, storage storage.StorageActions) error {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	decoder := json.NewDecoder(file)

// 	for {
// 		var paymentFile PaymentsFile
// 		if err := decoder.Decode(&paymentFile); err != nil {
// 			break
// 		}

// 		err := savePaymentToDB(paymentFile.Payment, storage)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func savePaymentToDB(payment PaymentData, storage storage.StorageActions) error {
// 	createdAt := time.Now()

// 	paymentInput := storage.Payment{
// 		Amount:      payment.TotalAmount,
// 		Description: "Imported from JSON",
// 		CreatedAt:   createdAt,
// 		ShopID:      payment.ShopID,
// 		Address:     payment.Address,
// 	}

// 	// Создаём платёж в БД
// 	paymentID, err := storage.CreatePayment(paymentInput)
// 	if err != nil {
// 		return err
// 	}

// 	// Добавляем товары в БД
// 	for _, item := range payment.Items {
// 		err = storage.CreateItem(storage.Item{
// 			PaymentID: paymentID,
// 			Name:      item.Name,
// 			Price:     item.Price,
// 			Quantity:  item.Quantity,
// 		})
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
