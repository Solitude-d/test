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
	"test/webook/internal/web/jwt"
	"test/webook/pkg/logger"
)

func Test_articleService_Publish(t *testing.T) {
	TestCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.ArticleService
		reqBody  string
		wantCode int
		wantRes  Result
	}{
		{
			name: "新建并发表",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "我新建的标题",
					Content: "我新建的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				return svc
			},
			reqBody: `
{
"title":"我新建的标题",
"content":"我新建的内容"
}
`,
			wantCode: 200,
			wantRes: Result{
				Data: float64(1),
				Msg:  "OK",
			},
		},
		{
			name: "publish失败",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "我新建的标题",
					Content: "我新建的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(0), errors.New("publish失败"))
				return svc
			},
			reqBody: `
{
"title":"我新建的标题",
"content":"我新建的内容"
}
`,
			wantCode: 200,
			wantRes: Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
	}
	for _, tc := range TestCases {
		t.Run(t.Name(), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := gin.Default()
			svc.Use(func(ctx *gin.Context) {
				ctx.Set("claims", &jwt.UserClaim{
					Uid: 123,
				})
			})
			h := NewArticleHandler(tc.mock(ctrl), &logger.NopLogger{})
			h.RegisterRoutes(svc)

			req, err := http.NewRequest(http.MethodPost, "/articles/publish",
				bytes.NewBuffer([]byte(tc.reqBody)))

			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")
			//拿到响应
			resp := httptest.NewRecorder()

			svc.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			if resp.Code != 200 {
				return
			}
			var webRes Result
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			assert.Equal(t, tc.wantRes, webRes)
		})
	}
}
