package main

import "fmt"

func DumpCategoryData(c *Category) {
	fmt.Println("Category", "id", c.ID, "name", c.Name)
	for _, product := range c.Products {
		fmt.Println("==Product", "id", product.ID, "name", product.Name)
		for _, item := range product.Items {
			fmt.Println("====Item", "id", item.ID, "name", item.Name)
		}
	}
}
