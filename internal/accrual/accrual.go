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
	accrualHost string
	log         domain.Logger
}

func InitialAccrualApi(ctx context.Context, accrualURL string, log domain.Logger) (*AgentAPI, error) {
	agent := fiber.AcquireAgent()

	return &AgentAPI{agent, accrualURL, log}, nil
}

func (a AgentAPI) GetOrder(order string) (dto.OrderDesc, error) {
	var orderDesc dto.OrderDesc
	var buf []byte
	req := a.Request()
	req.Header.SetMethod("GET")
	req.SetRequestURI("http://" + a.accrualHost + "/api/orders/" + order)
	if err := a.Parse(); err != nil {
		a.log.Warnln("can't init accrual api")
		return orderDesc, err
	}
	a.Body(buf)
	// if len(errs) > 0 {
	// 	return orderDesc, fmt.Errorf("%w: %w", domain.ErrGetAccrualOrders, errs[0])
	// }
	if err := json.Unmarshal(buf, &orderDesc); err != nil {
		return orderDesc, fmt.Errorf("%w: %w, %v", domain.ErrUnmarshalAccrualOrders, err, buf)
	}

	a.log.Infoln(orderDesc)
	return orderDesc, nil
}
