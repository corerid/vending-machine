package product

import (
	"errors"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func commonPrepData() []Product {
	return []Product{
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
}

func Test_SelectProduct(t *testing.T) {
	tests := []struct {
		description          string
		prepData             func()
		input                string
		expectedBuyedProduct map[Product]int8
		expectedTotalAmount  int64
		expectedError        string
		hasError             bool
	}{
		{
			description: "test_select_product_success",
			prepData: func() {
				commonPrepData()
			},
			input:               "1\n2\n4\n2\n4\n2\n\n",
			expectedTotalAmount: 65,
			expectedBuyedProduct: map[Product]int8{
				{
					ProductNo: 1,
					Name:      "Lays",
					Price:     5,
				}: 1,
				{
					ProductNo: 4,
					Name:      "Pepsi",
					Price:     15,
				}: 2,
				{
					ProductNo: 2,
					Name:      "Hanami",
					Price:     10,
				}: 3,
			},
			hasError: false,
		},
		{
			description: "test_select_product_success_with_product_out_of_stock",
			prepData: func() {
				commonPrepData()
			},
			input:               "1\n1\n2\n4\n2\n4\n2\n\n",
			expectedTotalAmount: 65,
			expectedBuyedProduct: map[Product]int8{
				{
					ProductNo: 1,
					Name:      "Lays",
					Price:     5,
				}: 1,
				{
					ProductNo: 4,
					Name:      "Pepsi",
					Price:     15,
				}: 2,
				{
					ProductNo: 2,
					Name:      "Hanami",
					Price:     10,
				}: 3,
			},
			hasError: false,
		},
		{
			description: "test_select_product_failed_user_not_select_any_product",
			prepData: func() {
				commonPrepData()
			},
			input:         "\n",
			expectedError: "you have not select any product",
			hasError:      true,
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

			_, err = io.WriteString(userInput, test.input)
			if err != nil {
				t.Fatal(err)
			}

			_, err = userInput.Seek(0, io.SeekStart)
			if err != nil {
				t.Fatal(err)
			}

			buyedProduct, totalAmount, err := SelectProduct(userInput)
			if test.hasError {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedTotalAmount, totalAmount)

				for k := range buyedProduct {
					assert.Equal(t, test.expectedBuyedProduct[k], buyedProduct[k])
				}
			}
		})
	}
}

func Test_checkProduct(t *testing.T) {
	type inputArgs struct {
		productNo    string
		productStock []Product
	}

	tests := []struct {
		description          string
		input                inputArgs
		expected             Product
		expectedProductStock []Product
		expectedError        error
		hasError             bool
	}{
		{
			description: "test_check_product_success",
			input: inputArgs{
				productNo:    "1",
				productStock: commonPrepData(),
			},
			expected: Product{

				ProductNo: 1,
				Name:      "Lays",
				Price:     5,
			},
			expectedProductStock: []Product{
				{
					ProductNo: 1,
					Name:      "Lays",
					Price:     5,
					Stock:     0,
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
			},
			hasError: false,
		},
		{
			description: "test_check_product_failed_invalid_product_no",
			input: inputArgs{
				productNo:    "xxx",
				productStock: commonPrepData(),
			},
			expected:             Product{},
			expectedProductStock: commonPrepData(),
			expectedError:        errors.New("invalid input"),
			hasError:             true,
		},
		{
			description: "test_check_product_failed_product_is_out_of_stock",
			input: inputArgs{
				productNo: "2",
				productStock: []Product{
					{
						ProductNo: 2,
						Name:      "Hanami",
						Price:     10,
						Stock:     0,
					},
				},
			},
			expected: Product{},
			expectedProductStock: []Product{
				{
					ProductNo: 2,
					Name:      "Hanami",
					Price:     10,
					Stock:     0,
				},
			},
			expectedError: errors.New("Hanami is out of stock"),
			hasError:      true,
		},
		{
			description: "test_check_product_failed_product_does_not_exist",
			input: inputArgs{
				productNo: "1",
				productStock: []Product{
					{
						ProductNo: 2,
						Name:      "Hanami",
						Price:     10,
						Stock:     0,
					},
				},
			},
			expected: Product{},
			expectedProductStock: []Product{
				{
					ProductNo: 2,
					Name:      "Hanami",
					Price:     10,
					Stock:     0,
				},
			},
			expectedError: errors.New("product doesn't exist"),
			hasError:      true,
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			output, err := checkProduct(test.input.productNo, test.input.productStock)
			if test.hasError {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError, err)
				assert.Equal(t, test.expected, output)
				assert.Equal(t, test.expectedProductStock, test.input.productStock)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, output)
				assert.Equal(t, test.expectedProductStock, test.input.productStock)
			}
		})
	}
}

func Test_DecreaseStock(t *testing.T) {
	tests := []struct {
		description   string
		prepData      func()
		input         map[Product]int8
		expected      []Product
		expectedError error
		hasError      bool
	}{
		{
			description: "test_decrese_stock_success",
			prepData: func() {
				ProductStock = []Product{
					{
						ProductNo: 1,
						Name:      "Sunbyte",
						Price:     10,
						Stock:     10,
					},
					{
						ProductNo: 2,
						Name:      "Papika",
						Price:     5,
						Stock:     10,
					},
				}
			},
			input: map[Product]int8{
				{
					ProductNo: 1,
					Name:      "Sunbyte",
				}: 1,
				{
					ProductNo: 2,
					Name:      "Papika",
				}: 10,
			},
			expected: []Product{
				{
					ProductNo: 1,
					Name:      "Sunbyte",
					Price:     10,
					Stock:     9,
				},
				{
					ProductNo: 2,
					Name:      "Papika",
					Price:     5,
					Stock:     0,
				},
			},
			hasError: false,
		},
		{
			description: "test_decrese_stock_failed_stock_is_less_than_zero",
			prepData: func() {
				ProductStock = []Product{
					{
						ProductNo: 1,
						Name:      "Sunbyte",
						Price:     10,
						Stock:     10,
					},
					{
						ProductNo: 2,
						Name:      "Papika",
						Price:     5,
						Stock:     10,
					},
				}
			},
			input: map[Product]int8{
				{
					ProductNo: 1,
					Name:      "Sunbyte",
				}: 1,
				{
					ProductNo: 2,
					Name:      "Papika",
				}: 11,
			},
			expected: []Product{
				{
					ProductNo: 1,
					Name:      "Sunbyte",
					Price:     10,
					Stock:     9,
				},
				{
					ProductNo: 2,
					Name:      "Papika",
					Price:     5,
					Stock:     10,
				},
			},
			expectedError: errors.New("Papika's stock is less than zero"),
			hasError:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			test.prepData()
			err := DecreaseStock(test.input)
			if test.hasError {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, ProductStock)
			}
		})
	}
}
