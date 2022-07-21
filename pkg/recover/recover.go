package recover

import (
	"blog/pkg/l"
	"fmt"
	"github.com/gorilla/handlers"
	"go.uber.org/zap"
	"net/http"
)

type recoverLogger struct {
	logger *zap.Logger
}

func (r *recoverLogger) Println(fields ...interface{}) {
	r.logger.Error(fmt.Sprintln(fields))
}

func NewRecoverHandler(printStack bool) func(h http.Handler) http.Handler {
	l := recoverLogger{l.GetLogger()}
	return handlers.RecoveryHandler(handlers.RecoveryLogger(&l), handlers.PrintRecoveryStack(printStack))
}
