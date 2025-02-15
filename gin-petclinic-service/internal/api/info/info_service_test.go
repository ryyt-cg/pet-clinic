package info

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func Test_getAppInfo(t *testing.T) {
	//logger := log.New().With(nil, "function", "Test_getAppInfo")
	logger := zap.NewNop()
	infoMock := MockServicer{}
	info := &Info{"Test App", "Info App", "1.0.0", "", ""}
	infoMock.On("getAppInfo").Return(*info, nil)

	infoService := NewService(logger)
	result, _ := infoService.getAppInfo()
	infoMock.AssertExpectations(t)
	infoMock.AssertNumberOfCalls(t, "getAppInfo", 1)

	assert.Equal(t, info.AppName, result.AppName)
	assert.Equal(t, info.Description, result.Description)
	assert.Equal(t, info.Version, result.Version)
}
