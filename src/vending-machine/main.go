package main

import (
	"fmt"
	"vending-machine/money"
	"vending-machine/payment"
	"vending-machine/product"
)

func main() {
	//loop until user want to exit
	for {
		//list of product's stock
		product.ListAllProducts()

		//list of money's stock
		money.ListAvailableMoney()

		//user select product
		buyedProducts, totalAmount, err := product.SelectProduct(nil)
		if err != nil {
			fmt.Print("error: ", err)
			return
		}

		//do payment process
		recievedMoney, changeList, isSuccessful, err := payment.Payment(totalAmount, buyedProducts)
		if err != nil {
			fmt.Print("error: ", err)
			return
		}

		//purchase summary
		payment.Summary(buyedProducts, totalAmount, recievedMoney, changeList, isSuccessful)

		//if user type "exit" then program'll terminate
		//if user ENTER then user can shop again
		var userContinue string
		fmt.Println("\nPress ENTER key to continue shopping or type \"exit\" to exit program")
		fmt.Scanln(&userContinue)
		if userContinue == "exit" {
			break
		}
	}
}
