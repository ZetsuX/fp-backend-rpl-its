package controller

import (
	"fp-rpl/common"
	"fp-rpl/dto"
	"fp-rpl/entity"
	"fp-rpl/service"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type transactionController struct {
	transactionService service.TransactionService
	sessionService     service.SessionService
	spotService        service.SpotService
	userService        service.UserService
}

type TransactionController interface {
	MakeTransaction(ctx *gin.Context)
	GetAllTransactions(ctx *gin.Context)
	GetTransactionsByUsername(ctx *gin.Context)
	GetMyTransactions(ctx *gin.Context)
	DeleteTransactionByID(ctx *gin.Context)
}

func NewTransactionController(transactionS service.TransactionService, sessionS service.SessionService, spotS service.SpotService, userS service.UserService) TransactionController {
	return &transactionController{
		transactionService: transactionS,
		sessionService:     sessionS,
		spotService:        spotS,
		userService:        userS,
	}
}

func (transactionC *transactionController) MakeTransaction(ctx *gin.Context) {
	var transactionDTO dto.TransactionMakeRequest
	err := ctx.ShouldBind(&transactionDTO)
	if err != nil {
		resp := common.CreateFailResponse("failed to process transaction make request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	userId := ctx.GetUint64("ID")
	transactionDTO.UserID = userId

	sessionId, err := strconv.ParseUint(ctx.Param("sessionid"), 10, 64)
	if err != nil {
		resp := common.CreateFailResponse("failed to process session id", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	session, err := transactionC.sessionService.GetSessionByID(ctx, sessionId)
	if err != nil {
		resp := common.CreateFailResponse(err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	if reflect.DeepEqual(session, entity.Session{}) {
		resp := common.CreateFailResponse("session not found", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	transactionDTO.SessionID = sessionId

	var spots []entity.Spot
	for _, spotName := range transactionDTO.SpotsName {
		spotNum, err := strconv.ParseInt(spotName[1:], 10, 0)
		if err != nil {
			resp := common.CreateFailResponse("failed to process spot name", http.StatusBadRequest)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		}

		spot, err := transactionC.spotService.GetSpotBySessionIDAndAttributes(ctx, sessionId, string(spotName[0]), int(spotNum))
		if err != nil {
			resp := common.CreateFailResponse("failed to process spot", http.StatusBadRequest)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		}

		if reflect.DeepEqual(spot, entity.Spot{}) {
			resp := common.CreateFailResponse("spot with name "+spotName+" not found", http.StatusBadRequest)
			ctx.JSON(http.StatusBadRequest, resp)
			return
		}

		if spot.TransactionID != nil {
			resp := common.CreateFailResponse("spot with name "+spotName+" is reserved", http.StatusBadRequest)
			ctx.JSON(http.StatusBadRequest, resp)
			return
		}

		spots = append(spots, spot)
	}

	transactionDTO.TotalPrice = (session.Price * float64(len(transactionDTO.SpotsName)))
	transactionDTO.Code = uuid.NewString()

	newTransaction, err := transactionC.transactionService.CreateNewTransaction(ctx, transactionDTO)
	if err != nil {
		resp := common.CreateFailResponse("failed to process transaction make request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	for _, spot := range spots {
		spot.TransactionID = &newTransaction.ID

		_, err = transactionC.spotService.UpdateSpot(ctx, spot)
		if err != nil {
			resp := common.CreateFailResponse("failed to reserve spot "+spot.Row+strconv.Itoa(spot.Number), http.StatusBadRequest)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		}
	}

	transaction, err := transactionC.transactionService.GetTransactionByID(ctx, newTransaction.ID)
	if err != nil {
		resp := common.CreateFailResponse("failed to process transaction make request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	resp := common.CreateSuccessResponse("successfully created transaction", http.StatusCreated, transaction)
	ctx.JSON(http.StatusCreated, resp)
}

func (transactionC *transactionController) GetAllTransactions(ctx *gin.Context) {
	transactions, err := transactionC.transactionService.GetAllTransactions(ctx)
	if err != nil {
		resp := common.CreateFailResponse("failed to fetch all transactions", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	var resp common.Response
	if len(transactions) == 0 {
		resp = common.CreateSuccessResponse("no transaction found", http.StatusOK, transactions)
	} else {
		resp = common.CreateSuccessResponse("successfully fetched all transactions", http.StatusOK, transactions)
	}
	ctx.JSON(http.StatusOK, resp)
}

func (transactionC *transactionController) GetTransactionsByUsername(ctx *gin.Context) {
	username := ctx.Param("username")

	user, err := transactionC.userService.GetUserByIdentifier(ctx, username)
	if err != nil {
		resp := common.CreateFailResponse("failed to identify user", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	var resp common.Response
	if len(user.Transactions) == 0 {
		resp = common.CreateSuccessResponse("no transaction found", http.StatusOK, user.Transactions)
	} else {
		resp = common.CreateSuccessResponse("successfully fetched user transactions", http.StatusOK, user.Transactions)
	}
	ctx.JSON(http.StatusOK, resp)
}

func (transactionC *transactionController) GetMyTransactions(ctx *gin.Context) {
	userID := ctx.GetUint64("ID")

	transactions, err := transactionC.transactionService.GetTransactionsByUserID(ctx, userID)
	if err != nil {
		resp := common.CreateFailResponse("failed to fetch my transactions", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	var resp common.Response
	if len(transactions) == 0 {
		resp = common.CreateSuccessResponse("no transaction found", http.StatusOK, transactions)
	} else {
		resp = common.CreateSuccessResponse("successfully fetched my transactions", http.StatusOK, transactions)
	}
	ctx.JSON(http.StatusOK, resp)
}

func (transactionC *transactionController) DeleteTransactionByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		resp := common.CreateFailResponse("failed to process id", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	transaction, err := transactionC.transactionService.GetTransactionByID(ctx, id)
	if err != nil {
		resp := common.CreateFailResponse("failed to get transaction", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	if reflect.DeepEqual(transaction, entity.Transaction{}) {
		resp := common.CreateFailResponse("transaction with given id not found", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	err = transactionC.transactionService.DeleteTransactionByID(ctx, id)
	if err != nil {
		resp := common.CreateFailResponse("failed to delete transaction", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	resp := common.CreateSuccessResponse("successfully deleted transaction", http.StatusOK, nil)
	ctx.JSON(http.StatusOK, resp)
}
