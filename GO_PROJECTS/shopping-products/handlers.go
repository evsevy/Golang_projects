package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
)


// Coul be better ?
// HandlerType is an enum to represent a specific handler type instead of using string names to reach.
type HandlerType int
// constants of type HandlerType
const(
	GET_ALL_PRODUCTS HandlerType = iota
	GET_PRODUCT_BY_ID
	GET_BASKET_PRODUCTS
	ADD_PRODUCT_TO_BASKET
	INCREASE_COUNT_OF_A_BASKET_PRODUCT
	DECREASE_COUNT_OF_A_BASKET_PRODUCT
	REMOVE_BASKET_PRODUCT_FROM_BASKET
	CLEAR_ALL_BASKET
	GET_BASKET_COUNT
)

// HandlerFunction represents a fiber handler
type HandlerFunction func(c *fiber.Ctx) error

// Handler represents a handle operation by keeping the Handling Route and Handling Function as pair.
type Handler struct{
	Route string
	Function HandlerFunction
	Method  string
}

type ShoppingCardsHandlers map[HandlerType]Handler

func(handlers ShoppingCardsHandlers) Init(){
	handlers[GET_ALL_PRODUCTS] = Handler{
		Route:    "/products",
		Function: GetAllProducts,
		Method: http.MethodGet,
	}
	handlers[GET_PRODUCT_BY_ID] = Handler{
		Route:    "/products/:id",
		Function: GetProductById,
		Method: http.MethodGet,
	}
	handlers[GET_BASKET_PRODUCTS] = Handler{
		Route:    "/basket",
		Function: GetBasketProducts,
		Method: http.MethodGet,
	}
	handlers[ADD_PRODUCT_TO_BASKET] = Handler{
		Route:    "/basket",
		Function: AddProductToBasket,
		Method: http.MethodPost,
	}
	handlers[INCREASE_COUNT_OF_A_BASKET_PRODUCT] = Handler{
		Route:    "/basket/:id/increase",
		Function: IncreaseCountOfABasketProduct,
		Method: http.MethodGet,
	}
	handlers[DECREASE_COUNT_OF_A_BASKET_PRODUCT] = Handler{
		Route:    "/basket/:id/decrease",
		Function: DecreaseCountOfABasketProduct,
		Method: http.MethodGet,
	}
	handlers[REMOVE_BASKET_PRODUCT_FROM_BASKET] = Handler{
		Route:    "/basket/:id/remove",
		Function: RemoveBasketProductFromBasket,
		Method: http.MethodDelete,
	}
	handlers[CLEAR_ALL_BASKET] = Handler{
		Route:    "/basket",
		Function: ClearAllBasket,
		Method: http.MethodDelete,
	}
	handlers[GET_BASKET_COUNT] = Handler{
		Route:    "/basket/count",
		Function: GetBasketCount,
		Method: http.MethodGet,
	}
}

func GetAllProducts(c *fiber.Ctx) error{
	return c.JSON(products.TransformToSlice())
}
func GetProductById(c *fiber.Ctx) error{
	productId,err := c.ParamsInt("id")
	if err != nil {
		log.Println("Product Id Taking With ParamsInt failed")
	}
	product := products[productId]
	return c.JSON(product)
}
func GetBasketProducts(c *fiber.Ctx) error{
	return c.JSON(basketProducts.TransformToSlice())
}
func AddProductToBasket(c *fiber.Ctx) error{
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
}
func IncreaseCountOfABasketProduct(c *fiber.Ctx) error{
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
}
func DecreaseCountOfABasketProduct(c *fiber.Ctx) error{
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
}
func RemoveBasketProductFromBasket(c *fiber.Ctx) error{
	productId,err := c.ParamsInt("id")
	if err!= nil {
		log.Println("Failed to get ParamsInt as productId ")
	}
	basketProductToRemove := basketProducts[productId]
	if basketProductToRemove == (BasketProduct{}){//if it is null
		log.Println("There is no basket product has the product Id")
	}
	delete(basketProducts,productId)

	basketCount.Value -= basketProductToRemove.Count
	return c.JSON(basketProductToRemove)
}
func ClearAllBasket(c *fiber.Ctx) error{
	basketProductsToSendBeforeDelete := basketProducts.TransformToSlice()
	basketProducts.Clear()
	basketCount.Value = 0
	return c.JSON(basketProductsToSendBeforeDelete)
}
func GetBasketCount(c *fiber.Ctx) error{
	return c.JSON(basketCount)
}

var shoppingCardsHandlers = ShoppingCardsHandlers{}

func InitializeHandlers(){
	shoppingCardsHandlers.Init()
}