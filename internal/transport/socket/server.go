package socket

import (
	"net/http"
	"rttask/internal/domain/service/task"

	"github.com/doquangtan/socketio/v4"
	"go.uber.org/zap"
)

type SocketServer struct {
	io     *socketio.Io
	logger *zap.Logger

	taskService *task.TaskService
}

func NewSocketServer(taskService *task.TaskService, logger *zap.Logger) *SocketServer {
	io := socketio.New()

	server := &SocketServer{
		io:          io,
		logger:      logger,
		taskService: taskService,
	}

	return server
}

func (s *SocketServer) HttpHandler() http.Handler {
	return s.io.HttpHandler()
}
