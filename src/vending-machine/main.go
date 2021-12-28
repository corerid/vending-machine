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
		boughtProducts, totalAmount, err := product.SelectProduct(nil)
		if err != nil {
			fmt.Print("error: ", err)
			return
		}

		//do payment process
		receivedMoney, changeList, isSuccessful, err := payment.Payment(totalAmount, boughtProducts)
		if err != nil {
			fmt.Print("error: ", err)
			return
		}

		//purchase summary
		payment.Summary(boughtProducts, totalAmount, receivedMoney, changeList, isSuccessful)

		//if user type "exit" then program will terminate
		//if user ENTER then user can shop again
		var userContinue string
		fmt.Println("\nPress ENTER key to continue shopping or type \"exit\" to exit program")
		fmt.Scanln(&userContinue)
		if userContinue == "exit" {
			break
		}
	}
}
