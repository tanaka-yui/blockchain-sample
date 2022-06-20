package transactionrepository

import (
	"blockchain/internal/domain"
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
		CreateTransaction(gateway string, input *CreatePutRequest) error
		PutTransaction(gateway string, input *CreatePutRequest) error
		ClearTransaction(gateway string) error
		GetAmount(gateWay string, input *GetAmountInput) (float32, error)
		GetChain(gateWay string) ([]*domain.Block, error)
		Consensus(gateWay string) error
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

func (r repositoryImpl) CreateTransaction(gateway string, input *CreatePutRequest) error {
	var requestBody bytes.Buffer
	if err := json.NewEncoder(&requestBody).Encode(input); err != nil {
		return errors.Wrap(err, "encode body")
	}

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/blockchain/transactions", gateway), &requestBody)
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

func (r repositoryImpl) PutTransaction(gateway string, input *CreatePutRequest) error {
	var requestBody bytes.Buffer
	if err := json.NewEncoder(&requestBody).Encode(input); err != nil {
		return errors.Wrap(err, "encode body")
	}

	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/blockchain/transactions", gateway), &requestBody)
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

func (r repositoryImpl) ClearTransaction(gateway string) error {
	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/blockchain/transactions", gateway), http.NoBody)
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

func (r repositoryImpl) GetAmount(gateWay string, input *GetAmountInput) (float32, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/blockchain/amount", gateWay), http.NoBody)
	if err != nil {
		return 0.0, err
	}
	params := request.URL.Query()
	params.Add("blockchainAddress", input.BlockchainAddress)

	request.URL.RawQuery = params.Encode()
	request.Header.Set("Content-Type", network.ApplicationJson)
	request.Close = true
	res, err := r.httpClient.Do(request.WithContext(context.Background()))
	if err != nil {
		return 0.0, err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	if res.StatusCode != http.StatusOK {
		return 0.0, fmt.Errorf("server returned a non-200 status code: %v", res.StatusCode)
	}
	var amountResponse AmountResponse
	if err := json.NewDecoder(res.Body).Decode(&amountResponse); err != nil {
		return 0.0, err
	}

	return amountResponse.Amount, nil
}

func (r repositoryImpl) GetChain(gateWay string) ([]*domain.Block, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/blockchain/chains", gateWay), http.NoBody)
	if err != nil {
		return []*domain.Block{}, err
	}
	request.Header.Set("Content-Type", network.ApplicationJson)
	request.Close = true
	res, err := r.httpClient.Do(request.WithContext(context.Background()))
	if err != nil {
		return []*domain.Block{}, err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	if res.StatusCode != http.StatusOK {
		return []*domain.Block{}, fmt.Errorf("server returned a non-200 status code: %v", res.StatusCode)
	}
	var blockRes GetChainResponse
	if err := json.NewDecoder(res.Body).Decode(&blockRes); err != nil {
		return []*domain.Block{}, err
	}

	return blockRes.Chain, nil
}

func (r repositoryImpl) Consensus(gateWay string) error {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/blockchain/consensus", gateWay), http.NoBody)
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
