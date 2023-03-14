package main

import "fmt"

func selectProductsV2() ([]Product, error) {
	var result []Product
	err := DBClient.SetupJoinTable(&Product{}, "Factories", &ProductFactory{})
	if err != nil {
		return nil, err
	}

	err = DBClient.Preload("Items").Preload("Factories").Preload("Factories.Workshops").Find(&result).Error
	return result, err
}

func testSetupJoinTable() {
	products, err := selectProductsV2()
	if err != nil {
		fmt.Println("ERROR", "Get product by setJoinTable failed", "error", err)
		return
	}

	fmt.Println(products)
}
