package web

import (
	"bytes"
	"encoding/json"
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

func TestUserHandler_LoginSMS(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		reqBody  string
		wantCode int
		wantBody Result
	}{
		{
			name: "椒盐成功辣！",
			reqBody: `{
					"phone":"15381818181",
					"inputCode":"123456"
						}`,
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				uSvc := svcmocks.NewMockUserService(ctrl)
				uSvc.EXPECT().FindOrCreate(gomock.Any(), "15381818181").
					Return(domain.User{
						Phone: "15381818181",
					}, nil)
				cSvc := svcmocks.NewMockCodeService(ctrl)
				cSvc.EXPECT().Verify(gomock.Any(), "login", "15381818181", "123456").
					Return(true, nil)

				return uSvc, cSvc
			},
			wantCode: http.StatusOK,
			wantBody: Result{Msg: "椒盐成功辣！"},
		},
		{
			name: "验证码有误",
			reqBody: `{
					"phone":"15381818181",
					"inputCode":"1234"
						}`,
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				uSvc := svcmocks.NewMockUserService(ctrl)
				cSvc := svcmocks.NewMockCodeService(ctrl)
				cSvc.EXPECT().Verify(gomock.Any(), "login", "15381818181", "1234").
					Return(false, nil)

				return uSvc, cSvc
			},
			wantCode: http.StatusOK,
			wantBody: Result{
				Code: 4,
				Msg:  "验证码有误"},
		},
		{
			name: "设置jwtToken失败，系统错误",
			reqBody: `{
					"phone":"15381818181",
					"inputCode":"123456"
						}`,
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {

				uSvc := svcmocks.NewMockUserService(ctrl)
				uSvc.EXPECT().FindOrCreate(gomock.Any(), "15381818181").
					Return(domain.User{
						Phone: "15381818181",
					}, nil)
				cSvc := svcmocks.NewMockCodeService(ctrl)
				cSvc.EXPECT().Verify(gomock.Any(), "login", "15381818181", "123456").
					Return(true, nil)

				return uSvc, cSvc
			},
			wantCode: http.StatusOK,
			wantBody: Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := gin.Default()
			uSvc, cSvc := tc.mock(ctrl)
			h := NewUserHandler(uSvc, cSvc)
			h.UserRouteRegister(svc)
			req, err := http.NewRequest(http.MethodPost, "/users/login_sms",
				bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			svc.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			var res Result
			err = json.NewDecoder(resp.Body).Decode(&res)
			assert.Equal(t, tc.wantBody, res)
		})
	}
}
