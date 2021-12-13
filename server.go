package todo

import (
	"context"
	"net/http"
	"time"
)

// создаем абстракция над структурой сервер из пакета http
type Server struct {
	httpServer *http.Server
}

// Добавляем метод запуска. Здесь мы инкапсулируем настройки
// сервера в числе которых номер порта и handler. Возвращает функция
// стандартный метод ListenAndServer пакета http , который
// запускает бесконечный цикл и слушает все входящие запросы
// для последующей обработки.
func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

// добавляем метод остановки работы
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
