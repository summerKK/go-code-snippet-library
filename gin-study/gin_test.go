package gin_test

import (
	"context"
	"testing"
	"time"

	"github.com/summerKK/go-code-snippet-library/gin-study"
)

func TestEngine_Run(t *testing.T) {
	engine := gin.Default()
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	engine.Run(ctx, ":8080")
}
