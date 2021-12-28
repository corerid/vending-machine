# Vending-Machine

### Project structure
```
    .
    └── src                    
        └── vending-machine        
            ├── money
            │    └── ...
            ├── payment
            │   └── ...
            ├── product
            │    └── ...
            ├── main.go            
            └── ...
```

### Sync dependencies
```
$ cd src/vending-machine
$ go mod tody
```

### Run program
```
$ go run main.go
```

### Program Instruction
```
1. start program
2. select product by type product no.
3. if you want to select more product, you can continue select product (same as 2.) 
   or if you finish select product, just press ENTER key to checkout
4. insert money (each at a time) that accepted (1, 5 or 10) until you insert money more than total product's amount
5. If sucessful, program will show purchase summary
```

[![IMG-0088.jpg](https://i.postimg.cc/zDwCk3Hp/IMG-0088.jpg)](https://postimg.cc/QVtK88bW)

### Config stock
```
you can change 
- product's stock at src/vending-machine/product/product.go variable "ProductStock"
- money's stock at src/vending-machine/money/money.go variable "MoneyStock"
```
