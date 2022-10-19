package storage

// Блюдо
type Dish struct {
	ID       int
	Name     string
	Price    int
	PersonID int
}

// Заказ
type Order struct {
	ID int
	// блюда которые заказали
	Dishes []Dish
	// количество гостей
	PersonCount int
}

// Счет
type Bill struct {
	ID       int
	PersonID int
	Amount   int
	Discount int
}
