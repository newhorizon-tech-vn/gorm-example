package main

import (
	"gorm.io/gorm"
)

type Item struct {
	ID         int            `json:"id" gorm:"column:id;primaryKey;"`
	ProductID  int            `json:"product_id" gorm:"column:product_id"`
	Name       string         `json:"name" gorm:"column:name"`
	Descripton string         `json:"descripton" gorm:"column:descripton"`
	DeletedAt  gorm.DeletedAt `json:"deleteDate" example:"31/12/9999 23:59" swaggertype:"string" gorm:"index;column:deleted_at"`
	Product    *Product       `json:"category" gorm:"foreignKey:id;references:product_id"`
}

func (*Item) TableName() string {
	return "Item"
}

type Category struct {
	ID        int            `json:"id" gorm:"column:id;primaryKey;"`
	Name      string         `json:"name" gorm:"column:name"`
	DeletedAt gorm.DeletedAt `json:"deleteDate" example:"31/12/9999 23:59" swaggertype:"string" gorm:"index;column:deleted_at"`
	Products  []*Product     `json:"products" gorm:"foreignKey:category_id;references:id;"`
}

func (*Category) TableName() string {
	return "Category"
}

type Product struct {
	ID         int    `json:"id" gorm:"column:id;primaryKey;"`
	CategoryID int    `json:"category_id" gorm:"column:category_id"`
	Name       string `json:"name" gorm:"column:name"`
	Price      int    `json:"price" gorm:"column:price"`
	// Library auto convert to plural words
	// Product_Factory_V2 => không đúng, hoặc là snake_case(product_factory_v2) hoặc là camelCase(ProductFactoryV2), không dùng phan trộn như vâyj
	// Factories []*Factory `gorm:"many2many:ProductFactory;"`

	Category *Category `json:"category" gorm:"foreignKey:id;references:category_id"` // colume name id, not í field name

	// DeletedAt  gorm.DeletedAt `json:"deleteDate" example:"31/12/9999 23:59" swaggertype:"string" gorm:"index;column:deleted_at"`
	Items []*Item `json:"items" gorm:"foreignKey:product_id;references:id"`

	// gorm:"many2many:<Join Table>;foreignKey:<primaryKey of Product>;joinForeignKey:<Product prefeKey of Product_Factory>;References:<primaryKey of Factory>;joinReferences:<Factory prefeKey of Product_Factory>
	Factories []*Factory `gorm:"many2many:product_factories;foreignKey:ID;joinForeignKey:ProductID;References:ID;joinReferences:FactoryID"`
}

func (*Product) TableName() string {
	return "Product"
}

type Factory struct {
	ID        int         `json:"id" gorm:"column:id;primaryKey;"`
	Name      string      `json:"name" gorm:"column:name"`
	Address   string      `json:"address" gorm:"column:address"`
	Workshops []*Workshop `json:"workshop" gorm:"foreignKey:FactoryID"`
	Products  []*Product  `gorm:"many2many:product_factories;foreignKey:ID;joinForeignKey:FactoryID;References:ID;joinReferences:ProductID"`
	// Products []Product `gorm:"many2many:ProductFactory;"`
}

func (*Factory) TableName() string {
	return "Factory"
}

type Workshop struct {
	ID        int    `json:"id" gorm:"column:id;primaryKey;"`
	FactoryID string `json:"factory_id" gorm:"column:factory_id"`
	Name      string `json:"name" gorm:"column:name"`
}

func (*Workshop) TableName() string {
	return "Workshop"
}

type ProductFactory struct {
	ID        int `json:"id" gorm:"column:id;primaryKey;"`
	ProductID int `json:"product_id" gorm:"column:product_id;"`
	FactoryID int `json:"factory_id" gorm:"column:factory_id;"`
	Price     int `json:"price" gorm:"column:price;"`
}

func (*ProductFactory) TableName() string {
	return "product_factories"
}
