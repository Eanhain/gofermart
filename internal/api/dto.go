package dto

import (
	"time"
)

//go:generate easyjson -all .

//easyjson:json

// Запрос на списание средств (прием на сервер)
// POST /api/user/balance/withdraw
type Withdrawn struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
	// на входе его нет, только на выходе
	Processed_at time.Time `json:"processed_at"`
}

// Аутентификация пользователя (прием на сервер)
// POST /api/user/login.
type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Hash     string
	Orders   string
	Amount   Amount
}

// структура на возврат заказов (отправка с сервера)
// GET /api/user/orders
type OrdersDesc []OrderDesc

// Взаимодействие с системой расчёта начислений баллов лояльности
// Для взаимодействия с системой доступен один хендлер:
// GET /api/orders/{number} — получение информации о расчёте начислений баллов лояльности.
type OrderDesc struct {
	Number      string    `json:"number"`
	Status      string    `json:"status"`
	Accrual     float64   `json:"accural,omitempty"`
	Uploaded_at time.Time `json:"uploaded_at"`
}

// Получение текущего баланса пользователя (отправка с сервера)
// GET /api/user/balance
type Amount struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

// Получение информации о выводе средств
// GET /api/user/withdrawals.
type Withdrawns []Withdrawn

//easyjson:json
type UserArray []User

// func ToDTO(username string, password string) User {
// 	return User{
// 		Login: username,
// 		Hash:  password}
// }
