package web

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"test/webook/internal/domain"
	"test/webook/internal/service"
	svcmocks "test/webook/internal/service/mocks"
)

func TestUserHandler(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "12345@qq.com",
					Password: "123456",
				}).Return(nil)
				return userSvc
			},
			reqBody: `{
					"email":"12345@qq.com",
					"password":"123456",
					"confirmPassword":"123456"}`,
			wantCode: http.StatusOK,
			wantBody: "道爷我成辣！！"},
		{
			name: "参数错误",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody: `{
					"email1":"12345@qq.com",
					"password":"123456"`,
			wantCode: http.StatusBadRequest},
		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody: `{
					"email":"12345@qq.c",
					"password":"123456"}`,
			wantCode: http.StatusOK,
			wantBody: "邮箱格式不对"},
		{
			name: "两次密码不一致",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody: `{
					"email":"12345@qq.com",
					"password":"123456",
					"confirmPassword":"12356"}`,
			wantCode: http.StatusOK,
			wantBody: "两次密码输入不一致"},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "12345@qq.com",
					Password: "123456",
				}).Return(service.ErrUserDuplicateEmail)
				return userSvc
			},
			reqBody: `{
					"email":"12345@qq.com",
					"password":"123456",
					"confirmPassword":"123456"}`,
			wantCode: http.StatusOK,
			wantBody: "邮箱冲突"},
		{
			name: "系统异常",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "12345@qq.com",
					Password: "123456",
				}).Return(errors.New("这里是系统异常"))
				return userSvc
			},
			reqBody: `{
					"email":"12345@qq.com",
					"password":"123456",
					"confirmPassword":"123456"}`,
			wantCode: http.StatusOK,
			wantBody: "系统异常"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := gin.Default()
			h := NewUserHandler(tc.mock(ctrl), nil)
			h.UserRouteRegister(svc)

			req, err := http.NewRequest(http.MethodPost, "/users/signup",
				bytes.NewBuffer([]byte(tc.reqBody)))

			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")
			//拿到响应
			resp := httptest.NewRecorder()

			svc.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())
		})
	}
}
