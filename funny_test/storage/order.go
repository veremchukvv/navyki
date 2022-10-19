package storage

// Блюдо
type Dish struct {
	ID       int
	Name     string
	Price    int64
	PersonID int
}

// Заказ
type Order struct {
	ID int
	// блюда которые заказали
	Dishes []Dish
	// количество гостей
	PersonCount int64
}

// Счет
type Bill struct {
	ID       int
	PersonID int
	Amount   int64
	Discount int64
	Tax      string
}
