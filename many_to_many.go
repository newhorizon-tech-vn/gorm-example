package main

import (
	"fmt"

	"gorm.io/gorm"
)

func addCategoryWithFactoriesAndWorkshops() (int, error) {
	category := &Category{
		Name: "Category_4",
	}

	factory5 := &Factory{
		Name:      "Factory_5",
		Workshops: []*Workshop{&Workshop{Name: "Workshop_5_1"}, &Workshop{Name: "Workshop_5_2"}},
	}
	factory6 := &Factory{
		Name:      "Factory_6",
		Workshops: []*Workshop{&Workshop{Name: "Workshop_6_1"}, &Workshop{Name: "Workshop_6_2"}},
	}

	// Factory phải dùng dạng tham chiếu để sau khi Create Factory ID sẽ được cập nhật
	// Nếu để dạng tham trị thì các factory tự đầu đã được gán và clone cho các product
	category.Products = append(category.Products, &Product{
		Name:      "Product_4_1",
		Price:     1000,
		Items:     []*Item{&Item{Name: "Item_4_1_1"}, &Item{Name: "Item_4_1_2"}},
		Factories: []*Factory{factory5, factory6},
	})
	category.Products = append(category.Products, &Product{
		Name:      "Product_4_2",
		Price:     1000,
		Items:     []*Item{&Item{Name: "Item_4_2_1"}, &Item{Name: "Item_4_2_2"}},
		Factories: []*Factory{factory5},
	})

	err := DBClient.Transaction(func(tx *gorm.DB) error {
		factories := []*Factory{factory5, factory6}
		// generate factory 5 id
		if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Create(factories).Error; err != nil {
			return err
		}

		if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Create(category).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return category.ID, nil
}

func getAllProducts() ([]Product, error) {
	var result []Product
	err := DBClient.Preload("Factories").Find(&result).Error
	return result, err
}

func getAllProductsAndFactoriesAndWorkshopsByFactoryID(factoryID int) ([]Product, error) {
	var result []Product
	err := DBClient.Preload("Factories", "id = ?", factoryID).Preload("Factories.Workshops", "name LIKE ?", "%_1%").Find(&result).Error
	return result, err
}

func getProductsAndFactoriesAndWorkshopsByFactoryID(productID, factoryID int) ([]Product, error) {
	var result []Product
	err := DBClient.Preload("Factories", "id = ?", factoryID).Preload("Factories.Workshops", "name LIKE ?", "%_1%").Where("id = ?", productID).Find(&result).Error
	return result, err
}

func getProducts() ([]Product, error) {
	var result []Product
	err := DBClient.Preload("Items").Preload("Factories").Preload("Factories.Workshops").Find(&result).Error
	return result, err
}

func getCategoryByID(id int) (Category, error) {
	var result Category
	err := DBClient.Preload("Products").Preload("Products.Items").Preload("Products.Factories").Preload("Products.Factories.Workshops").Where("id = ?", id).Find(&result).Error
	return result, err
}

func testManyToMany() {
	products, err := getProductsAndFactoriesAndWorkshopsByFactoryID(1, 2)
	if err != nil {
		fmt.Println("ERROR", "Get all product failed", err)
		return
	}
	fmt.Println("DEBUG", "Get all product success", "products", products)
}
