package equipment

import (
	"net/http"

	"slices"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateEquipment - Create new equipment to an ambulance
func (this *implEquipmentAPI) CreateEquipment(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		var newEquipment Equipment

		if err := c.ShouldBindJSON(&newEquipment); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		// TODO: Add other checks
		if newEquipment.Name == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Name is required",
			}, http.StatusBadRequest
		}

		if newEquipment.Id == "" || newEquipment.Id == "@new" {
			newEquipment.Id = uuid.NewString()
		}

		conflictIndx := slices.IndexFunc(ambulance.Equipment, func(equipment Equipment) bool {
			return newEquipment.Id == equipment.Id
		})

		if conflictIndx >= 0 {
			return nil, gin.H{
				"status":  http.StatusConflict,
				"message": "Equipment already exists",
			}, http.StatusConflict
		}

		ambulance.Equipment = append(ambulance.Equipment, newEquipment)

		newEquipmentIndx := slices.IndexFunc(ambulance.Equipment, func(equipment Equipment) bool {
			return newEquipment.Id == equipment.Id
		})

		if newEquipmentIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to save entry",
			}, http.StatusInternalServerError
		}
		return ambulance, ambulance.Equipment[newEquipmentIndx], http.StatusOK
	})
}

// DeleteEquipmentById - Delete specific equipment
func (this *implEquipmentAPI) DeleteEquipmentById(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		equipmentId := ctx.Param("equipmentId")

		if equipmentId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Equipment ID is required",
			}, http.StatusBadRequest
		}

		equipmentIndx := slices.IndexFunc(ambulance.Equipment, func(equipment Equipment) bool {
			return equipmentId == equipment.Id
		})

		if equipmentIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Entry not found",
			}, http.StatusNotFound
		}

		ambulance.Equipment = append(ambulance.Equipment[:equipmentIndx], ambulance.Equipment[equipmentIndx+1:]...)
		return ambulance, nil, http.StatusNoContent
	})
}

// GetEquipmentById - Get specific equipment details
func (this *implEquipmentAPI) GetEquipmentById(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		equipmentId := ctx.Param("equipmentId")

		if equipmentId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Equipment ID is required",
			}, http.StatusBadRequest
		}

		equipmentIndx := slices.IndexFunc(ambulance.Equipment, func(equipment Equipment) bool {
			return equipmentId == equipment.Id
		})

		if equipmentIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Equipment not found",
			}, http.StatusNotFound
		}

		// return nil ambulance - no need to update it in db
		return nil, ambulance.Equipment[equipmentIndx], http.StatusOK
	})
}

// GetEquipmentList - Provides the equipment list for a specific ambulance
func (this *implEquipmentAPI) GetEquipmentList(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(
		ctx *gin.Context,
		ambulance *Ambulance,
	) (updatedAmbulance *Ambulance, responseContent interface{}, status int) {
		result := ambulance.Equipment
		if result == nil {
			result = []Equipment{}
		}
		return nil, result, http.StatusOK
	})
}

// UpdateEquipmentById - Update specific equipment details
func (this *implEquipmentAPI) UpdateEquipmentById(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		var equipment Equipment

		if err := c.ShouldBindJSON(&equipment); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		equipmentId := ctx.Param("equipmentId")

		if equipmentId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Equipment ID is required",
			}, http.StatusBadRequest
		}

		equipmentIndx := slices.IndexFunc(ambulance.Equipment, func(equipment Equipment) bool {
			return equipmentId == equipment.Id
		})

		if equipmentIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Equipment not found",
			}, http.StatusNotFound
		}

		if equipment.Id != "" {
			ambulance.Equipment[equipmentIndx].Id = equipment.Id
		}

		if equipment.Name != "" {
			ambulance.Equipment[equipmentIndx].Name = equipment.Name
		}

		if equipment.Availability != "" {
			ambulance.Equipment[equipmentIndx].Availability = equipment.Availability
		}

		if !equipment.LastInspectionDate.IsZero() {
			ambulance.Equipment[equipmentIndx].LastInspectionDate = equipment.LastInspectionDate
		}

		if equipment.TechnicalCondition != 0 {
			ambulance.Equipment[equipmentIndx].TechnicalCondition = equipment.TechnicalCondition
		}

		if equipment.InspectionInterval != 0 {
			ambulance.Equipment[equipmentIndx].InspectionInterval = equipment.InspectionInterval
		}

		return ambulance, ambulance.Equipment[equipmentIndx], http.StatusOK
	})
}
