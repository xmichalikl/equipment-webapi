package equipment

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAmbulanceList - Provides the list of ambulances
func (this *implAmbulanceListAPI) GetAmbulanceList(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNotImplemented)
}
