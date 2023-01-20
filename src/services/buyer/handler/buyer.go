package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/db/yugabyte"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/buyer/domain"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/buyer/usecase"
	"go.uber.org/fx"
	"gorm.io/datatypes"
)

type Handler interface {
	Auth() gin.HandlerFunc
	Login(ctx *gin.Context)
	Products(ctx *gin.Context)
	ProductByID(ctx *gin.Context)
	Orders(ctx *gin.Context)
	OrderByID(ctx *gin.Context)
	UpdateOrderStatus(ctx *gin.Context)
	CreateOrder(ctx *gin.Context)
}

type handler struct {
	BuyerUsecase usecase.BuyerUsecase
	OrderUsecase usecase.OrderUsecase
}

type Params struct {
	fx.In
	BuyerUsecase usecase.BuyerUsecase
	OrderUsecase usecase.OrderUsecase
}

func NewBuyerHandler(param Params) Handler {
	return &handler{
		BuyerUsecase: param.BuyerUsecase,
		OrderUsecase: param.OrderUsecase,
	}
}

func (h *handler) Login(ctx *gin.Context) {
	session := sessions.Default(ctx)
	request := new(LoginRequest)
	if err := ctx.Bind(request); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest,
			LoginResponse{
				Error: "invalid body type",
			})
		return
	}

	res, err := h.BuyerUsecase.Login(ctx, request.Username)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError,
			LoginResponse{
				Error: "something happened on our end, please try at a later time",
			})
		return
	}

	// set buyer object to cookie
	session.Set(domain.BuyerKey, res.ID)
	session.Save()

	ctx.JSON(http.StatusOK, LoginResponse{
		Data: res,
	})
}

// Products
func (h *handler) Products(ctx *gin.Context) {
	var (
		res []domain.Product
		err error
	)

	res, err = h.OrderUsecase.Products(ctx)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, ProductsResponse{
			Error: "something happened on our end, please try at a later time",
		})
		return
	}

	ctx.JSON(http.StatusOK, ProductsResponse{
		Data: &res,
	})
}

// Products
func (h *handler) ProductByID(ctx *gin.Context) {
	var (
		parsedId uint64
		err      error
	)

	strProductId := ctx.Param("id")
	if strProductId != "" {
		parsedId, err = strconv.ParseUint(strProductId, 10, 0)
		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusBadRequest, ProductByIDResponse{
				Error: "please pass a valid id",
			})
			return
		}
	}

	res, err := h.OrderUsecase.ProductByID(ctx, uint(parsedId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ProductByIDResponse{
			Error: "something happened on our end, please try at a later time",
		})
		return
	}

	ctx.JSON(http.StatusOK, ProductByIDResponse{
		Data: res,
	})
}

// Orders
func (h *handler) Orders(ctx *gin.Context) {
	var (
		res []domain.Order
		err error
	)

	session := sessions.Default(ctx)
	buyerId := session.Get(domain.BuyerKey).(uint)

	res, err = h.OrderUsecase.OrdersByBuyer(ctx, buyerId)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, OrdersResponse{
			Error: "something happened on our end, please try at a later time",
		})
		return
	}

	ctx.JSON(http.StatusOK, OrdersResponse{
		Data: &res,
	})
}

// UpdateOrder
func (h *handler) UpdateOrderStatus(ctx *gin.Context) {
	request := new(UpdateOrderRequest)
	if err := ctx.Bind(request); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, OrderResponse{
			Error: "invalid body type",
		})
		return
	}

	res, err := h.OrderUsecase.UpdateOrderStatus(ctx, request.ID, request.Status)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, OrderResponse{
			Error: "something happened on our end, please try at a later time",
		})
		return
	}

	ctx.JSON(http.StatusOK, OrderResponse{
		Data: res,
	})
}

// CreateOrder
func (h *handler) CreateOrder(ctx *gin.Context) {
	var (
		err error
	)

	request := new(CreateOrderRequest)
	if err := ctx.Bind(request); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, OrderResponse{
			Error: "invalid body type",
		})
		return
	}

	session := sessions.Default(ctx)
	buyerId := session.Get(domain.BuyerKey).(uint)

	// convert handler request to domain order
	var orderDetails []domain.OrderDetail
	for _, v := range request.Products {
		item := domain.OrderDetail{
			Product: domain.Product{
				Model: yugabyte.Model{
					ID: v.ProductID,
				},
			},
			ProductID:       v.ProductID,
			ProductQuantity: int(v.ProductQty),
		}

		orderDetails = append(orderDetails, item)
	}

	order := domain.Order{
		OrderDate:    datatypes.Date(time.Now()),
		BuyerID:      buyerId,
		Status:       domain.OrderStatusNew,
		InvoiceNo:    "INV-" + time.Now().Format("20060102150405"),
		OrderDetails: orderDetails,
	}

	res, err := h.OrderUsecase.CreateOrder(ctx, order)
	if err != nil {
		if strings.Contains(err.Error(), "ongoing order") {
			ctx.JSON(http.StatusBadRequest, OrderResponse{
				Error: "cannot create order, theres an ongoing order",
			})
		} else {
			ctx.Error(err)
			ctx.JSON(http.StatusInternalServerError, OrderResponse{
				Error: "something happened on our end, please try at a later time : " + err.Error(),
			})
		}

		return
	}

	ctx.JSON(http.StatusOK, OrderResponse{
		Data: res,
	})
}

func (h *handler) OrderByID(ctx *gin.Context) {
	var (
		parsedId uint64
		err      error
	)

	strOrderId := ctx.Param("id")
	if strOrderId != "" {
		parsedId, err = strconv.ParseUint(strOrderId, 10, 0)
		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusBadRequest, OrderResponse{
				Error: "please pass order id to path",
			})
			return
		}
	}

	res, err := h.OrderUsecase.OrderByID(ctx, uint(parsedId))
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, OrderResponse{
			Error: "something happened on our end, please try at a later time",
		})
		return
	}

	ctx.JSON(http.StatusOK, OrderResponse{
		Data: res,
	})
}
