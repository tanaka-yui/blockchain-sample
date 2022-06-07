package transactionrepository

import (
	"blockchain/pkg/network"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

type (
	Repository interface {
		PutTransaction(endPoint string, input *PutRequest) error
	}
	repositoryImpl struct {
		httpClient *http.Client
	}
)

func NewRepository() Repository {
	return &repositoryImpl{
		httpClient: http.DefaultClient,
	}
}

func (r repositoryImpl) PutTransaction(endPoint string, input *PutRequest) error {
	var requestBody bytes.Buffer
	if err := json.NewEncoder(&requestBody).Encode(input); err != nil {
		return errors.Wrap(err, "encode body")
	}

	request, err := http.NewRequest(http.MethodPut, endPoint, &requestBody)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", network.ApplicationJson)
	request.Close = true
	res, err := r.httpClient.Do(request.WithContext(context.Background()))
	if err != nil {
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned a non-200 status code: %v", res.StatusCode)
	}
	return nil
}
