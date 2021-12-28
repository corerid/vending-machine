package payment

import (
	"errors"
	"fmt"
	"os"
	"vending-machine/money"
	"vending-machine/product"
)

//Payment - payment process that
//1. receive payment from user
//2. change
//3. restock of product and money
func Payment(totalProductAmount int64, buyedProducts map[product.Product]int8, userInputList ...*os.File) (map[money.Money]int8, []money.Money, bool, error) {

	//if userInput is not from file (for test purpose) then use from stdin instead
	var userInputPayment, userInputContinue *os.File
	if userInputList == nil {
		userInputPayment, userInputContinue = os.Stdin, os.Stdin
	} else {
		userInputPayment = userInputList[0]
		userInputContinue = userInputList[1]
	}

	var (
		recievedMoney = map[money.Money]int8{}
		totalPayment  int64
		changeList    []money.Money
		isSuccessful  = true
		err           error
	)

	fmt.Println("------------ Checkout ------------")
	//print product details bought by the customer
	product.PrintBuyedProdcut(buyedProducts)

	for {
		//receive payment from user
		totalPayment, recievedMoney = recievePayment(totalProductAmount, userInputPayment)

		//change the remaining money to the user
		changeAmount := totalPayment - totalProductAmount
		changeList, err = change(changeAmount, money.MoneyStock, recievedMoney)
		if err != nil {
			fmt.Printf("%+v, press ENTER key to checkout again or type \"exit\" to cancel\n", err)

			var userContinueCheckout string
			fmt.Fscanln(userInputContinue, &userContinueCheckout)

			if userContinueCheckout == "exit" {
				isSuccessful = false
				return recievedMoney, []money.Money{}, isSuccessful, nil
			}
		} else {
			break
		}
	}

	//restock
	product.DecreaseStock(buyedProducts)
	money.IncreaseStock(recievedMoney)
	money.DecreaseStock(changeList)

	return recievedMoney, changeList, isSuccessful, nil
}

func recievePayment(totalProductAmount int64, userInputList ...*os.File) (int64, map[money.Money]int8) {

	//if userInput is not from file (for test purpose) then use from stdin instead
	var userInput *os.File
	if userInputList == nil {
		userInput = os.Stdin
	} else {
		userInput = userInputList[0]
	}

	var paymentAmount int64
	recievedMoney := make(map[money.Money]int8)

	//loop until user pay more than total product's amount
	for paymentAmount < totalProductAmount {
		fmt.Println("\nTotal amount left: ", totalProductAmount-paymentAmount)
		fmt.Printf("Please select money to insert (1, 5, 10): ")

		var selectedCoin string
		fmt.Fscanln(userInput, &selectedCoin)

		//validate money that user insert
		money, err := money.CheckMoney(selectedCoin)
		if err != nil {
			fmt.Printf("%+v, please try again\n\n", err)
			continue
		}

		//map recievedMoney for count the same money that user insert
		moneyMap, ok := recievedMoney[money]
		if !ok {
			recievedMoney[money] = 1
		} else {
			recievedMoney[money] = moneyMap + 1
		}

		//accumate the payment amount from user
		paymentAmount = paymentAmount + money.Value
	}

	return paymentAmount, recievedMoney
}

//change - change the remaining money to the user
func change(changeAmount int64, availableMoney []money.Money, recievedMoney map[money.Money]int8) ([]money.Money, error) {

	var changeList []money.Money
	changeMoney := money.Money{}

	//create tmpAvailableMoney from receiving money's stock because we'll change the real stock when everything is success
	tmpAvailableMoney := make([]money.Money, len(availableMoney))
	copy(tmpAvailableMoney, availableMoney)

	//add tmp money's stock from receiving money from user
	for i, tmpAvailMoney := range tmpAvailableMoney {
		for recMoney, amount := range recievedMoney {
			if recMoney.Name == tmpAvailMoney.Name {
				tmpAvailableMoney[i].Stock = tmpAvailMoney.Stock + int64(amount)
			}
		}
	}

	//no change
	if changeAmount == 0 {
		return []money.Money{}, nil
	}

	for i, availMoney := range tmpAvailableMoney {
		//if money's value enough for change and money's stock is greater than zero
		if changeAmount-availMoney.Value >= 0 && availMoney.Stock > 0 {
			//decrese remaing amount for change
			changeAmount = changeAmount - availMoney.Value

			//decrease tmp money's stock
			tmpAvailableMoney[i].Stock = availMoney.Stock - 1

			changeMoney = money.Money{
				MoneyType: availMoney.MoneyType,
				Name:      availMoney.Name,
				Value:     availMoney.Value,
			}
			break
		}
	}

	//if there is no change that meet criteria
	if changeMoney == (money.Money{}) {
		return []money.Money{}, errors.New("insufficient change")
	}

	//there is remaing amount for change the recursive
	if changeAmount != 0 {
		otherChange, err := change(changeAmount, tmpAvailableMoney, map[money.Money]int8{})
		if err != nil {
			return []money.Money{}, err
		}
		changeList = append(otherChange, changeMoney)
		return changeList, nil

	}

	return []money.Money{changeMoney}, nil
}

func Summary(buyedProducts map[product.Product]int8, totalAmount int64, receiveMoney map[money.Money]int8, changeList []money.Money, isSuccessful bool) {
	fmt.Println("------------ Summary ------------")
	//Product details bought by the customer
	product.PrintBuyedProdcut(buyedProducts)
	fmt.Println("total price: ", totalAmount)

	if isSuccessful {
		//User payment detail
		fmt.Println("\nYou've paid")
		for money, value := range receiveMoney {
			fmt.Printf("%+v %+v for %+v %+v", money.MoneyType, money.Name, value, money.MoneyType)
			if value > 1 {
				fmt.Println("s")
			} else {
				fmt.Println("")
			}
		}

		//change detail
		fmt.Println("\nChange")
		if len(changeList) == 0 {
			fmt.Println("no change")
		} else {
			changeMap := make(map[money.Money]int8)
			for _, change := range changeList {
				val, ok := changeMap[change]
				if !ok {
					changeMap[change] = 1
				} else {
					changeMap[change] = val + 1
				}
			}
			for money, amount := range changeMap {
				fmt.Printf("%+v %+v for %+v %+v", money.MoneyType, money.Name, amount, money.MoneyType)
				if amount > 1 {
					fmt.Println("s")
				} else {
					fmt.Println("")
				}
			}
		}
	} else {
		fmt.Println("unsuccessful!")
		fmt.Println("\nreturn")
		for money, value := range receiveMoney {
			fmt.Printf("%+v %+v for %+v %+v", money.MoneyType, money.Name, value, money.MoneyType)
			if value > 1 {
				fmt.Println("s")
			} else {
				fmt.Println("")
			}
		}
	}
	fmt.Println("---------------------------------")
}
