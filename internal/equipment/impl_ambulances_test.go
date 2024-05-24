package equipment

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/xmichalikl/equipment-webapi/internal/db_service"
)

type AmbulanceSuite struct {
	suite.Suite
	dbServiceMock *DbServiceMock[Ambulance]
}

func TestEquipmentSuite(t *testing.T) {
	suite.Run(t, new(AmbulanceSuite))
}

type DbServiceMock[DocType interface{}] struct {
	mock.Mock
}

func (this *DbServiceMock[DocType]) CreateDocument(ctx context.Context, id string, document *DocType) error {
	args := this.Called(ctx, id, document)
	return args.Error(0)
}

func (this *DbServiceMock[DocType]) FindDocument(ctx context.Context, id string) (*DocType, error) {
	args := this.Called(ctx, id)
	return args.Get(0).(*DocType), args.Error(1)
}

func (this *DbServiceMock[DocType]) UpdateDocument(ctx context.Context, id string, document *DocType) error {
	args := this.Called(ctx, id, document)
	return args.Error(0)
}

func (this *DbServiceMock[DocType]) DeleteDocument(ctx context.Context, id string) error {
	args := this.Called(ctx, id)
	return args.Error(0)
}

func (this *DbServiceMock[DocType]) GetAllDocuments(ctx context.Context) ([]DocType, error) {
	args := this.Called(ctx)
	return args.Get(0).([]DocType), args.Error(1)
}

func (this *DbServiceMock[DocType]) Disconnect(ctx context.Context) error {
	args := this.Called(ctx)
	return args.Error(0)
}

func (suite *AmbulanceSuite) SetupTest() {
	suite.dbServiceMock = &DbServiceMock[Ambulance]{}

	// Compile time Assert that the mock is of type db_service.DbService[Ambulance]
	var _ db_service.DbService[Ambulance] = suite.dbServiceMock
	// lastInspectionDate, _ := time.Parse(time.RFC3339, "2024-05-22T08:17:23.950Z")

	suite.dbServiceMock.
		On("FindDocument", mock.Anything, mock.Anything).
		Return(
			&Ambulance{
				Id:   "test-ambulance",
				Name: "test-name",
				Equipment: []Equipment{
					{
						Id:                 "test-equipment",
						Name:               "test-name",
						Availability:       "available",
						LastInspectionDate: time.Now(),
						TechnicalCondition: 3,
						InspectionInterval: 6,
					},
				},
			},
			nil,
		)
}

func (suite *AmbulanceSuite) Test_UpdateEquipment_DbServiceUpdateCalled() {
	// ARRANGE
	suite.dbServiceMock.
		On("UpdateDocument", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	json := `{
		"id": "test-equipment",
		"name": "test-name",
		"availability": "available",
		"lastInspectionDate": "2024-05-22T08:17:23.950Z",
		"technicalCondition": 3,
		"inspectionInterval": 6
	}`

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", suite.dbServiceMock)
	ctx.Params = []gin.Param{
		{Key: "ambulanceId", Value: "test-ambulance"},
		{Key: "equipmentId", Value: "test-equipment"},
	}
	ctx.Request = httptest.NewRequest("POST", "/ambulances/test-ambulance/equipment/test-equipment", strings.NewReader(json))

	sut := implEquipmentAPI{}

	// ACT
	sut.UpdateEquipmentById(ctx)

	// ASSERT
	suite.dbServiceMock.AssertCalled(suite.T(), "UpdateDocument", mock.Anything, "test-ambulance", mock.Anything)
}
