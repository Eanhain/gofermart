package accrual

import (
	"context"
	"encoding/json"
	"fmt"

	dto "github.com/Eanhain/gofermart/internal/api"
	"github.com/Eanhain/gofermart/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type AgentAPI struct {
	*fiber.Agent
	log domain.Logger
}

func InitialAccrualApi(ctx context.Context, accrualURL string, log domain.Logger) (*AgentAPI, error) {
	agent := fiber.AcquireAgent()
	req := agent.Request()
	req.Header.SetMethod("GET")
	req.SetRequestURI(accrualURL)
	if err := agent.Parse(); err != nil {
		log.Warnln("can't init accrual api")
		return nil, err
	}
	return &AgentAPI{agent, log}, nil
}

func (a AgentAPI) GetOrder(order string) (dto.OrderDesc, error) {
	var orderDesc dto.OrderDesc
	_, body, errs := a.Bytes()
	if len(errs) > 0 {
		return orderDesc, fmt.Errorf("%w: %w", domain.ErrGetAccrualOrders, errs)
	}
	if err := json.Unmarshal(body, &orderDesc); err != nil {
		return orderDesc, fmt.Errorf("%w: %w", domain.ErrUnmarshalAccrualOrders, err)
	}
	return orderDesc, nil
}
