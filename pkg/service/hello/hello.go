package hello

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-app-template/internal/config"
	"net/http"

	"github.com/go-kit/kit/endpoint"
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

	HelloService struct{}
)

var (
	ErrEmpty = errors.New("Empty string")
)

// 这里可以引入各种依赖
func NewHelloService(fileConfig *config.FileConfig) *HelloService {
	fmt.Println(fileConfig)
	return &HelloService{}
}

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
		return nil, err
	}

	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
