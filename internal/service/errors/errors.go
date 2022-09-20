// Package errors Пакет предназначен для хранения всех типов ошибок сервиса в одном месте.
package errors

type AppError string

// Error implements error interface
func (e AppError) Error() string {
	return string(e)
}

const (
	MetricNotImplemented   AppError = "metric not implemented"
	ErrEmptyFilestoreName  AppError = "empty filestore name"
	ErrFileStoreNotDefined AppError = "filestore not defined"
)
