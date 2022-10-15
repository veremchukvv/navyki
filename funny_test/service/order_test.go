package service

import (
	"context"
	"funny_test/service/mock"
	"funny_test/storage"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestOrder_Split(t *testing.T) {
	// prepare test data
	tests := []struct {
		description string

		order     storage.Order
		splitType SplitType

		expectedBills []storage.Bill
	}{
		{
			description: "Простое разделение с 2 людьми и несколькими блюдами",

			order: storage.Order{
				ID: 1,
				Dishes: []storage.Dish{
					{
						ID:       1,
						Name:     "колбаса",
						Price:    500,
						PersonID: 1,
					},
					{
						ID:       2,
						Name:     "сыр",
						Price:    300,
						PersonID: 2,
					},
				},
				PersonCount: 2,
			},
			splitType: Simple,

			expectedBills: []storage.Bill{
				{
					ID:       0,
					PersonID: 0,
					Amount:   800,
				},
			},
		},
		{
			description: "Простое разделение",

			order: storage.Order{
				ID: 1,
				Dishes: []storage.Dish{
					{
						ID:       1,
						Name:     "колбаса",
						Price:    500,
						PersonID: 1,
					},
				},
				PersonCount: 2,
			},
			splitType: Simple,

			expectedBills: []storage.Bill{
				{
					ID:       0,
					PersonID: 0,
					Amount:   500,
				},
			},
		},
		{
			description: "Разделение по блюдам",

			order: storage.Order{
				ID: 1,
				Dishes: []storage.Dish{
					{
						ID:       1,
						Name:     "колбаса",
						Price:    500,
						PersonID: 1,
					},
					{
						ID:       2,
						Name:     "сыр",
						Price:    300,
						PersonID: 1,
					},
					{
						ID:       2,
						Name:     "творог",
						Price:    400,
						PersonID: 2,
					},
				},
				PersonCount: 2,
			},
			splitType: ByDishes,

			expectedBills: []storage.Bill{
				{
					ID:       0,
					PersonID: 1,
					Amount:   800,
				},
				{
					ID:       0,
					PersonID: 2,
					Amount:   400,
				},
			},
		},
		{
			description: "Разделение счёта поровну",

			order: storage.Order{
				ID: 1,
				Dishes: []storage.Dish{
					{
						ID:       1,
						Name:     "колбаса",
						Price:    600,
						PersonID: 1,
					},
					{
						ID:       2,
						Name:     "сыр",
						Price:    300,
						PersonID: 2,
					},
					{
						ID:       3,
						Name:     "творог",
						Price:    400,
						PersonID: 3,
					},
				},
				PersonCount: 3,
			},
			splitType: ByPerson,

			expectedBills: []storage.Bill{
				{
					ID:       0,
					PersonID: 1,
					Amount:   433,
				},
				{
					ID:       1,
					PersonID: 2,
					Amount:   433,
				},
				{
					ID:       2,
					PersonID: 3,
					Amount:   434,
				},
			},
		},
	}
	cmpOpts := []cmp.Option{
		cmpopts.IgnoreFields(storage.Bill{}, "ID"),
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			mockStorage := mock.NewMockOrderStorage(ctrl)
			mockStorage.EXPECT().Order(gomock.Any(), gomock.Any()).Return(&tt.order).AnyTimes()

			svc := NewOrderService(mockStorage)

			svc.Split(context.Background(), tt.order.ID, tt.splitType)

			// execute method
			bills := svc.Split(context.Background(), tt.order.ID, tt.splitType)

			asrt := assert.New(t)

			asrt.NotNil(bills, "Счета не должны быть nil")
			asrt.NotEmpty(bills, "Счета должны существовать")

			if diff := cmp.Diff(tt.expectedBills, bills, cmpOpts...); diff != "" {
				t.Errorf("Счета не совпадают (-want +got):\n%s", diff)
			}
		})
	}
}
