package hello

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type (
	ServiceInterface interface {
		Hello(string) (string, error)
	}

	helloRequest struct {
		S string `json:"s"`
	}

	helloResponse struct {
		V   string `json:"v"`
		Err string `json:"err,omitempty"`
	}

	HelloService struct {
		logger *zap.Logger
	}
)

var (
	ErrEmpty = errors.New("Empty string")
)

// 这里可以引入各种依赖
func NewHelloService(log *zap.Logger) *HelloService {
	return &HelloService{log}
}

// Hello
// @Summary Say Hello
// @Description this is description
// @Tags hello
// @Accept  octet-stream
// @Produce octet-stream
// @Param   req body string false "string enums"
// @Success 200 {object} string
// @Router  /hello [get]
func (h *HelloService) Hello(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}

	return fmt.Sprintf("Hello %s", s), nil
}

func makeHelloEndpoint(svc ServiceInterface) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(helloRequest)
		v, err := svc.Hello(req.S)
		if err != nil {
			return helloResponse{
				V:   v,
				Err: err.Error(),
			}, nil
		}
		return helloResponse{v, ""}, nil
	}
}

func decodeHelloRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request helloRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, errors.Wrap(err, "decodeHelloRequest failed")
	}

	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
