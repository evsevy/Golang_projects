package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)
//GET ALL PRODUCTS *
//GET PRODUCT BY ID *
//GET BASKET PRODUCTS *
//ADD PRODUCT TO BASKET  t
//INCREASE COUNT OF A BASKET PRODUCT
//DECREASE COUNT OF A BASKET PRODUCT
//REMOVE BASKET PRODUCT FROM BASKET del
//CLEAR ALL BASKET del
//GET BASKET COUNT
type Test struct{
	Description string
	Route string
	ExpectedCode int
	Method string `default:"GET"`
	Body	io.Reader
}

func TestIGetAllProducts(t *testing.T){
	categories.Init()
	products.Init()

	test := Test{
		Description: "get http status 200, when successfully get all products",
		Route: "/products",
		ExpectedCode: http.StatusOK,
		Method: http.MethodGet,
	}
	app:= fiber.New()
	app.Use(cors.New())

	app.Get("/products",func(c *fiber.Ctx) error{
		return c.JSON(products)
	})

	req := httptest.NewRequest(test.Method,test.Route, nil)
	resp,_ := app.Test(req,1)

	currentProductsAsByte,_ := json.Marshal(products)
	sendProductsAsByte,_ := ioutil.ReadAll(resp.Body)
	assert.Equalf(t, test.ExpectedCode,resp.StatusCode,test.Description)
	assert.JSONEq(t,string(currentProductsAsByte),string(sendProductsAsByte),"send products and current products are same" )
}

func TestIGetProductByProductID(t *testing.T){
	categories.Init()
	products.Init()
	productId := 1
	test := Test{
		Description: "get http status 200, when successfully get product by id",
		Route: fmt.Sprintf("/products/%d",productId),
		ExpectedCode: http.StatusOK,
		Method: http.MethodGet,
	}
	app:= fiber.New()
	app.Use(cors.New())

	app.Get("/products/:id",func(c *fiber.Ctx) error{
		productId,err := c.ParamsInt("id")
		if err != nil {
			log.Println("Product Id Taking With ParamsInt failed")
		}
		product := products[productId]
		return c.JSON(product)
	})

	req := httptest.NewRequest(test.Method,test.Route, nil)
	resp,_ := app.Test(req,1)

	currentProductAsByte,_ := json.Marshal(products[productId])
	sentProductAsByte,_ := ioutil.ReadAll(resp.Body)
	assert.Equalf(t, test.ExpectedCode,resp.StatusCode,test.Description)
	assert.JSONEq(t,string(currentProductAsByte),string(sentProductAsByte),"sent product and current product are same" )
}

func TestIGetBasketProducts(t *testing.T) {
	categories.Init()
	products.Init()
	test := Test{
		Description: "get http status 200, when successfully get products in basket",
		Route: "/basket",
		ExpectedCode: http.StatusOK,
		Method: http.MethodGet,
	}
	app:= fiber.New()
	app.Use(cors.New())
	//adding an item to basket
	productId := 1
	basketProduct := BasketProduct{Product: products[productId], Count: 1}
	basketProducts[productId] = basketProduct

	app.Get("/basket",func(c *fiber.Ctx) error{
		return c.JSON(basketProducts)
	})

	req := httptest.NewRequest(test.Method,test.Route, nil)
	resp,_ := app.Test(req,1)

	currentBasketProductsAsByte,_ := json.Marshal(basketProducts)
	sentBasketProductsAsByte,_ := ioutil.ReadAll(resp.Body)
	assert.Equalf(t, test.ExpectedCode,resp.StatusCode,test.Description)
	assert.JSONEq(t,string(currentBasketProductsAsByte),string(sentBasketProductsAsByte),"sent product and current product are same" )
}

func TestIAddedProductToBasket(t *testing.T){

	InitializeModels()
	test := Test{
		Description: "get http status 200, when successfully add new product to the basket",
		Route: "/basket",
		ExpectedCode: http.StatusOK,
		Method: http.MethodPost,
	}
	app:= fiber.New()
	app.Use(cors.New())
	//adding an item to bask

	app.Post("/basket",func(c *fiber.Ctx) error{
		var product Product
		if err := c.BodyParser(&product); err != nil {
			log.Println("body parser product parsing error")
		}
		basketProduct := BasketProduct{
			Product: product,
			Count:   1,
		}
		basketProducts[product.ProductId] = basketProduct
		basketCount.Value+=1
		return c.JSON(basketCount)
	})
	productToAddBasketAsByte,_ := json.Marshal(products[1])
	bytes.NewReader(productToAddBasketAsByte)
	req := httptest.NewRequest(test.Method,test.Route, strings.NewReader(string(productToAddBasketAsByte)))
	req.Header.Set("Content-Type","application/json")
	resp,_ := app.Test(req,1)

	responseCountAsByte,_ := ioutil.ReadAll(resp.Body)
	countAsByte,_ := json.Marshal(basketCount)


	assert.Equalf(t, test.ExpectedCode,resp.StatusCode,test.Description)
	assert.JSONEq(t,string(countAsByte),string(responseCountAsByte),"sent product and current product are same" )
	assert.Equalf(t, products[1],basketProducts[1].Product,"Added successfully to the basket products")
}

func TestIIncreasedCountOfProductInBasketByProductID(t *testing.T) {
	InitializeModels()
	testProductID := 1
	test := Test{
		Description: "get http status 200, when successfully increase count of an item in basket by product id",
		Route: fmt.Sprintf("/basket/%d/increase",testProductID),
		ExpectedCode: http.StatusOK,
		Method: http.MethodGet,
	}
	app:= fiber.New()
	app.Use(cors.New())

	app.Get("/basket/:id/increase",func(c *fiber.Ctx) error{
		productId,err := c.ParamsInt("id")
		if err!= nil {
			log.Println("Failed to get ParamsInt as productId ")
		}
		increasedBasketProduct := basketProducts[productId]
		if increasedBasketProduct == (BasketProduct{}){//if it is null
			log.Println("There is no basket product has the product Id")
		}
		increasedBasketProduct.Count += 1
		basketProducts[productId] = increasedBasketProduct
		basketCount.Value+=1

		return c.JSON(increasedBasketProduct)
	})
	//Before send the request, we need to add the product has testProductID to the basketProducts
	testBasketProduct := BasketProduct{
		Product: products[testProductID],
		Count:   1,
	}
	basketProducts[testProductID] = testBasketProduct
	req := httptest.NewRequest(test.Method,test.Route,nil)
	resp,_ := app.Test(req,1)


	assert.Equalf(t, test.ExpectedCode,resp.StatusCode,test.Description)
	assert.GreaterOrEqual(t, 2,basketProducts[testProductID].Count,"Count of basket product that has the test product id increased successfully")
}

func TestIDecreasedCountOfProductInBasketByProductID(t *testing.T){
	InitializeModels()
	testProductID := 1
	test := Test{
		Description: "get http status 200, when successfully increase count of an item in basket by product id",
		Route: fmt.Sprintf("/basket/%d/decrease",testProductID),
		ExpectedCode: http.StatusOK,
		Method: http.MethodGet,
	}
	app:= fiber.New()
	app.Use(cors.New())

	app.Get("/basket/:id/decrease",func(c *fiber.Ctx) error{
		productId,err := c.ParamsInt("id")
		if err!= nil {
			log.Println("Failed to get ParamsInt as productId ")
		}
		decreasedBasketProduct := basketProducts[productId]
		if decreasedBasketProduct == (BasketProduct{}){//if it is null
			log.Println("There is no basket product has the product Id")
		}
		decreasedBasketProduct.Count -= 1
		basketProducts[productId] = decreasedBasketProduct
		basketCount.Value-=1
		return c.JSON(decreasedBasketProduct)
	})
	//Before send the request, we need to add the product has testProductID to the basketProducts
	testBasketProduct := BasketProduct{
		Product: products[testProductID],
		Count:   2,
	}
	basketProducts[testProductID] = testBasketProduct
	req := httptest.NewRequest(test.Method,test.Route,nil)
	resp,_ := app.Test(req,1)


	assert.Equalf(t, test.ExpectedCode,resp.StatusCode,test.Description)
	assert.LessOrEqual(t, 1,basketProducts[testProductID].Count,"Count of basket product that has the test product id increased successfully")
}

func TestIRemoveProductFromBasketByProductID(t *testing.T){
	InitializeModels()
	testProductID := 1
	test := Test{
		Description: "get http status 200, when successfully increase count of an item in basket by product id",
		Route: fmt.Sprintf("/basket/%d/remove",testProductID),
		ExpectedCode: http.StatusOK,
		Method: http.MethodDelete,
	}
	app:= fiber.New()
	app.Use(cors.New())

	app.Delete("/basket/:id/remove",func(c *fiber.Ctx) error{
		productId,err := c.ParamsInt("id")
		if err!= nil {
			log.Println("Failed to get ParamsInt as productId ")
		}
		basketProductToRemove := basketProducts[productId]
		if basketProductToRemove == (BasketProduct{}){//if it is null
			log.Println("There is no basket product has the product Id")
		}
		delete(basketProducts,productId)

		return c.JSON(basketProductToRemove)
	})
	//Before send the request, we need to add the product has testProductID to the basketProducts
	testBasketProduct := BasketProduct{
		Product: products[testProductID],
		Count:   1,
	}
	basketProducts[testProductID] = testBasketProduct

	req := httptest.NewRequest(test.Method,test.Route,nil)
	resp,_ := app.Test(req,1)

	sentBasketProductAsByte,_ := ioutil.ReadAll(resp.Body)
	testBasketProductAsByte,_ := json.Marshal(testBasketProduct)

	assert.Equalf(t, test.ExpectedCode,resp.StatusCode,test.Description)
	assert.JSONEq(t,string(testBasketProductAsByte),string(sentBasketProductAsByte),"Sent Basket Product as response in remove operation not same with deleted one" )
	assert.Emptyf(t, basketProducts[testProductID],"Not Deleted Successfully")
}

func TestICleanedAllProductsInBasket(t *testing.T){
	InitializeModels()
	test := Test{
		Description: "get http status 200, when successfully clear all basket",
		Route: fmt.Sprintf("/basket"),
		ExpectedCode: http.StatusOK,
		Method: http.MethodDelete,
	}
	app:= fiber.New()
	app.Use(cors.New())

	app.Delete("/basket",func(c *fiber.Ctx) error{
		basketProductsToSendBeforeDelete := basketProducts.Copy()
		basketProducts.Clear()

		return c.JSON(basketProductsToSendBeforeDelete)
	})
	// Add to item before request to control is basket products get cleaned with request
	testBasketProduct := BasketProduct{
		Product: products[1],
		Count:   1,
	}
	basketProducts[1] = testBasketProduct

	testBasketProduct2 := BasketProduct{
		Product: products[2],
		Count:   1,
	}
	basketProducts[2] = testBasketProduct2

	basketProductsAsByte,_ := json.Marshal(basketProducts)

	req := httptest.NewRequest(test.Method,test.Route,nil)
	resp,_ := app.Test(req,1)

	sentBasketProductsAsByte,_ := ioutil.ReadAll(resp.Body)

	lengthAfterClearRequest := len(basketProducts)

	assert.Equalf(t, test.ExpectedCode,resp.StatusCode,test.Description)
	assert.Equalf(t, 0,lengthAfterClearRequest,"clear operation is done if the basket products length is zero")
	assert.JSONEq(t,string(basketProductsAsByte),string(sentBasketProductsAsByte),"sent basket products and the current basket products before clear them are not equal")
}

func TestIGetProductCountOnBasket(t *testing.T){
	InitializeModels()
	test := Test{
		Description: "get http status 200, when successfully get basket count",
		Route: fmt.Sprintf("/basket/count"),
		ExpectedCode: http.StatusOK,
		Method: http.MethodGet,
	}
	app:= fiber.New()
	app.Use(cors.New())

	app.Get("/basket/count",func(c *fiber.Ctx) error{
		return c.JSON(basketCount)
	})

	testBasketProduct := BasketProduct{
		Product: products[1],
		Count:   1,
	}
	basketProducts[1] = testBasketProduct

	testBasketProduct2 := BasketProduct{
		Product: products[2],
		Count:   1,
	}
	basketProducts[2] = testBasketProduct2


	req := httptest.NewRequest(test.Method,test.Route,nil)
	resp,_ := app.Test(req,1)
	sentBasketCountAsByte,_ := ioutil.ReadAll(resp.Body)

	currentBasketCountAsByte,_ := json.Marshal(basketCount)

	assert.Equalf(t, test.ExpectedCode,resp.StatusCode,test.Description)
	assert.JSONEq(t, string(currentBasketCountAsByte),string(sentBasketCountAsByte),"sent basket count value and current are same if the request operation succeed")

}