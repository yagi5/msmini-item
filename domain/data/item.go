package data

import (
	"fmt"
	"time"
)

// Category is custom type that represents item category
type Category string

const (
	// Book book
	Book Category = "Book"
	// Food food
	Food Category = "Food"
)

// Item is entity
type Item struct {
	ID          string
	Name        string
	Description string
	Price       int64
	Category    Category
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (i *Item) String() string {
	return fmt.Sprintf("%s: %s($%d)", i.Category, i.Name, i.Price)
}
