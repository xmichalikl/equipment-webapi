package equipment

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// // GetEquipmentList - Provides the equipment list for a specific ambulance
func (this *implEquipmentListAPI) GetEquipmentList(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNotImplemented)
}
