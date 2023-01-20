package domain

import "github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/db/yugabyte"

// Buyer, represents a buyer entity
type Buyer struct {
	yugabyte.Model
	Username string  `json:"username"`
	Orders   []Order `json:"orders,omitempty"`
}
