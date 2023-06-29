package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)


func CreateShoppingCardServerWithGivenHandlers(handlers ShoppingCardsHandlers) * fiber.App{
	var app = fiber.New()
	app.Use(cors.New())
	for _ , handler := range handlers {
		app.Add(handler.Method,handler.Route,handler.Function)
	}
	return app

}


func CreateAndStartServer(port int) error{
	InitializeModels()
	InitializeHandlers()
	app := CreateShoppingCardServerWithGivenHandlers(shoppingCardsHandlers)
	app.Static("/","./public")
	err := app.Listen(fmt.Sprintf(":%d",port))
	app.Shutdown()
	return err
}