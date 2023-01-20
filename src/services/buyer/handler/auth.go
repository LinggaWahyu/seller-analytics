package handler

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/buyer/domain"
)

// Auth, add the middleware function
func (h *handler) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		buyerIdRaw := session.Get(domain.BuyerKey)
		if buyerIdRaw == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": "no authentication found"})
			return
		}

		buyerId, ok := buyerIdRaw.(uint)
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": "invalid cookie"})
			return
		}
		valid, err := h.BuyerUsecase.IsUserAuthenticated(c.Request.Context(), buyerId)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"errors": err})
			return
		}

		if !valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "request does not have valid authentication"})
			return
		}

		// set buyer key into context
		c.Set(domain.BuyerKey, buyerId)

		c.Next()
	}
}
