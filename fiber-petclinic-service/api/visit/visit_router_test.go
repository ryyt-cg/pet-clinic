package visit

import (
	"encoding/json"
	"fiber-petclinic-service/pkg/test"
	"go.uber.org/zap"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

// config the gin engine for testing purpose
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func Test_VisitById(t *testing.T) {
	//logger := log.New().With(nil, "function", "Test_VisitById")
	logger := zap.NewNop()
	visitMock := MockServicer{}

	visitResponse := &Response{
		ID:          1,
		Description: "spayed",
		VisitDate:   "2010/09/07",
	}

	visitMock.On("getVisitById", 1).Return(visitResponse, nil)
	visitRouter := NewRouter(logger, &visitMock)

	r := setupRouter()
	v1 := r.Group("/v1")
	visitRouter.Register(v1.Group("/visits"))

	tc1, _ := json.Marshal(visitResponse)

	tests := []test.APITestCase{
		{"Get Visit By ID", "GET", "/v1/visits/1", "", nil, http.StatusOK, string(tc1)},
	}
	for _, tc := range tests {
		test.Endpoint(t, r, tc)
	}
}
