package dto

import (
	"time"
)

//go:generate easyjson -all .

//easyjson:json

// Запрос на списание средств (прием на сервер)
// POST /api/user/balance/withdraw
type Withdrawn struct {
	Order string  `json:"order" db:"order"`
	Sum   float64 `json:"sum" db:"sum"`
	// на входе его нет, только на выходе
	ProcessedAt time.Time `json:"processed_at" db:"processed_at"`
}

// Аутентификация пользователя (прием на сервер)
// POST /api/user/login.
type UserInput struct {
	Login    string `json:"login" db:"login"`
	Password string `json:"password" db:"password"`
}

type User struct {
	Login string `json:"login" db:"username"`
	Hash  string `json:"hash" db:"hash"`
}

// Взаимодействие с системой расчёта начислений баллов лояльности
// Для взаимодействия с системой доступен один хендлер:
// GET /api/orders/{number} — получение информации о расчёте начислений баллов лояльности.
type OrderDesc struct {
	Number     string    `json:"number" db:"number"`
	Status     string    `json:"status" db:"status"`
	Accrual    float64   `json:"accural,omitempty" db:"accural"`
	UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`
}

// Получение текущего баланса пользователя (отправка с сервера)
// GET /api/user/balance
type Amount struct {
	Current   float64 `json:"current" db:"balance"`
	Withdrawn float64 `json:"withdrawn" db:"withdrawn"`
}

// Получение информации о выводе средств
// GET /api/user/withdrawals.
type Withdrawns []Withdrawn

// структура на возврат заказов (отправка с сервера)
// GET /api/user/orders
type OrdersDesc []OrderDesc

//easyjson:json
type UserArray []User
