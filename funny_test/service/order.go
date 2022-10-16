package service

import (
	"context"
	"funny_test/storage"
	"sort"
)

type OrderStorage interface {
	Order(ctx context.Context, ID int) *storage.Order
}

type OrderService struct {
	store OrderStorage
}

func NewOrderService(store OrderStorage) *OrderService {
	return &OrderService{store: store}
}

// способ разделения счета
type SplitType int64

const (
	Simple   SplitType = iota // счет не разделяется
	ByDishes                  // cчет разделяется по людям которые заказали блюда
	ByPerson                  // счет разделяется по всем людям поровну
)

func simpleSplit(o *storage.Order) []storage.Bill {
	amount := 0
	for _, dish := range o.Dishes {
		amount = amount + dish.Price
	}

	return []storage.Bill{{
		ID:       0,
		PersonID: 0,
		Amount:   amount,
	}}
}

func dishSplit(o *storage.Order) []storage.Bill {
	billsByUser := make(map[int]storage.Bill)

	for _, dish := range o.Dishes {
		bill := storage.Bill{}
		if b, ok := billsByUser[dish.PersonID]; ok {
			bill = b
		}

		bill.PersonID = dish.PersonID
		bill.Amount = bill.Amount + dish.Price

		billsByUser[dish.PersonID] = bill
	}

	bills := make([]storage.Bill, 0)

	for _, b := range billsByUser {
		bills = append(bills, b)
	}

	sort.Slice(bills, func(i, j int) (less bool) {
		return bills[i].PersonID < bills[j].PersonID
	})

	return bills
}

func dishSplitEqual(o *storage.Order) []storage.Bill {
	billsByUser := make(map[int]storage.Bill)

	amount := 0

	for _, dish := range o.Dishes {
		amount = amount + dish.Price
	}

	//считаем дробную часть общей суммы
	change := amount % o.PersonCount

	for _, dish := range o.Dishes {
		bill := storage.Bill{}

		bill.ID = dish.PersonID
		bill.PersonID = dish.PersonID
		bill.Amount = amount / o.PersonCount
		billsByUser[dish.PersonID] = bill
	}

	bills := make([]storage.Bill, 0)

	for _, b := range billsByUser {
		bills = append(bills, b)
	}

	sort.Slice(bills, func(i, j int) (less bool) {
		return bills[i].PersonID < bills[j].PersonID
	})

	//добавляем дробную часть от суммы в счёт последнему гостю
	bills[len(bills)-1].Amount += change

	return bills
}

func (svc *OrderService) Split(ctx context.Context, ID int, splitType SplitType) []storage.Bill {
	order := svc.store.Order(ctx, ID)

	switch splitType {
	case Simple:
		return simpleSplit(order)
	case ByDishes:
		return dishSplit(order)
	case ByPerson:
		return dishSplitEqual(order)
	}
	return nil
}
