package v1

import (
	"log/slog"
	"net/http"

	"github.com/Homyakadze14/PsyhoApp/ApiGatewate/internal/common"
	"github.com/Homyakadze14/PsyhoApp/ApiGatewate/internal/entities"
	authv1 "github.com/Homyakadze14/PsyhoApp/ApiGatewate/proto/gen/auth"
	"github.com/gin-gonic/gin"
)

type authRoutes struct {
	s   authv1.AuthClient
	log *slog.Logger
}

func NewAuthRoutes(log *slog.Logger, handler *gin.RouterGroup, s authv1.AuthClient) {
	r := &authRoutes{
		log: log,
		s:   s,
	}

	g := handler.Group("/auth")
	{
		g.POST("/register", r.register)
		g.POST("/login", r.login)
		g.POST("/logout", r.logout)
		g.POST("/activate_account", r.activateAccount)
		g.POST("/refresh", r.refresh)
		g.POST("/send_password_link", r.sndPwdLink)
		g.POST("/change_password", r.changePwd)
	}
}

// @Summary     Register
// @Description Register
// @ID          Register
// @Tags  	    Auth
// @Accept      json
// @Param 		register body entities.RegisterRequest false "register"
// @Produce     json
// @Success     200 {object} authv1.RegisterResponse
// @Failure     400
// @Failure     404
// @Failure     500
// @Failure     503
// @Router      /auth/register [post]
func (r *authRoutes) register(c *gin.Context) {
	const op = "authRoutes.register"

	log := r.log.With(
		slog.String("op", op),
	)

	var req *entities.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	resp, err := r.s.Register(c.Request.Context(), req.ToGRPC())
	if err != nil {
		code, err := common.GetProtoErrWithStatusCode(err)
		log.Error(err.Error())
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary     Login
// @Description Login
// @ID          Login
// @Tags  	    Auth
// @Accept      json
// @Param 		login body entities.LoginRequest false "login"
// @Produce     json
// @Success     200 {object} authv1.LoginResponse
// @Failure     400
// @Failure     401
// @Failure     404
// @Failure     500
// @Failure     503
// @Router      /auth/login [post]
func (r *authRoutes) login(c *gin.Context) {
	const op = "authRoutes.login"

	log := r.log.With(
		slog.String("op", op),
	)

	var req *entities.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	resp, err := r.s.Login(c.Request.Context(), req.ToGRPC())
	if err != nil {
		code, err := common.GetProtoErrWithStatusCode(err)
		log.Error(err.Error())
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary     Logout
// @Description Logout
// @ID          Logout
// @Tags  	    Auth
// @Accept      json
// @Param 		logout body entities.LogoutRequest false "logout"
// @Produce     json
// @Success     200 {object} authv1.LogoutResponse
// @Failure     400
// @Failure     404
// @Failure     500
// @Failure     503
// @Router      /auth/logout [post]
func (r *authRoutes) logout(c *gin.Context) {
	const op = "authRoutes.logout"

	log := r.log.With(
		slog.String("op", op),
	)

	var req *entities.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	resp, err := r.s.Logout(c.Request.Context(), req.ToGRPC())
	if err != nil {
		code, err := common.GetProtoErrWithStatusCode(err)
		log.Error(err.Error())
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary     Activate account
// @Description Activate account
// @ID          Activate account
// @Tags  	    Auth
// @Accept      json
// @Param 		activate body entities.ActivateAccountRequest false "activate"
// @Produce     json
// @Success     200 {object} authv1.ActivateAccountResponse
// @Failure     400
// @Failure     404
// @Failure     500
// @Failure     503
// @Router      /auth/activate_account [post]
func (r *authRoutes) activateAccount(c *gin.Context) {
	const op = "authRoutes.activateAccount"

	log := r.log.With(
		slog.String("op", op),
	)

	var req *entities.ActivateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	resp, err := r.s.ActivateAccount(c.Request.Context(), req.ToGRPC())
	if err != nil {
		code, err := common.GetProtoErrWithStatusCode(err)
		log.Error(err.Error())
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary     Refresh token
// @Description Refresh token
// @ID          Refresh token
// @Tags  	    Auth
// @Accept      json
// @Param 		refresh body entities.RefreshRequest false "refresh"
// @Produce     json
// @Success     200 {object} authv1.RefreshResponse
// @Failure     400
// @Failure     401
// @Failure     404
// @Failure     500
// @Failure     503
// @Router      /auth/refresh [post]
func (r *authRoutes) refresh(c *gin.Context) {
	const op = "authRoutes.activateArefreshccount"

	log := r.log.With(
		slog.String("op", op),
	)

	var req *entities.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	resp, err := r.s.Refresh(c.Request.Context(), req.ToGRPC())
	if err != nil {
		code, err := common.GetProtoErrWithStatusCode(err)
		log.Error(err.Error())
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary     Send password link
// @Description Send password link
// @ID          Send password link
// @Tags  	    Auth
// @Accept      json
// @Param 		email body entities.SendPasswordLinkRequest false "email"
// @Produce     json
// @Success     200 {object} authv1.SendPasswordLinkResponse
// @Failure     400
// @Failure     401
// @Failure     404
// @Failure     500
// @Failure     503
// @Router      /auth/send_password_link [post]
func (r *authRoutes) sndPwdLink(c *gin.Context) {
	const op = "authRoutes.sndPwdLink"

	log := r.log.With(
		slog.String("op", op),
	)

	var req *entities.SendPasswordLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	resp, err := r.s.SendPasswordLink(c.Request.Context(), req.ToGRPC())
	if err != nil {
		code, err := common.GetProtoErrWithStatusCode(err)
		log.Error(err.Error())
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary     Change password
// @Description Change password
// @ID          Change password
// @Tags  	    Auth
// @Accept      json
// @Param 		password body entities.ChangePasswordRequest false "password"
// @Produce     json
// @Success     200 {object} authv1.ChangePasswordResponse
// @Failure     400
// @Failure     401
// @Failure     404
// @Failure     500
// @Failure     503
// @Router      /auth/change_password [post]
func (r *authRoutes) changePwd(c *gin.Context) {
	const op = "authRoutes.changePwd"

	log := r.log.With(
		slog.String("op", op),
	)

	var req *entities.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	resp, err := r.s.ChangePassword(c.Request.Context(), req.ToGRPC())
	if err != nil {
		code, err := common.GetProtoErrWithStatusCode(err)
		log.Error(err.Error())
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
