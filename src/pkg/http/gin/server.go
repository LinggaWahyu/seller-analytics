package gin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/http/domain"
)

// ServeHTTP, runs a gin engine r based on provided config cfg
func ServeHTTP(r *gin.Engine, cfg domain.HTTPServerConfig) {
	err := r.Run(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		panic(err)
	}
}
