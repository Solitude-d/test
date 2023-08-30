package web

import (
	"net/http"
	"time"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"test/webook/internal/domain"
	"test/webook/internal/service"
)

const biz = "login"

type UserHandler struct {
	svc              service.UserService
	codeSvc          service.CodeService
	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
	usernameRegexExp *regexp.Regexp
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService) *UserHandler {
	const (
		emailRegexPattern    = `^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z0-9]{2,6}$`
		passwordRegexPattern = `^(?=.*\d)(?=.*[a-zA-Z])(?=.*[^\da-zA-Z\s]).{6,12}$`
		usernameRegexPattern = `^[a-zA-Z0-9_-]{4,10}$`
	)
	return &UserHandler{
		svc:              svc,
		codeSvc:          codeSvc,
		emailRegexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		usernameRegexExp: regexp.MustCompile(usernameRegexPattern, regexp.None),
	}
}

func (u *UserHandler) UserRouteRegister(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	ug.GET("/profile", u.ProfileJWT)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/edit", u.Edit)
	ug.POST("/login_sms/code/send", u.SendLoginSMSCode)
	ug.POST("/login_sms", u.LoginSMS)
}

func (u *UserHandler) LoginSMS(c *gin.Context) {
	type Req struct {
		Phone     string `json:"phone"`
		InputCode string `json:"inputCode"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return
	}

	ok, err := u.codeSvc.Verify(c, biz, req.Phone, req.InputCode)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统异常",
		})
		return
	}
	if !ok {
		c.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码有误",
		})
		return
	}
	user, err := u.svc.FindOrCreate(c, req.Phone)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统异常",
		})
		return
	}
	if err = u.setJWTToken(c, user.Id); err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	c.JSON(http.StatusOK, Result{
		Msg: "椒盐成功辣！",
	})
}

func (u *UserHandler) SendLoginSMSCode(c *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return
	}

	err := u.codeSvc.Send(c, biz, req.Phone)
	switch err {
	case nil:
		c.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case service.ErrCodeSendTooMany:
		c.JSON(http.StatusOK, Result{
			Msg: "发送太频繁",
		})
	default:
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}

	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统异常",
		})
		return
	}
	c.JSON(http.StatusOK, Result{
		Msg: "发送成功",
	})
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

func (u *UserHandler) ProfileJWT(c *gin.Context) {
	cc, ok := c.Get("claim")
	if !ok {
		c.String(http.StatusOK, "系统异常")
		return
	}
	claims, ok := cc.(*UserClaim)
	if !ok {
		c.String(http.StatusOK, "系统异常")
		return
	}
	user, err := u.svc.Profile(c, claims.Uid)
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
	sess.Options(sessions.Options{
		//Secure: true,
		//HttpOnly: true,
		MaxAge: 30 * 60,
	})
	sess.Set("userId", user.Id)
	sess.Save()
	c.String(http.StatusOK, "道爷上线辣")
	return
}

func (u *UserHandler) LoginJWT(c *gin.Context) {
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
	if err = u.setJWTToken(c, user.Id); err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
	c.String(http.StatusOK, "道爷上线辣")
}

func (u *UserHandler) setJWTToken(c *gin.Context, id int64) error {
	claims := UserClaim{
		Uid:       id,
		UserAgent: c.Request.UserAgent(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenStr, err := token.SignedString([]byte("xHd&^OrleeXM@Yq40gfww%8S%eND1*md"))
	if err != nil {
		return err
	}
	c.Header("x-jwt-token", tokenStr)
	return nil
}

func (u *UserHandler) LogOut(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Options(sessions.Options{
		MaxAge: -1,
	})
	sess.Save()
	c.String(http.StatusOK, "道爷走辣")
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

type UserClaim struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}
