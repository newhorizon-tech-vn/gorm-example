package main

import (
	"fmt"

	"golang.org/x/exp/slices"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func getCategory(id int) (*Category, error) {
	result := &Category{}

	err := DBClient.Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getCategoryWithProductsV1(id int) (*Category, error) {
	result := &Category{}

	// load to children
	err := DBClient.Preload(clause.Associations).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getCategoryWithProducts(id int) (*Category, error) {
	result := &Category{}

	// load to children
	err := DBClient.Preload("Products").Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getCategoryWithProductsByCond(id int) (*Category, error) {
	result := &Category{}

	// load to children
	err := DBClient.Preload("Products", "name LIKE ?", "%_2%").Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getCategoryWithProductsAndItems(id int) (*Category, error) {
	result := &Category{}
	err := DBClient.Preload("Products").Preload("Products.Items").Preload(clause.Associations).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getCategoryWithProductsAndItemsByCond(id int) (*Category, error) {
	result := &Category{}

	// Load to grandchildren
	// err := DBClient.Preload("Product.Item").Preload(clause.Associations).Where("id = ?", id).First(&result).Error
	// err := DBClient.Preload("Product.Items").Preload(clause.Associations).Where("id = ?", id).First(&result).Error
	err := DBClient.Preload("Products").Preload("Products.Items", "name LIKE ?", "%Item_1_1%").Preload(clause.Associations).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func addCategoryWithProductsAndItems() (int, error) {
	category := &Category{
		Name: "Category_4",
	}
	category.Products = append(category.Products, &Product{
		Name:  "Product_4_1",
		Price: 1000,
		Items: []*Item{&Item{Name: "Item_4_1_1"}, &Item{Name: "Item_4_1_2"}},
	})
	category.Products = append(category.Products, &Product{
		Name:  "Product_4_2",
		Price: 1000,
		Items: []*Item{&Item{Name: "Item_4_2_1"}, &Item{Name: "Item_4_2_2"}},
	})

	err := DBClient.Transaction(func(tx *gorm.DB) error {
		// if err := tx.Create(category).Error; err != nil {
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

func updateCategoryWithProductsAndItems(category Category) error {
	// remove last element
	category.Products = slices.Delete(category.Products, len(category.Products)-1, len(category.Products))
	// add new element
	category.Products = append(category.Products, &Product{
		Name:  "Product_4_3",
		Price: 1000,
		Items: []*Item{&Item{Name: "Item_4_3_1"}, &Item{Name: "Item_4_3_2"}},
	})

	return DBClient.Transaction(func(tx *gorm.DB) error {
		// bắt buộc Product phải có delete_at
		return tx.Table((&Product{}).TableName()).Model(&category).Association("Products").Replace(category.Products)
		// return tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(category).Error
	})
}

func testHasMany() {
	// add category
	categoryId, err := addCategoryWithProductsAndItems()
	if err != nil {
		fmt.Println("ERROR:", "add category failed", err)
		return
	}
	fmt.Println("DEBUG:", "add category success", "category_id", categoryId)

	// get category deatails
	category, err := getCategoryWithProductsAndItems(categoryId)
	if err != nil {
		fmt.Println("ERROR:", "get category after add failed", err)
		return
	}
	fmt.Println("DEBUG:", "get category after add success", category)
	DumpCategoryData(category)

	// update category
	err = updateCategoryWithProductsAndItems(*category)
	if err != nil {
		fmt.Println("ERROR:", "update category failed", err)
		return
	}
	fmt.Println("DEBUG:", "update category success", "category_id", categoryId)

	// get category deatails
	category, err = getCategoryWithProductsAndItems(categoryId)
	if err != nil {
		fmt.Println("ERROR:", "get category after update failed", err)
		return
	}
	fmt.Println("DEBUG:", "get category after update success", category)

	DumpCategoryData(category)
}
