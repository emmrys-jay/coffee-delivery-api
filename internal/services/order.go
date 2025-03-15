package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/emmrys-jay/coffee-delivery-api/internal/database/models"
	"github.com/emmrys-jay/coffee-delivery-api/internal/database/repository"
	"github.com/emmrys-jay/coffee-delivery-api/util"
	"github.com/shopspring/decimal"
)

type OrderService struct {
	repo        *repository.OrderRepository
	userRepo    *repository.UserRepository
	coffeeRepo  *repository.CoffeeRepository
	reserveRepo *repository.ReservationRepository
}

func NewOrderService(
	repo *repository.OrderRepository,
	userRepo *repository.UserRepository,
	coffeeRepo *repository.CoffeeRepository,
	reserveRepo *repository.ReservationRepository,
) *OrderService {
	return &OrderService{
		repo:        repo,
		userRepo:    userRepo,
		coffeeRepo:  coffeeRepo,
		reserveRepo: reserveRepo,
	}
}

func (os *OrderService) PlaceOrder(ctx context.Context, userId uint, req *models.CreateOrderRequest) (*models.OrderResponse, error) {
	_, err := os.userRepo.GetUserByID(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("error fetching user by id, %w", err)
	}

	var idMap = make(map[uint]uint)
	var ids = make([]uint, 0, len(req.Coffees))
	for _, v := range req.Coffees {
		ids = append(ids, v.CoffeeID)
		idMap[v.CoffeeID] = v.Quantity
	}

	coffees, err := os.coffeeRepo.GetMany(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("error fetching coffee products, %w", err)
	}

	if len(coffees) == 0 {
		return nil, errors.New("none of the coffee products specified was found")
	}

	// Check for the integrity of order quantity with quantity in stock
	// Calculate total order amount
	// Populate order items
	var totalAmount = decimal.Zero
	var orderItems = make([]models.OrderItem, 0, len(coffees))

	for _, v := range coffees {
		quantityOrdered := idMap[v.Id]

		if int(v.Quantity)-int(quantityOrdered) < 0 {
			errMsg := fmt.Sprintf("The quantity specified for '%s' is more than the quantity in stock: %v (specified) for %v (in stock)",
				v.Name, idMap[v.Id], v.Quantity)
			return nil, errors.New(errMsg)
		}

		price, _ := util.ParseDecimal(v.Price)
		totalAmount = totalAmount.Add(price.Mul(decimal.NewFromUint64(uint64(quantityOrdered))))

		orderItems = append(orderItems, models.OrderItem{
			CoffeeID:  v.Id,
			Name:      v.Name,
			Quantity:  quantityOrdered,
			UnitPrice: v.Price,
			CreatedAt: util.CurrentTime(),
		})
	}

	order := models.Order{
		UserID:      userId,
		TotalAmount: totalAmount.String(),
		OrderItems:  orderItems,
		Status:      models.ORDER_STATUS_PENDING,
	}

	_, err = os.repo.CreateOrder(ctx, &order)
	if err != nil {
		return nil, fmt.Errorf("error creating order, %w", err)
	}

	// now := util.CurrentTime()
	// reservations := make([]models.StockReservation, 0, len(order.OrderItems))
	// for _, v := range order.OrderItems {
	// 	r := models.StockReservation{
	// 		CoffeeId:         v.CoffeeID,
	// 		ReservedQuantity: v.Quantity,
	// 		OrderId:          order.Id,
	// 		CreatedAt:        now,
	// 	}

	// 	reservations = append(reservations, r)
	// }

	// if err := os.reserveRepo.ReserveProducts(ctx, reservations); err != nil {
	// 	return nil, fmt.Errorf("error reserving products: %w", err)
	// }

	orderResponse := order.ToOrderResponse()
	return &orderResponse, nil
}

func (os *OrderService) GetOrder(ctx context.Context, id uint) (*models.Order, error) {
	return os.repo.GetOrder(ctx, id)
}

func (os *OrderService) ListUserOrders(ctx context.Context, userId uint) ([]models.Order, error) {
	return os.repo.ListUserOrders(ctx, userId)
}

func (os *OrderService) UpdateOrderStatus(ctx context.Context, orderId uint, status string) (*models.Order, error) {
	retOrder, err := os.GetOrder(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("error fetching order, %w", err)
	}

	if !models.IsValidStatus(status) {
		return nil, errors.New("invalid status")
	}

	err = os.repo.UpdateOrderStatus(ctx, orderId, status)
	if err != nil {
		return nil, fmt.Errorf("error updating status, %w", err)
	}

	retOrder.Status = status // Add the updated status to the order struct to be returned
	return retOrder, nil
}

func (os *OrderService) CancelOrder(ctx context.Context, id uint) (*models.Order, error) {
	retOrder, err := os.GetOrder(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching order, %w", err)
	}

	if retOrder.Status != models.ORDER_STATUS_PENDING {
		return nil, errors.New("You cannot cancel this order again since it has already been processed. Please contact admin")
	}

	err = os.repo.UpdateOrderStatus(ctx, id, models.ORDER_STATUS_CANCELED)
	if err != nil {
		return nil, err
	}

	retOrder.Status = models.ORDER_STATUS_CANCELED // Add the updated status to the order struct to be returned
	return retOrder, nil
}
