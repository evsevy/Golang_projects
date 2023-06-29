package main

import (
	"fmt"
	"github.com/k0kubun/pp"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"github.com/pact-foundation/pact-go/utils"
	"log"
	"testing"
)

type StateHandlers map[string]types.StateHandler

func (s StateHandlers) Init(){
	 emptyStateHandler := func() error {//if there is nothing do in this provider state use this empty function.
		return nil
	 }
	s["i get products successfully"] = emptyStateHandler
	s["i get product successfully"] = emptyStateHandler
	s["i get basket products successfully"] = func() error { //adding a product to basket to control provider state on product id :1
		testProductId := 1
		basketProducts[testProductId] = BasketProduct{
			Count: 1,
			Product: products[testProductId],
		}
		return nil
	}
	s["i add product to basket products successfully"] = emptyStateHandler
	s["i increased count of product successfully by given product id"] = func() error { //there should be a basket product with id 1 to be able to increase it
		testProductId := 1
		basketProducts[testProductId] = BasketProduct{
			Count: 1,
			Product: products[testProductId],
		}
		return nil
	}
	s["i decreased count of product successfully by given product id"] = func() error {//there should be a basket product with id 1 to be able to increase it.
		testProductId := 1
		basketProducts[testProductId] = BasketProduct{
			Count: 2,// not control if less then 1 then remove condition. So I gave bigger count
			Product: products[testProductId],
		}
		return nil
	}
	s["i removed a product successfully from basket products by the products id"] = func() error {
		testProductId := 1
		basketProducts[testProductId] = BasketProduct{
			Count: 1,
			Product: products[testProductId],
		}
		return nil
	}
	s["i cleared all the products in basket successfully"] = func() error { // adding to product randomly to be sure the cleaning operation done in success
		firstTestProductId,secondTestProductId := 1,2
		basketProducts[firstTestProductId] = BasketProduct{
			Count: 1,
			Product: products[firstTestProductId],
		}
		basketProducts[secondTestProductId] = BasketProduct{
			Count: 1,
			Product: products[secondTestProductId],
		}
		return nil
	}
	s["i get count of products in basket sucessfully"] = emptyStateHandler
}
func (s StateHandlers) TransformToStringStateHandlerMap() map[string]types.StateHandler{
	var stateHandlers = map[string]types.StateHandler{}
	for state, handler := range s {
		stateHandlers[state] = handler
	}
	return stateHandlers
}
type Settings struct{
	Host string
	ConsumerName string
	ProviderName string
	PactURL string
	PublishVerificationResults	bool
	FailIfNoPactsFound	bool
	DisableToolValidityCheck bool
	BrokerBaseURL string
	BrokerToken string
	ProviderVersion string
	PactFileWriteMode string
	StateHandlers * StateHandlers
}

func (s * Settings) Init(){
	s.Host = "127.0.0.1"
	s.ConsumerName = "ShoppingCardClient"
	s.ProviderName = "ShoppingCardApi"
	s.PactURL = "https://eneskzlcn.pactflow.io/pacts/provider/ShoppingCardApi/consumer/ShoppingCardClient/version/a1aa11ed35c452df8177871ab88eabb32323501c"
	s.PublishVerificationResults = true
	s.FailIfNoPactsFound = true
	s.DisableToolValidityCheck = true
	s.BrokerBaseURL = "https://eneskzlcn.pactflow.io"
	s.BrokerToken = "L0IzB6WxiCRX7sEdAQoWlQ"
	s.ProviderVersion = "1.0.0"
	s.PactFileWriteMode = "merge"
	s.StateHandlers = &StateHandlers{}
	s.StateHandlers.Init()

}
func TestProvider(t *testing.T){
	port,_ := utils.GetFreePort()

	go CreateAndStartServer(port)

	settings :=	Settings{}
	settings.Init()

	pact := dsl.Pact{

		Consumer:                 settings.ConsumerName,
		Provider:                 settings.ProviderName,
		Host:                     settings.Host,
		DisableToolValidityCheck: settings.DisableToolValidityCheck,
	}

	log.Println(pact.Host)
	verifyRequest := types.VerifyRequest{
		ProviderBaseURL:           fmt.Sprintf("http://%s:%d", settings.Host, port),
		PactURLs:                   []string{settings.PactURL},
		BrokerURL:                  settings.BrokerBaseURL,
		BrokerToken:                settings.BrokerToken,
		FailIfNoPactsFound:         settings.FailIfNoPactsFound,
		PublishVerificationResults: settings.PublishVerificationResults,
		ProviderVersion:            settings.ProviderVersion,
		StateHandlers:              settings.StateHandlers.TransformToStringStateHandlerMap(),
	}

	verifyResponses, err := pact.VerifyProvider(t, verifyRequest)
	if err != nil {
		t.Fatal(err)
	}

	pp.Println(len(verifyResponses), "pact tests run")
}