package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"test/webook/internal/domain"
	"test/webook/internal/integration/startup"
	"test/webook/internal/repository/dao/article"
	ijwt "test/webook/internal/web/jwt"
)

// ArticleTestSuite 单元测试套件
type ArticleTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

// SetupSuite  在所有测试之前初始化一些内容
func (s *ArticleTestSuite) SetupSuite() {
	//s.server = startup.InitWebServer()

	s.server = gin.Default()
	s.server.Use(func(ctx *gin.Context) {
		ctx.Set("claims", &ijwt.UserClaim{
			Uid: 217,
		})
	})
	s.db = startup.InitTestDB()
	artHdl := startup.InitArticleHandler(article.NewGORMArticleDAO(s.db))
	artHdl.RegisterRoutes(s.server)
}

// TearDownTest 每一个测试用例都会执行
func (s *ArticleTestSuite) TearDownTest() {
	//清空所有数据 并且自增主键恢复到1
	s.db.Exec("TRUNCATE TABLE articles")
	s.db.Exec("TRUNCATE TABLE publish_articles")
}

func (s *ArticleTestSuite) TestEdit() {
	t := s.T()
	testCases := []struct {
		name string

		//准备数据
		before func(t *testing.T)
		//验证数据
		after func(t *testing.T)

		art Article

		wantCode int
		// 新建之后返回帖子的id
		wantRes Result[int64]
	}{
		{
			name: "新建帖子，保存成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				var art article.Article
				err := s.db.Where("id = ?", 1).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				art.Ctime = 0
				art.Utime = 0

				assert.Equal(t, article.Article{
					Id:       1,
					Title:    "新建1",
					Content:  "内容1",
					AuthorId: 217,
					Status:   domain.ArticleStatusUnpublished.ToUint8(),
				}, art)
			},

			art: Article{
				Title:   "新建1",
				Content: "内容1",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 1,
				Msg:  "OK",
			},
		},
		{
			name: "修改已有帖子，并保存",
			before: func(t *testing.T) {
				err := s.db.Create(article.Article{
					Id:       2,
					Title:    "我的标题",
					Content:  "内容",
					AuthorId: 217,
					Ctime:    123,
					Utime:    234,
					Status:   domain.ArticleStatusPublished.ToUint8(),
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var art article.Article
				err := s.db.Where("id = ?", 2).First(&art).Error
				assert.NoError(t, err)
				// 确保更新了 更新时间(Utime)
				assert.True(t, art.Utime > 234)
				art.Utime = 0

				assert.Equal(t, article.Article{
					Id:       2,
					Title:    "更新后的标题",
					Content:  "更新后的内容",
					Ctime:    123,
					AuthorId: 217,
					Status:   domain.ArticleStatusUnpublished.ToUint8(),
				}, art)
			},

			art: Article{
				Id:      2,
				Title:   "更新后的标题",
				Content: "更新后的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 2,
				Msg:  "OK",
			},
		},
		{
			name: "偷偷修改别人的帖子",
			before: func(t *testing.T) {
				err := s.db.Create(article.Article{
					Id:      3,
					Title:   "我的标题",
					Content: "我的内容",
					// 模拟用户上217  这里是123 说明我再偷偷改别人的数据
					AuthorId: 123,
					Ctime:    123,
					Utime:    234,
					Status:   domain.ArticleStatusPublished.ToUint8(),
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var art article.Article
				err := s.db.Where("id = ?", 3).First(&art).Error
				assert.NoError(t, err)

				assert.Equal(t, article.Article{
					Id:       3,
					Title:    "我的标题",
					Content:  "我的内容",
					Ctime:    123,
					Utime:    234,
					AuthorId: 123,
					Status:   domain.ArticleStatusUnpublished.ToUint8(),
				}, art)
			},

			art: Article{
				Id:      3,
				Title:   "更新后的标题",
				Content: "更新后的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Code: 5,
				Data: 0,
				Msg:  "系统错误",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)

			reqBody, err := json.Marshal(tc.art)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost,
				"/articles/edit", bytes.NewReader(reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()

			s.server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)

			if resp.Code != 200 {
				return
			}
			var webRes Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			assert.Equal(t, tc.wantRes, webRes)
			tc.after(t)
		})
	}
}

func (s *ArticleTestSuite) TestABC() {
	s.T().Log("这是测试套件")
}

func TestArticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuite{})
}

type Article struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}
