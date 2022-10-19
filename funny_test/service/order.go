package service

import (
	"context"
	"funny_test/storage"
	"sort"
	"strconv"
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
	var amount int64
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

	//создаём массив счетов 0 размера, чтобы избежать добавления пустых счетов
	bills := make([]storage.Bill, 0)

	for _, b := range billsByUser {
		bills = append(bills, b)
	}

	//сортируем список счетов для дальнейших операций
	sort.Slice(bills, func(i, j int) (less bool) {
		return bills[i].PersonID < bills[j].PersonID
	})

	return bills
}

func dishSplitEqual(o *storage.Order) []storage.Bill {
	billsByUser := make(map[int]storage.Bill)

	var amount int64

	for _, dish := range o.Dishes {
		amount = amount + dish.Price
	}

	for _, dish := range o.Dishes {
		bill := storage.Bill{}

		bill.ID = dish.PersonID
		bill.PersonID = dish.PersonID
		// bill.Amount = amount / o.PersonCount
		billsByUser[dish.PersonID] = bill
	}

	//создаём массив счетов 0 размера, чтобы избежать добавления пустых счетов
	bills := make([]storage.Bill, 0)

	for _, b := range billsByUser {
		bills = append(bills, b)
	}

	//сортируем список счетов для дальнейших операций
	sort.Slice(bills, func(i, j int) (less bool) {
		return bills[i].PersonID < bills[j].PersonID
	})

	//считаем дробную часть общей суммы
	fraction := amount % o.PersonCount

	//объявляем нулевую скидку
	var discount int64

	//добиваемся деления общей суммы счёта на количество гостей без остатка и считаем размер скидки
	if fraction != 0 {
		var i int64
		for {
			i++
			amount -= 1
			newFraction := amount % o.PersonCount
			if newFraction == 0 {
				discount = i
				break
			}
		}
	}

	// считаем сумму счёта для каждого гостя от общей суммы счёта за вычетом скидки
	amountForBill := amount / o.PersonCount

	// добавляем в счёт каждому гостю сумму поделённую на количество гостей без остатка
	for i := range billsByUser {
		bills[i-1].Amount += amountForBill
	}

	// добавляем в счёт последнему гостю информацию о размере скидки (для учёта в бухгалтерии)
	bills[len(bills)-1].Discount += discount

	TaxFraction := ((bills[1].Amount) * 20) % 100

	for i := range billsByUser {
		NDS := ((bills[i-1].Amount) * 20) / 100
		bills[i-1].Tax = strconv.FormatInt(NDS, 10) + "," + strconv.FormatInt(TaxFraction, 10)
	}

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
