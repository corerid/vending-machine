package money

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_checkMoney(t *testing.T) {
	type inputArgs struct {
		moneyName string
	}

	tests := []struct {
		description   string
		prepData      func()
		input         inputArgs
		expected      Money
		expectedError error
		hasError      bool
	}{
		{
			description: "test_check_money_success",
			prepData: func() {
				MoneyStock = []Money{
					{
						MoneyType: COIN,
						Name:      "1",
						Value:     1,
						Stock:     10,
					},
					{
						MoneyType: COIN,
						Name:      "5",
						Value:     5,
						Stock:     10,
					},
				}
			},
			input: inputArgs{
				moneyName: "1",
			},
			expected: Money{
				MoneyType: COIN,
				Name:      "1",
				Value:     1,
				Stock:     10,
			},
			hasError: false,
		},
		{
			description: "test_check_money_failed_money_does_not_accepted",
			prepData: func() {
				MoneyStock = []Money{
					{
						MoneyType: COIN,
						Name:      "1",
						Value:     1,
						Stock:     10,
					},
					{
						MoneyType: COIN,
						Name:      "5",
						Value:     5,
						Stock:     10,
					},
				}
			},
			input: inputArgs{
				moneyName: "100",
			},
			expected:      Money{},
			expectedError: errors.New("money doesn't excepted"),
			hasError:      true,
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			test.prepData()
			output, err := CheckMoney(test.input.moneyName)
			if test.hasError {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError, err)
				assert.Equal(t, test.expected, output)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, output)
			}
		})
	}
}

func Test_IncreaseStock(t *testing.T) {
	tests := []struct {
		description string
		prepData    func()
		input       map[Money]int8
		expected    []Money
		hasError    bool
	}{
		{
			description: "test_increse_stock_success",
			prepData: func() {
				MoneyStock = []Money{
					{
						MoneyType: COIN,
						Name:      "1",
						Value:     1,
						Stock:     0,
					},
					{
						MoneyType: COIN,
						Name:      "5",
						Value:     5,
						Stock:     0,
					},
					{
						MoneyType: COIN,
						Name:      "10",
						Value:     10,
						Stock:     0,
					},
					{
						MoneyType: BANK,
						Name:      "20",
						Value:     20,
						Stock:     10,
					},
				}
			},
			input: map[Money]int8{
				{
					Name: "1",
				}: 1,
				{
					Name: "5",
				}: 2,
				{
					Name: "10",
				}: 3,
			},
			expected: []Money{
				{
					MoneyType: COIN,
					Name:      "1",
					Value:     1,
					Stock:     1,
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
					Stock:     3,
				},
				{
					MoneyType: BANK,
					Name:      "20",
					Value:     20,
					Stock:     10,
				},
			},
			hasError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			test.prepData()
			err := IncreaseStock(test.input)
			if test.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, MoneyStock)
			}
		})
	}
}

func Test_DecreaseStock(t *testing.T) {
	tests := []struct {
		description   string
		prepData      func()
		input         []Money
		expected      []Money
		expectedError error
		hasError      bool
	}{
		{
			description: "test_decrese_stock_success",
			prepData: func() {
				MoneyStock = []Money{
					{
						MoneyType: COIN,
						Name:      "1",
						Value:     1,
						Stock:     1,
					},
					{
						MoneyType: COIN,
						Name:      "5",
						Value:     5,
						Stock:     3,
					},
					{
						MoneyType: COIN,
						Name:      "10",
						Value:     10,
						Stock:     5,
					},
					{
						MoneyType: BANK,
						Name:      "20",
						Value:     20,
						Stock:     10,
					},
				}
			},
			input: []Money{
				{
					MoneyType: COIN,
					Name:      "1",
					Value:     1,
				},
				{
					MoneyType: COIN,
					Name:      "5",
					Value:     5,
				},
				{
					MoneyType: COIN,
					Name:      "5",
					Value:     5,
				},
				{
					MoneyType: COIN,
					Name:      "10",
					Value:     10,
				},
				{
					MoneyType: COIN,
					Name:      "10",
					Value:     10,
				},
				{
					MoneyType: COIN,
					Name:      "10",
					Value:     10,
				},
			},
			expected: []Money{
				{
					MoneyType: COIN,
					Name:      "1",
					Value:     1,
					Stock:     0,
				},
				{
					MoneyType: COIN,
					Name:      "5",
					Value:     5,
					Stock:     1,
				},
				{
					MoneyType: COIN,
					Name:      "10",
					Value:     10,
					Stock:     2,
				},
				{
					MoneyType: BANK,
					Name:      "20",
					Value:     20,
					Stock:     10,
				},
			},
			hasError: false,
		},
		{
			description: "test_decrese_stock_failed_stock_is_less_than_zero",
			prepData: func() {
				MoneyStock = []Money{
					{
						MoneyType: COIN,
						Name:      "1",
						Value:     1,
						Stock:     1,
					},
					{
						MoneyType: COIN,
						Name:      "5",
						Value:     5,
						Stock:     1,
					},
					{
						MoneyType: COIN,
						Name:      "10",
						Value:     10,
						Stock:     5,
					},
				}
			},
			input: []Money{
				{
					MoneyType: COIN,
					Name:      "1",
					Value:     1,
				},
				{
					MoneyType: COIN,
					Name:      "5",
					Value:     5,
				},
				{
					MoneyType: COIN,
					Name:      "5",
					Value:     5,
				},
			},
			expectedError: errors.New("5's stock is less than zero"),
			hasError:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			test.prepData()
			err := DecreaseStock(test.input)
			if test.hasError {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, MoneyStock)
			}
		})
	}
}
