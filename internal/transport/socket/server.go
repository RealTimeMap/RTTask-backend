package socket

import (
	"net/http"
	"rttask/internal/domain/service/task"
	"rttask/internal/infrastructure/security"

	"github.com/zishang520/socket.io/servers/socket/v3"
	"go.uber.org/zap"
)

type Server struct {
	io         *socket.Server
	logger     *zap.Logger
	jwtManager security.JWTManager

	// Namespaces
	TaskNamespace *TaskNamespace
}

func NewServer(logger *zap.Logger, jwtManager security.JWTManager, taskService *task.Service) *Server {
	io := socket.NewServer(nil, nil)

	s := &Server{
		io:         io,
		logger:     logger,
		jwtManager: jwtManager,
	}

	// Setup default namespace handlers
	s.setupHandlers()

	// Setup namespaces (auth middleware is inside each namespace)
	s.TaskNamespace = NewTaskNamespace(io, taskService, jwtManager, logger)

	return s
}

func (s *Server) Handler() http.Handler {
	return s.io.ServeHandler(nil)
}

func (s *Server) IO() *socket.Server {
	return s.io
}

func (s *Server) setupHandlers() {
	s.io.On("connection", func(clients ...any) {
		client := clients[0].(*socket.Socket)
		s.logger.Info("Socket.IO client connected",
			zap.String("socket_id", string(client.Id())),
		)

		// Handle disconnect
		client.On("disconnect", func(reason ...any) {
			s.logger.Info("Socket.IO client disconnected",
				zap.String("socket_id", string(client.Id())),
				zap.Any("reason", reason),
			)
		})
	})
}

// emitError sends error event to client
func (s *Server) emitError(client *socket.Socket, event string, message string, code int) {
	client.Emit("error", map[string]any{
		"event":   event,
		"message": message,
		"code":    code,
	})
}
