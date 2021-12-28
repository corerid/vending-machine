package product

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Product struct {
	ProductNo int8
	Name      string
	Price     int64
	Stock     int8
}

var ProductStock = []Product{
	{
		ProductNo: 1,
		Name:      "Lays",
		Price:     5,
		Stock:     1,
	},
	{
		ProductNo: 2,
		Name:      "Hanami",
		Price:     10,
		Stock:     10,
	},
	{
		ProductNo: 3,
		Name:      "Kitkat",
		Price:     25,
		Stock:     10,
	},
	{
		ProductNo: 4,
		Name:      "Pepsi",
		Price:     15,
		Stock:     10,
	},
}

func ListAllProducts() {
	fmt.Println("List of products")
	fmt.Println("No        Name      Price     Stock")
	fmt.Println("-----------------------------------")
	for _, product := range ProductStock {
		fmt.Printf("%-10v%-10v%-10v%-10v\n", product.ProductNo, product.Name, product.Price, product.Stock)
	}
	fmt.Println("-----------------------------------")
}

func SelectProduct(userInput *os.File) (map[Product]int8, int64, error) {

	//if userInput is not from file (for test purpose) then use from stdin instead
	if userInput == nil {
		userInput = os.Stdin
	}

	//create tmpProductStock from product's stock because we'll change the real stock when everything is success
	tmpProductStock := make([]Product, len(ProductStock))
	copy(tmpProductStock, ProductStock)

	//loop for select product until user ENTER for checkout
	var totalAmount int64
	boughtProducts := make(map[Product]int8)
	fmt.Println("Please Select Product No: ")
	for {
		var selectedProduct string
		fmt.Fscanln(userInput, &selectedProduct)

		//if user ENTER then finish the loop for next process (checkout)
		if selectedProduct == "" {
			//if user does not select any product then error (assume that user doesn't want to continue shopping)
			if len(boughtProducts) == 0 {
				return boughtProducts, 0, errors.New("you have not select any product")
			}
			break
		}

		//validate product that user's selected
		product, err := checkProduct(selectedProduct, tmpProductStock)
		if err != nil {
			//if product validation is not pass, then let user to select product again
			fmt.Printf("%+v, please select product no. again\n", err)
			continue
		}

		//accumulate the price of the product selected by user
		totalAmount = totalAmount + product.Price

		//map boughtProducts for count the same product
		productMap, ok := boughtProducts[product]
		if !ok {
			boughtProducts[product] = 1
		} else {
			boughtProducts[product] = productMap + 1
		}

		fmt.Println("Press ENTER to checkout or continue select product")
	}

	return boughtProducts, totalAmount, nil
}

//checkProduct - for check product no. from receiving product's stock is available or not
//and return product for founded product no.
func checkProduct(productNo string, productStock []Product) (Product, error) {
	for i, product := range productStock {
		productNoInt, err := strconv.Atoi(productNo)
		if err != nil {
			return Product{}, errors.New("invalid input")
		}
		if int8(productNoInt) == product.ProductNo {
			//if product's stock is zero then error
			if product.Stock == 0 {
				return Product{}, errors.New(product.Name + " is out of stock")
			}

			//return and decrease stock of target product
			productStock[i].Stock = productStock[i].Stock - 1
			boughtProduct := Product{
				ProductNo: product.ProductNo,
				Name:      product.Name,
				Price:     product.Price,
			}
			return boughtProduct, nil
		}
	}

	return Product{}, errors.New("product doesn't exist")
}

//DecreaseStock - decrease global product's stock by buyedProducts map (products that user buy)
func DecreaseStock(buyedProducts map[Product]int8) error {
	for boughtProduct, amount := range buyedProducts {
		for i, product := range ProductStock {
			if boughtProduct.ProductNo == product.ProductNo {
				if ProductStock[i].Stock-amount < 0 {
					return errors.New(product.Name + "'s stock is less than zero")
				}
				ProductStock[i].Stock = ProductStock[i].Stock - amount
				break
			}
		}
	}
	return nil
}

func PrintBoughtProduct(boughtProducts map[Product]int8) {
	fmt.Printf("You've bought\n")
	for product, value := range boughtProducts {
		fmt.Printf("%+v price %+v for %+v piece", product.Name, product.Price, value)
		if value > 1 {
			fmt.Println("s")
		} else {
			fmt.Println("")
		}
	}
}
