package errors

var (
	// ErrInvalidCredentials - неверный email или пароль
	// Используем generic сообщение для безопасности (не раскрываем существует ли пользователь)
	ErrInvalidCredentials = NewUnauthorizedError("invalid email or password")

	// ErrUserAlreadyExists - пользователь с таким email уже существует
	ErrUserAlreadyExists = NewAlreadyExistsError("user", "email", "")

	// ErrUserNotFound - пользователь не найден
	ErrUserNotFound = NewNotFoundError("user", "")

	// ErrInvalidToken - невалидный или истекший токен
	ErrInvalidToken = NewUnauthorizedError("invalid or expired token")

	// ErrPasswordTooWeak - пароль не соответствует требованиям
	ErrPasswordTooWeak = NewValidationError("password does not meet requirements")
)
