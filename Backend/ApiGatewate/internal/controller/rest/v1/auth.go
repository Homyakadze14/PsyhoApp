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
	s   authv1.AuthServiceClient
	log *slog.Logger
}

func NewAuthRoutes(log *slog.Logger, handler *gin.RouterGroup, s authv1.AuthServiceClient) {
	r := &authRoutes{
		log: log,
		s:   s,
	}

	g := handler.Group("/auth")
	{
		g.POST("/register", r.register)
		g.POST("/login", r.login)
		g.POST("/logout", r.logout)

		// Additional auth endpoints
		g.POST("/generate_auth_code", r.generateAuthCode)
		g.POST("/verify", r.verify)
		g.POST("/generate_service_token", r.generateServiceToken)
		g.POST("/get_role", r.getRole)
		g.POST("/set_role", r.setRole)
		g.POST("/check_access_token", r.checkAccessToken)
		g.POST("/check_service_token", r.checkServiceToken)
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

// @Summary     Generate Auth Code
// @Description Generate authentication code
// @ID          GenerateAuthCode
// @Tags  	    Auth
// @Accept      json
// @Param 		request body authv1.GenerateAuthCodeRequest false "request"
// @Produce     json
// @Success     200 {object} authv1.GenerateAuthCodeResponse
// @Failure     400
// @Failure     404
// @Failure     500
// @Failure     503
// @Router      /auth/generate_auth_code [post]
func (r *authRoutes) generateAuthCode(c *gin.Context) {
	const op = "authRoutes.generateAuthCode"

	log := r.log.With(
		slog.String("op", op),
	)

	var req *authv1.GenerateAuthCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	resp, err := r.s.GenerateAuthCode(c.Request.Context(), req)
	if err != nil {
		code, err := common.GetProtoErrWithStatusCode(err)
		log.Error(err.Error())
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary     Verify
// @Description Verify authentication
// @ID          Verify
// @Tags  	    Auth
// @Accept      json
// @Param 		request body authv1.VerifyRequest false "request"
// @Produce     json
// @Success     200 {object} authv1.VerifyResponse
// @Failure     400
// @Failure     404
// @Failure     500
// @Failure     503
// @Router      /auth/verify [post]
func (r *authRoutes) verify(c *gin.Context) {
	const op = "authRoutes.verify"

	log := r.log.With(
		slog.String("op", op),
	)

	var req *authv1.VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	resp, err := r.s.Verify(c.Request.Context(), req)
	if err != nil {
		code, err := common.GetProtoErrWithStatusCode(err)
		log.Error(err.Error())
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary     Generate Service Token
// @Description Generate service token
// @ID          GenerateServiceToken
// @Tags  	    Auth
// @Accept      json
// @Param 		request body authv1.GenerateServiceTokenRequest false "request"
// @Produce     json
// @Success     200 {object} authv1.GenerateServiceTokenResponse
// @Failure     400
// @Failure     404
// @Failure     500
// @Failure     503
// @Router      /auth/generate_service_token [post]
func (r *authRoutes) generateServiceToken(c *gin.Context) {
	const op = "authRoutes.generateServiceToken"

	log := r.log.With(
		slog.String("op", op),
	)

	var req *authv1.GenerateServiceTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	resp, err := r.s.GenerateServiceToken(c.Request.Context(), req)
	if err != nil {
		code, err := common.GetProtoErrWithStatusCode(err)
		log.Error(err.Error())
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary     Get Role
// @Description Get user role
// @ID          GetRole
// @Tags  	    Auth
// @Accept      json
// @Param 		request body authv1.GetRoleRequest false "request"
// @Produce     json
// @Success     200 {object} authv1.GetRoleResponse
// @Failure     400
// @Failure     404
// @Failure     500
// @Failure     503
// @Router      /auth/get_role [post]
func (r *authRoutes) getRole(c *gin.Context) {
	const op = "authRoutes.getRole"

	log := r.log.With(
		slog.String("op", op),
	)

	var req *authv1.GetRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	resp, err := r.s.GetRole(c.Request.Context(), req)
	if err != nil {
		code, err := common.GetProtoErrWithStatusCode(err)
		log.Error(err.Error())
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary     Set Role
// @Description Set user role
// @ID          SetRole
// @Tags  	    Auth
// @Accept      json
// @Param 		request body authv1.SetRoleRequest false "request"
// @Produce     json
// @Success     200 {object} authv1.SetRoleResponse
// @Failure     400
// @Failure     404
// @Failure     500
// @Failure     503
// @Router      /auth/set_role [post]
func (r *authRoutes) setRole(c *gin.Context) {
	const op = "authRoutes.setRole"

	log := r.log.With(
		slog.String("op", op),
	)

	var req *authv1.SetRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	resp, err := r.s.SetRole(c.Request.Context(), req)
	if err != nil {
		code, err := common.GetProtoErrWithStatusCode(err)
		log.Error(err.Error())
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary     Check Access Token
// @Description Check access token validity
// @ID          CheckAccessToken
// @Tags  	    Auth
// @Accept      json
// @Param 		request body authv1.CheckAccessTokenRequest false "request"
// @Produce     json
// @Success     200 {object} authv1.CheckAccessTokenResponse
// @Failure     400
// @Failure     404
// @Failure     500
// @Failure     503
// @Router      /auth/check_access_token [post]
func (r *authRoutes) checkAccessToken(c *gin.Context) {
	const op = "authRoutes.checkAccessToken"

	log := r.log.With(
		slog.String("op", op),
	)

	var req *authv1.CheckAccessTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	resp, err := r.s.CheckAccessToken(c.Request.Context(), req)
	if err != nil {
		code, err := common.GetProtoErrWithStatusCode(err)
		log.Error(err.Error())
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary     Check Service Token
// @Description Check service token validity
// @ID          CheckServiceToken
// @Tags  	    Auth
// @Accept      json
// @Param 		request body authv1.CheckServiceTokenRequest false "request"
// @Produce     json
// @Success     200 {object} authv1.CheckServiceTokenResponse
// @Failure     400
// @Failure     404
// @Failure     500
// @Failure     503
// @Router      /auth/check_service_token [post]
func (r *authRoutes) checkServiceToken(c *gin.Context) {
	const op = "authRoutes.checkServiceToken"

	log := r.log.With(
		slog.String("op", op),
	)

	var req *authv1.CheckServiceTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	resp, err := r.s.CheckServiceToken(c.Request.Context(), req)
	if err != nil {
		code, err := common.GetProtoErrWithStatusCode(err)
		log.Error(err.Error())
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
