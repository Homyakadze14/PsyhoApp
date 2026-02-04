package v1

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/Homyakadze14/PsyhoApp/ApiGatewate/internal/common"
	authv1 "github.com/Homyakadze14/PsyhoApp/ApiGatewate/proto/gen/auth"
	"github.com/gin-gonic/gin"
)

func authMiddleware(log *slog.Logger, s authv1.AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authH := c.GetHeader("Authorization")
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if !strings.Contains(authH, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "bad token"})
			return
		}

		token := strings.Split(authH, "Bearer ")[1]
		slog.Info(token)

		_, err := s.Verify(ctx, &authv1.VerifyRequest{AccessToken: token})
		if err != nil {
			status, err := common.GetProtoErrWithStatusCode(err)
			log.Error(err.Error())
			c.AbortWithStatusJSON(status, gin.H{"error": err.Error()})
			return
		}

		c.Next()
	}
}
