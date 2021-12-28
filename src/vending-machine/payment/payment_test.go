package payment

import (
	"errors"
	"io"
	"io/ioutil"
	"sort"
	"testing"
	"vending-machine/money"
	"vending-machine/product"

	"github.com/stretchr/testify/assert"
)

func Test_Payment(t *testing.T) {
	type inputArgs struct {
		totalAmount       int64
		buyedProducts     map[product.Product]int8
		userInputPayment  string
		userInputContinue string
	}

	type expectedArgs struct {
		expectedRecievedMoney map[money.Money]int8
		expectedChangeList    []money.Money
		expectedProductStock  []product.Product
		expectedMoneyStock    []money.Money
		expectedIsSuccessful  bool
		expectedError         error
	}

	tests := []struct {
		description string
		prepData    func()
		input       inputArgs
		expected    expectedArgs
		hasError    bool
	}{
		{
			description: "test_payment_success",
			prepData: func() {
				money.AvailableMoney = []money.Money{
					{
						MoneyType: money.COIN,
						Name:      "10",
						Value:     10,
						Stock:     10,
					},
					{
						MoneyType: money.COIN,
						Name:      "5",
						Value:     5,
						Stock:     10,
					},
					{
						MoneyType: money.COIN,
						Name:      "1",
						Value:     1,
						Stock:     10,
					},
				}

				//sort descending input available money first for prioritize chage
				sort.Slice(money.AvailableMoney, func(i, j int) bool {
					return money.AvailableMoney[i].Value > money.AvailableMoney[j].Value
				})

				product.Products = []product.Product{
					{
						ProductNo: 3,
						Name:      "Kitkat",
						Price:     25,
						Stock:     10,
					},
				}
			},
			input: inputArgs{
				totalAmount: 25,
				buyedProducts: map[product.Product]int8{
					{
						ProductNo: 3,
						Name:      "Kitkat",
						Price:     25,
						Stock:     10,
					}: 1,
				},
				userInputContinue: "\n",
				userInputPayment:  "1\n1\n1\n1\n1\n5\n10\n10\n",
			},
			expected: expectedArgs{
				expectedChangeList: []money.Money{
					{
						MoneyType: money.COIN,
						Name:      "5",
						Value:     5,
					},
				},
				expectedRecievedMoney: map[money.Money]int8{
					{
						MoneyType: money.COIN,
						Name:      "1",
						Value:     1,
						Stock:     10,
					}: 5,
					{
						MoneyType: money.COIN,
						Name:      "5",
						Value:     5,
						Stock:     10,
					}: 1,
					{
						MoneyType: money.COIN,
						Name:      "10",
						Value:     10,
						Stock:     10,
					}: 2,
				},
				expectedProductStock: []product.Product{
					{
						ProductNo: 3,
						Name:      "Kitkat",
						Price:     25,
						Stock:     9,
					},
				},
				expectedMoneyStock: []money.Money{
					{
						MoneyType: money.COIN,
						Name:      "10",
						Value:     10,
						Stock:     12,
					},
					{
						MoneyType: money.COIN,
						Name:      "5",
						Value:     5,
						Stock:     10,
					},
					{
						MoneyType: money.COIN,
						Name:      "1",
						Value:     1,
						Stock:     15,
					},
				},
				expectedIsSuccessful: true,
			},
			hasError: false,
		},
		{
			description: "test_payment_success_with_insufficient_change",
			prepData: func() {
				money.AvailableMoney = []money.Money{
					{
						MoneyType: money.COIN,
						Name:      "10",
						Value:     10,
						Stock:     0,
					},
					{
						MoneyType: money.COIN,
						Name:      "5",
						Value:     5,
						Stock:     0,
					},
					{
						MoneyType: money.COIN,
						Name:      "1",
						Value:     1,
						Stock:     0,
					},
				}

				//sort descending input available money first for prioritize chage
				sort.Slice(money.AvailableMoney, func(i, j int) bool {
					return money.AvailableMoney[i].Value > money.AvailableMoney[j].Value
				})

				product.Products = []product.Product{
					{
						ProductNo: 1,
						Name:      "Lays",
						Price:     5,
						Stock:     10,
					},
				}
			},
			input: inputArgs{
				totalAmount: 5,
				buyedProducts: map[product.Product]int8{
					{
						ProductNo: 1,
						Name:      "Lays",
						Price:     5,
					}: 1,
				},
				userInputContinue: "exit\n",
				userInputPayment:  "10\n",
			},
			expected: expectedArgs{
				expectedChangeList: []money.Money{},
				expectedRecievedMoney: map[money.Money]int8{
					{
						MoneyType: money.COIN,
						Name:      "10",
						Value:     10,
						Stock:     0,
					}: 1,
				},
				expectedProductStock: []product.Product{
					{
						ProductNo: 1,
						Name:      "Lays",
						Price:     5,
						Stock:     10,
					},
				},
				expectedMoneyStock: []money.Money{
					{
						MoneyType: money.COIN,
						Name:      "10",
						Value:     10,
						Stock:     0,
					},
					{
						MoneyType: money.COIN,
						Name:      "5",
						Value:     5,
						Stock:     0,
					},
					{
						MoneyType: money.COIN,
						Name:      "1",
						Value:     1,
						Stock:     0,
					},
				},
				expectedIsSuccessful: false,
			},
			hasError: false,
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			test.prepData()

			//create mock user input
			userInputPayment, err := ioutil.TempFile("", "")
			if err != nil {
				t.Fatal(err)
			}
			defer userInputPayment.Close()

			_, err = io.WriteString(userInputPayment, test.input.userInputPayment)
			if err != nil {
				t.Fatal(err)
			}

			_, err = userInputPayment.Seek(0, io.SeekStart)
			if err != nil {
				t.Fatal(err)
			}

			userInputContinue, err := ioutil.TempFile("", "")
			if err != nil {
				t.Fatal(err)
			}
			defer userInputContinue.Close()

			_, err = io.WriteString(userInputContinue, test.input.userInputContinue)
			if err != nil {
				t.Fatal(err)
			}

			_, err = userInputContinue.Seek(0, io.SeekStart)
			if err != nil {
				t.Fatal(err)
			}

			actualRecievedMoney, actualChangeList, actualIsSuccessful, err := Payment(test.input.totalAmount, test.input.buyedProducts, userInputPayment, userInputContinue)

			if test.hasError {
				assert.Error(t, err)
				assert.Equal(t, test.expected.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.expected.expectedChangeList, actualChangeList)
			for k := range actualRecievedMoney {
				assert.Equal(t, test.expected.expectedRecievedMoney[k], actualRecievedMoney[k])
			}
			assert.Equal(t, test.expected.expectedProductStock, product.Products)
			assert.Equal(t, test.expected.expectedMoneyStock, money.AvailableMoney)
			assert.Equal(t, test.expected.expectedIsSuccessful, actualIsSuccessful)

		})
	}
}

func Test_recievePayment(t *testing.T) {
	type inputArgs struct {
		totalAmount int64
		userInput   string
	}

	type expectedArgs struct {
		expectedPaymentAmount int64
		expectedRecieveMoney  map[money.Money]int8
		expectedError         error
	}

	tests := []struct {
		description string
		prepData    func()
		input       inputArgs
		expected    expectedArgs
		hasError    bool
	}{
		{
			description: "test_recieve_payment_success",
			prepData: func() {
				money.AvailableMoney = []money.Money{
					{
						MoneyType: money.COIN,
						Name:      "1",
						Value:     1,
						Stock:     10,
					},
					{
						MoneyType: money.COIN,
						Name:      "5",
						Value:     5,
						Stock:     10,
					},
					{
						MoneyType: money.COIN,
						Name:      "10",
						Value:     10,
						Stock:     10,
					},
				}
			},
			input: inputArgs{
				totalAmount: 25,
				userInput:   "1\n1\n1\n1\n1\n5\n10\n10\n",
			},
			expected: expectedArgs{
				expectedPaymentAmount: 30,
				expectedRecieveMoney: map[money.Money]int8{
					{
						MoneyType: money.COIN,
						Name:      "1",
						Value:     1,
						Stock:     10,
					}: 5,
					{
						MoneyType: money.COIN,
						Name:      "5",
						Value:     5,
						Stock:     10,
					}: 1,
					{
						MoneyType: money.COIN,
						Name:      "10",
						Value:     10,
						Stock:     10,
					}: 2,
				},
			},
			hasError: false,
		},
		{
			description: "test_recieve_payment_success_with_money_does_not_excepted",
			prepData: func() {
				money.AvailableMoney = []money.Money{
					{
						MoneyType: money.COIN,
						Name:      "1",
						Value:     1,
						Stock:     10,
					},
					{
						MoneyType: money.COIN,
						Name:      "5",
						Value:     5,
						Stock:     10,
					},
					{
						MoneyType: money.COIN,
						Name:      "10",
						Value:     10,
						Stock:     10,
					},
				}
			},
			input: inputArgs{
				totalAmount: 25,
				userInput:   "20\n1\n1\n1\n1\n1\n5\n10\n10\n",
			},
			expected: expectedArgs{
				expectedPaymentAmount: 30,
				expectedRecieveMoney: map[money.Money]int8{
					{
						MoneyType: money.COIN,
						Name:      "1",
						Value:     1,
						Stock:     10,
					}: 5,
					{
						MoneyType: money.COIN,
						Name:      "5",
						Value:     5,
						Stock:     10,
					}: 1,
					{
						MoneyType: money.COIN,
						Name:      "10",
						Value:     10,
						Stock:     10,
					}: 2,
				},
			},
			hasError: false,
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			test.prepData()

			//create mock user input
			userInput, err := ioutil.TempFile("", "")
			if err != nil {
				t.Fatal(err)
			}
			defer userInput.Close()

			_, err = io.WriteString(userInput, test.input.userInput)
			if err != nil {
				t.Fatal(err)
			}

			_, err = userInput.Seek(0, io.SeekStart)
			if err != nil {
				t.Fatal(err)
			}

			actualPaymentAmount, actualRecieveMoney := recievePayment(test.input.totalAmount, userInput)
			if test.hasError {
				assert.Error(t, err)
				assert.Equal(t, test.expected.expectedError, err)
				assert.Equal(t, test.expected.expectedPaymentAmount, actualPaymentAmount)
				assert.Equal(t, test.expected.expectedRecieveMoney, actualRecieveMoney)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected.expectedPaymentAmount, actualPaymentAmount)
				for k := range actualRecieveMoney {
					assert.Equal(t, test.expected.expectedRecieveMoney[k], actualRecieveMoney[k])
				}
			}
		})
	}
}

func Test_change(t *testing.T) {
	type inputArgs struct {
		changeAmount   int64
		availableMoney []money.Money
		recievedMoney  map[money.Money]int8
	}

	type expectedArgs struct {
		expectedMoney []money.Money
		expectedError error
	}

	tests := []struct {
		description string
		prepData    func()
		input       inputArgs
		expected    expectedArgs
		hasError    bool
	}{
		{
			description: "test_change_success_with_no_recieve_money",
			input: inputArgs{
				changeAmount: 88,
				availableMoney: []money.Money{
					{
						MoneyType: money.COIN,
						Name:      "1",
						Value:     1,
						Stock:     10,
					},
					{
						MoneyType: money.COIN,
						Name:      "5",
						Value:     5,
						Stock:     10,
					},
					{
						MoneyType: money.COIN,
						Name:      "10",
						Value:     10,
						Stock:     10,
					},
					{
						MoneyType: money.BANK,
						Name:      "20",
						Value:     20,
						Stock:     10,
					},
					{
						MoneyType: money.BANK,
						Name:      "50",
						Value:     50,
						Stock:     10,
					},
				},
				recievedMoney: map[money.Money]int8{},
			},
			expected: expectedArgs{
				expectedMoney: []money.Money{
					{
						MoneyType: money.COIN,
						Name:      "1",
						Value:     1,
					},
					{
						MoneyType: money.COIN,
						Name:      "1",
						Value:     1,
					},
					{
						MoneyType: money.COIN,
						Name:      "1",
						Value:     1,
					},
					{
						MoneyType: money.COIN,
						Name:      "5",
						Value:     5,
					},
					{
						MoneyType: money.COIN,
						Name:      "10",
						Value:     10,
					},
					{
						MoneyType: money.BANK,
						Name:      "20",
						Value:     20,
					},
					{
						MoneyType: money.BANK,
						Name:      "50",
						Value:     50,
					},
				},
			},
			hasError: false,
		},
		{
			description: "test_change_success_with_recieve_money",
			input: inputArgs{
				changeAmount: 8,
				availableMoney: []money.Money{
					{
						MoneyType: money.COIN,
						Name:      "1",
						Value:     1,
						Stock:     0,
					},
					{
						MoneyType: money.COIN,
						Name:      "5",
						Value:     5,
						Stock:     10,
					},
					{
						MoneyType: money.COIN,
						Name:      "10",
						Value:     10,
						Stock:     10,
					},
				},
				recievedMoney: map[money.Money]int8{
					{
						MoneyType: money.COIN,
						Name:      "1",
					}: 3,
				},
			},
			expected: expectedArgs{
				expectedMoney: []money.Money{
					{
						MoneyType: money.COIN,
						Name:      "1",
						Value:     1,
					},
					{
						MoneyType: money.COIN,
						Name:      "1",
						Value:     1,
					},
					{
						MoneyType: money.COIN,
						Name:      "1",
						Value:     1,
					},
					{
						MoneyType: money.COIN,
						Name:      "5",
						Value:     5,
					},
				},
			},
			hasError: false,
		},
		{
			description: "test_change_success_with_change_amount_equal_to_zero",
			input: inputArgs{
				changeAmount: 0,
				availableMoney: []money.Money{
					{
						MoneyType: money.COIN,
						Name:      "1",
						Value:     1,
						Stock:     10,
					},
					{
						MoneyType: money.COIN,
						Name:      "5",
						Value:     5,
						Stock:     10,
					},
					{
						MoneyType: money.COIN,
						Name:      "10",
						Value:     10,
						Stock:     10,
					},
				},
				recievedMoney: map[money.Money]int8{},
			},
			expected: expectedArgs{
				expectedMoney: []money.Money{},
			},
			hasError: false,
		},
		{
			description: "test_change_failed_insufficient_change",
			input: inputArgs{
				changeAmount: 8,
				availableMoney: []money.Money{
					{
						MoneyType: money.COIN,
						Name:      "1",
						Value:     1,
						Stock:     0,
					},
					{
						MoneyType: money.COIN,
						Name:      "5",
						Value:     5,
						Stock:     10,
					},
					{
						MoneyType: money.COIN,
						Name:      "10",
						Value:     10,
						Stock:     10,
					},
				},
				recievedMoney: map[money.Money]int8{},
			},
			expected: expectedArgs{
				expectedMoney: []money.Money{},
				expectedError: errors.New("insufficient change"),
			},
			hasError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {

			//sort descending input available money first for prioritize chage
			sort.Slice(test.input.availableMoney, func(i, j int) bool {
				return test.input.availableMoney[i].Value > test.input.availableMoney[j].Value
			})

			output, err := change(test.input.changeAmount, test.input.availableMoney, test.input.recievedMoney)
			if test.hasError {
				assert.Error(t, err)
				assert.Equal(t, test.expected.expectedError, err)
				assert.Equal(t, test.expected.expectedMoney, output)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected.expectedMoney, output)
			}
		})
	}
}
