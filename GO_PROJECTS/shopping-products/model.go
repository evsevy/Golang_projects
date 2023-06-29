package main

//CATEGORY DEFINITION

type Category int

const(
	Technology Category = iota
	Accessory
	Clothes
)
type Categories map[Category]string

func(c Categories) Init(){

	c[Technology] = "Technology"
	c[Accessory] = "Accessory"
	c[Clothes] = "Clothes"
}

var categories = Categories{}

// PRODUCT DEFINITION

type Product struct {
	ProductId          int         `json:"product_id"`
	ProductImageSource string 	   `json:"product_image_source"`
	ProductName        string      `json:"product_name"`
	ProductPrice       string      `json:"product_price"`
	ProductCategory    string      `json:"product_category"`
	ProductColor       string      `json:"product_color"`
}

type Products map[int]Product

func(p Products) Init(){
	if categories == nil {
		categories.Init()
	}
	p[0] = Product{
		ProductId: 0,
		ProductColor:"Gray",
		ProductName: "ASUS MG- Laptop",
		ProductPrice: "10500",
		ProductCategory: categories[Technology],
		ProductImageSource: "img/asus.jpeg",
	}
	p[1] = Product{
		ProductId: 1,
		ProductColor:"Turquoise",
		ProductName: "iPhone 11",
		ProductPrice: "10500",
		ProductCategory: categories[Technology],
		ProductImageSource: "img/iphone.jpeg",
	}
	p[2] = Product{
		ProductId: 2,
		ProductColor:"Beige",
		ProductName: "Smart Bag",
		ProductPrice: "500",
		ProductCategory: categories[Accessory],
		ProductImageSource: "img/iphone.jpeg",
	}
	p[3] = Product{
		ProductId: 3,
		ProductColor:"Black",
		ProductName: "Sony LCD-TV",
		ProductPrice: "15000",
		ProductCategory: categories[Technology],
		ProductImageSource: "img/sony.jpeg",
	}
	p[4] = Product{
		ProductId: 4,
		ProductColor:"Black",
		ProductName: "T-Shirt",
		ProductPrice: "50",
		ProductCategory: categories[Clothes],
		ProductImageSource: "img/tshirt.jpeg",
	}
}
func (p Products) TransformToSlice() []Product{
	var productSlice = []Product{}
	for _, product := range p {
		productSlice = append(productSlice,product)
	}
	return productSlice
}
func(p Products) Clear(){
	//p = map[int]Product{}
	for k := range p {
		delete(p, k)
	}
}
var products = Products{}

// BASKET DEFINITION

type BasketProduct struct{
	Product Product `json:"product"`
	Count int `json:"count"`
}
type BasketProducts map[int]BasketProduct

func(bp BasketProducts) Clear(){
	//p = map[int]Product{}
	for k := range bp {
		delete(bp, k)
	}
}
func(bp BasketProducts) Copy() map[int]BasketProduct {
	var copiedBasketProduct = BasketProducts{}
	for index, element := range bp {
		copiedBasketProduct[index] = element
	}
	return copiedBasketProduct
}
func(bp BasketProducts) TransformToSlice() []BasketProduct{
	var basketProductSlice = []BasketProduct{}
	for _,basketProduct := range bp {
		basketProductSlice = append(basketProductSlice,basketProduct)
	}
	return basketProductSlice
}
var basketProducts = BasketProducts{}

type BasketCount struct {

	Value int `json:"count"`
}

var basketCount = BasketCount{}


func InitializeModels(){
	categories.Init()
	products.Init()
	basketCount.Value = 0
}



