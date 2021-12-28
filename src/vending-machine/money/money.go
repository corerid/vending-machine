package money

import (
	"errors"
	"fmt"
	"sort"
)

const (
	COIN = "coin"
	BANK = "bank"
)

type Money struct {
	MoneyType string
	Name      string
	Value     int64
	Stock     int64
}

var MoneyStock = []Money{
	{
		MoneyType: COIN,
		Name:      "1",
		Value:     1,
		Stock:     2,
	},
	{
		MoneyType: COIN,
		Name:      "5",
		Value:     5,
		Stock:     2,
	},
	{
		MoneyType: COIN,
		Name:      "10",
		Value:     10,
		Stock:     10,
	},
}

func init() {
	//sort descending AvailableMoney to prioritize the change from the most valuable amount to lowest
	//ex. 10 > 5 > 1
	sort.Slice(MoneyStock, func(i, j int) bool {
		return MoneyStock[i].Value > MoneyStock[j].Value
	})
}

func ListAvailableMoney() {
	fmt.Println("List of money")
	fmt.Println("MoneyType   Name        Value       Stock")
	fmt.Println("------------------------------------------")
	for _, money := range MoneyStock {
		fmt.Printf("%-12v%-12v%-12v%-12v\n", money.MoneyType, money.Name, money.Value, money.Stock)
	}
	fmt.Println("-----------------------------------")
}

//CheckMoney - for validate money is exist in stock or not
func CheckMoney(moneyName string) (Money, error) {
	for _, availMoney := range MoneyStock {
		if moneyName == availMoney.Name {
			return availMoney, nil
		}
	}
	return Money{}, errors.New("money doesn't excepted")
}

//IncreaseStock - increase global money's stock from recievedMoney map (money received from user)
func IncreaseStock(recievedMoney map[Money]int8) error {
	for recMoney, amount := range recievedMoney {
		for i, availMoney := range MoneyStock {
			if recMoney.Name == availMoney.Name {
				MoneyStock[i].Stock = MoneyStock[i].Stock + int64(amount)
				break
			}
		}
	}
	return nil
}

//IncreaseStock - decrease global money's stock from changeList (money thaะ changed to user)
func DecreaseStock(changeList []Money) error {
	for _, change := range changeList {
		for i, availMoney := range MoneyStock {
			if change.Name == availMoney.Name {
				if MoneyStock[i].Stock-1 < 0 {
					return errors.New(availMoney.Name + "'s stock is less than zero")
				}
				MoneyStock[i].Stock = MoneyStock[i].Stock - 1
				break
			}
		}
	}

	return nil
}
