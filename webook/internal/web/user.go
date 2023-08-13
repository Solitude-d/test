package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"test/webook/internal/domain"
	"test/webook/internal/service"
	"time"
)

type UserHandler struct {
	svc              *service.UserService
	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
	usernameRegexExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern    = `^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z0-9]{2,6}$`
		passwordRegexPattern = `^(?=.*\d)(?=.*[a-zA-Z])(?=.*[^\da-zA-Z\s]).{6,12}$`
		usernameRegexPattern = `^[a-zA-Z0-9_-]{4,10}$`
	)
	return &UserHandler{
		svc:              svc,
		emailRegexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		usernameRegexExp: regexp.MustCompile(usernameRegexPattern, regexp.None),
	}
}

func (u *UserHandler) UserRouteRegister(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	ug.GET("/profile", u.Profile)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
}

func (u *UserHandler) SignUp(c *gin.Context) {
	type sign struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req sign
	if err := c.Bind(&req); err != nil {
		c.String(http.StatusOK, "入参错误")
		return
	}
	isEmail, err := u.emailRegexExp.MatchString(req.Email)
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		c.String(http.StatusOK, "邮箱格式不对")
		return
	}
	err = u.svc.SignUp(c, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrUserDuplicateEmail {
		c.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		c.String(http.StatusOK, "系统异常")
		return
	}
	c.String(http.StatusOK, "道爷我成辣！！")
}

func (u *UserHandler) Profile(c *gin.Context) {
	sess := sessions.Default(c)
	id := sess.Get("userId").(int64)
	user, err := u.svc.Profile(c, id)
	if err != nil {
		c.String(http.StatusOK, "系统异常")
		return
	}
	c.JSON(http.StatusOK, user)
}

func (u *UserHandler) Login(c *gin.Context) {
	type login struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req login

	if err := c.Bind(&req); err != nil {
		c.String(http.StatusOK, "入参错误")
		return
	}

	user, err := u.svc.Login(c, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		c.String(http.StatusOK, "用户名或密码错误")
		return
	}
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	sess := sessions.Default(c)
	sess.Set("userId", user.Id)
	sess.Save()
	c.String(http.StatusOK, "道爷上线辣")
	return
}

func (u *UserHandler) Edit(c *gin.Context) {
	type edit struct {
		ID       int64  `json:"ID"`
		NickName string `json:"NickName"`
		Birth    string `json:"Birth"`
		Synopsis string `json:"Synopsis"`
	}
	var req edit

	if err := c.Bind(&req); err != nil {
		c.String(http.StatusOK, "入参错误")
		return
	}

	checkName, err := u.usernameRegexExp.MatchString(req.NickName)
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	if !checkName {
		c.String(http.StatusOK, "用户名格式不对")
		return
	}
	_, err = time.Parse("2006-01-02", req.Birth)
	if err != nil {
		c.String(http.StatusOK, "生日输入格式错误")
		return
	}
	if len(req.Synopsis) > 30 {
		c.String(http.StatusOK, "个人简介过长")
		return
	}
	err = u.svc.Edit(c, req.ID, req.NickName, req.Birth, req.Synopsis)
	if err != nil {
		c.String(http.StatusOK, "系统错误，修改个人信息失败")
		return
	}
	c.String(http.StatusOK, "个人信息修改成功")
	return
}
