package equipment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xmichalikl/equipment-webapi/internal/db_service"
)

// CreateAmbulance - Create new ambulance
func (this *implAmbulancesAPI) CreateAmbulance(ctx *gin.Context) {
	value, exists := ctx.Get("db_service")
	if !exists {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db not found",
				"error":   "db not found",
			})
		return
	}

	db, ok := value.(db_service.DbService[Ambulance])
	if !ok {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db context is not of required type",
				"error":   "cannot cast db context to db_service.DbService",
			})
		return
	}

	ambulance := Ambulance{}
	err := ctx.BindJSON(&ambulance)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  "Bad Request",
				"message": "Invalid request body",
				"error":   err.Error(),
			})
		return
	}

	if ambulance.Name == "" {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  "Bad Request",
				"message": "Name is required",
				"error":   "Name field is empty",
			})
		return
	}

	if ambulance.Equipment == nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  "Bad Request",
				"message": "Equipment is required",
				"error":   "Equipment field is null",
			})
		return
	}

	if ambulance.Id == "" {
		ambulance.Id = uuid.New().String()
	}

	err = db.CreateDocument(ctx, ambulance.Id, &ambulance)

	switch err {
	case nil:
		ctx.JSON(
			http.StatusCreated,
			ambulance,
		)
	case db_service.ErrConflict:
		ctx.JSON(
			http.StatusConflict,
			gin.H{
				"status":  "Conflict",
				"message": "Ambulance already exists",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to create ambulance in database",
				"error":   err.Error(),
			},
		)
	}
}

// DeleteAmbulance - Deletes specific ambulance
func (this *implAmbulancesAPI) DeleteAmbulance(ctx *gin.Context) {
	value, exists := ctx.Get("db_service")
	if !exists {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service not found",
				"error":   "db_service not found",
			})
		return
	}

	db, ok := value.(db_service.DbService[Ambulance])
	if !ok {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service context is not of type db_service.DbService",
				"error":   "cannot cast db_service context to db_service.DbService",
			})
		return
	}

	ambulanceId := ctx.Param("ambulanceId")
	err := db.DeleteDocument(ctx, ambulanceId)

	switch err {
	case nil:
		ctx.AbortWithStatus(http.StatusNoContent)
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Ambulance not found",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to delete ambulance from database",
				"error":   err.Error(),
			})
	}
}

// GetAmbulanceList - Provides the list of ambulances
func (this *implAmbulancesAPI) GetAmbulanceList(ctx *gin.Context) {
	value, exists := ctx.Get("db_service")
	if !exists {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service not found",
				"error":   "db_service not found",
			})
		return
	}

	db, ok := value.(db_service.DbService[Ambulance])
	if !ok {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service context is not of type db_service.DbService",
				"error":   "cannot cast db_service context to db_service.DbService",
			})
		return
	}

	ambulances, err := db.GetAllDocuments(ctx)

	switch err {
	case nil:
		ctx.JSON(
			http.StatusOK,
			ambulances,
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to retrieve ambulances from database",
				"error":   err.Error(),
			},
		)
	}
}
