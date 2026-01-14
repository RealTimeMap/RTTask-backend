package socket

import (
	"context"

	"rttask/internal/domain/service/task"
	"rttask/internal/infrastructure/security"
	"rttask/internal/transport/dto"

	"github.com/zishang520/socket.io/servers/socket/v3"
	"go.uber.org/zap"
)

// TASK NAMESPACE
// EVENTS
// Client -> Server
// 1. personalTasks:get - принудительное получение задач
//
// Server -> Client
// 1. personalTasks - отправляется при подключении и по запросу personalTasks:get

type TaskNamespace struct {
	nsp         socket.Namespace
	taskService *task.TaskService
	jwtManager  security.JWTManager
	logger      *zap.Logger
}

func NewTaskNamespace(io *socket.Server, taskService *task.TaskService, jwtManager security.JWTManager, logger *zap.Logger) *TaskNamespace {
	nsp := io.Of("/tasks", nil)

	tn := &TaskNamespace{
		nsp:         nsp,
		taskService: taskService,
		jwtManager:  jwtManager,
		logger:      logger,
	}

	tn.setupAuthMiddleware()
	tn.setupHandlers()

	return tn
}

// setupAuthMiddleware мидлвейр для проврки авторизации внутри сокета
func (tn *TaskNamespace) setupAuthMiddleware() {
	tn.nsp.Use(func(client *socket.Socket, next func(*socket.ExtendedError)) {
		tn.logger.Info("Auth middleware for /tasks",
			zap.String("socketID", string(client.Id())),
		)

		claims := tn.authenticateClient(client)
		if claims == nil {
			tn.logger.Warn("Unauthorized connection attempt to /tasks")
			next(socket.NewExtendedError("Unauthorized: valid token required", nil))
			return
		}

		next(nil)
	})
}

func (tn *TaskNamespace) authenticateClient(client *socket.Socket) *security.Claims {
	authMap := client.Handshake().Auth
	if authMap == nil {
		return nil
	}

	token, ok := authMap["token"].(string)
	if !ok || token == "" {
		return nil
	}

	claims, err := tn.jwtManager.ValidateToken(token)
	if err != nil {
		tn.logger.Warn("Invalid token in /tasks namespace",
			zap.String("socketID", string(client.Id())),
			zap.Error(err),
		)
		return nil
	}

	client.SetData(map[string]any{
		"userID": claims.UserID,
		"email":  claims.Email,
	})

	tn.logger.Info("Socket authenticated for /tasks",
		zap.String("socketID", string(client.Id())),
		zap.Uint("userID", claims.UserID),
		zap.String("email", claims.Email),
	)

	return claims
}

func (tn *TaskNamespace) setupHandlers() {
	tn.nsp.On("connection", func(clients ...any) {
		client := clients[0].(*socket.Socket)
		tn.logger.Info("Client connected to /tasks namespace",
			zap.String("socket_id", string(client.Id())),
		)

		// Send personal tasks immediately on connection
		tn.sendPersonalTasks(client)

		// Register task event handlers
		tn.registerGetUserTasksHandler(client)

		client.On("disconnect", func(reason ...any) {
			tn.logger.Info("Client disconnected from /tasks namespace",
				zap.String("socket_id", string(client.Id())),
				zap.Any("reason", reason),
			)
		})
	})
}

// sendPersonalTasks отправляет задачи пользователя при подключении
func (tn *TaskNamespace) sendPersonalTasks(client *socket.Socket) {
	userID := tn.getUserID(client)
	if userID == 0 {
		return
	}

	ctx := context.Background()
	groupedTasks, err := tn.taskService.GetUserTasks(ctx, userID)
	if err != nil {
		tn.logger.Error("Failed to send personal tasks on connection",
			zap.Uint("userID", userID),
			zap.Error(err),
		)
		tn.emitError(client, "personalTasks", err.Error(), 500)
		return
	}

	response := dto.NewGroupedTasksByStatusResponse(groupedTasks)
	client.Emit("personalTasks", response)

	tn.logger.Info("Personal tasks sent on connection",
		zap.Uint("userID", userID),
		zap.Int64("totalCount", response.TotalCount),
	)
}

// getUserID получение ID пользователя из сессии сокета
func (tn *TaskNamespace) getUserID(client *socket.Socket) uint {
	userData := client.Data()
	if userData == nil {
		return 0
	}

	userMap, ok := userData.(map[string]any)
	if !ok {
		return 0
	}

	userID, ok := userMap["userID"].(uint)
	if !ok {
		return 0
	}

	return userID
}

// emitError хэндлер отправки ошибок
func (tn *TaskNamespace) emitError(client *socket.Socket, event string, message string, code int) {
	client.Emit("error", map[string]any{
		"event":   event,
		"message": message,
		"code":    code,
	})
}

// registerGetUserTasksHandler метод регистрации обработки ивента для получения персональных задач
func (tn *TaskNamespace) registerGetUserTasksHandler(client *socket.Socket) {
	client.On("personalTasks:get", func(data ...any) {
		tn.sendPersonalTasks(client)
	})
}

// Namespace returns underlying socket.io namespace for advanced usage
func (tn *TaskNamespace) Namespace() socket.Namespace {
	return tn.nsp
}
