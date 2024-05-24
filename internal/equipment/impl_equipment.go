package equipment

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateEquipment - Create new equipment to an ambulance
func (this *implEquipmentAPI) CreateEquipment(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNotImplemented)
}

// DeleteEquipmentById - Delete specific equipment
func (this *implEquipmentAPI) DeleteEquipmentById(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNotImplemented)
}

// GetEquipmentById - Get specific equipment details
func (this *implEquipmentAPI) GetEquipmentById(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNotImplemented)
}

// UpdateEquipmentById - Update specific equipment details
func (this *implEquipmentAPI) UpdateEquipmentById(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNotImplemented)
}
