package main

import (
	"fmt"

	"gorm.io/gorm"
)

func addProductWithFactoriesAndWorkshops() (int, error) {

	factory5 := &Factory{
		Name:      "Factory_7",
		Workshops: []*Workshop{&Workshop{Name: "Workshop_7_1"}, &Workshop{Name: "Workshop_7_2"}},
	}
	factory6 := &Factory{
		Name:      "Factory_8",
		Workshops: []*Workshop{&Workshop{Name: "Workshop_8_1"}, &Workshop{Name: "Workshop_8_2"}},
	}

	product := &Product{
		Name:      "Product_8_1",
		Price:     1000,
		Items:     []*Item{&Item{Name: "Item_8_1_1"}, &Item{Name: "Item_8_1_2"}},
		Factories: []*Factory{factory5, factory6},
	}

	err := DBClient.Transaction(func(tx *gorm.DB) error {
		if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Create(product).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return product.ID, nil
}

func selectProductsV2() ([]Product, error) {
	var result []Product
	err := DBClient.SetupJoinTable(&Product{}, "Factories", &ProductFactory{})
	if err != nil {
		return nil, err
	}

	err = DBClient.Preload("Items").Preload("Factories").Preload("Factories.Workshops").Find(&result).Error
	return result, err
}

func selectProductByID(productId int) (Product, error) {
	var result Product
	err := DBClient.Preload("Items").Preload("Factories").Preload("Factories.Workshops").First(&result, productId).Error
	return result, err
}

func testSetupJoinTable() {
	productID, err := addProductWithFactoriesAndWorkshops()
	if err != nil {
		fmt.Println("ERROR", "Add product with factory failed", "error", err)
		return
	}

	product, err := selectProductByID(productID)
	if err != nil {
		fmt.Println("ERROR", "Get product by ID failed", "error", err)
		return
	}

	DumpProductData(&product)
}
