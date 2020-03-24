package hello

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	helloPb "go-app-template/api"

	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type (
	ServiceInterface interface {
		Hello(request helloPb.Request) (helloPb.Response, error)
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
// @Param   req body Request true "Request"
// @Success 200 {object} Response
// @Router  /hello [get]
func (h *HelloService) Hello(req helloPb.Request) (helloPb.Response, error) {
	if req.Name == "" {
		return helloPb.Response{
			Greeting: "Hello, nobody",
		}, ErrEmpty
	}

	return helloPb.Response{
		Greeting: fmt.Sprintf("Hello %s", req.Name),
	}, nil
}

func makeHelloEndpoint(svc ServiceInterface) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(helloPb.Request)
		rsp, err := svc.Hello(req)
		if err != nil {
			return nil, err
		}
		return rsp, nil
	}
}

func decodeHelloRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request helloPb.Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, errors.Wrap(err, "decodeHelloRequest failed")
	}

	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
