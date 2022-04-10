package httpserver

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	server       *http.Server
	graceTimeout time.Duration // время на штатное завершения работы сервера
}

type ServerOption func(server *Server)

func New(handler http.Handler, opts ...ServerOption) *Server {
	const (
		defaultAddress      = "127.0.0.1:8080"
		defaultGraceTimeout = 20 * time.Second
	)
	s := &Server{
		server: &http.Server{
			Addr:           defaultAddress,
			ReadTimeout:    time.Second * 10,
			WriteTimeout:   time.Second * 10,
			IdleTimeout:    time.Second * 10, // максимальное время ожидания следующего запроса
			MaxHeaderBytes: 1 << 20,          // 2^20 = 128 Kb
			Handler:        handler,
		},
		graceTimeout: defaultGraceTimeout,
	}

	// Применяем в цикле каждую опцию
	for _, opt := range opts {
		// вызываем функцию, предоставляющую экземпляр *Server в качестве аргумента
		opt(s)
	}

	// вернуть измененный экземпляр Server
	return s
}

func WithAddress(addr string) ServerOption {
	return func(s *Server) {
		if addr != "" {
			s.server.Addr = addr
		}
	}
}

func (s *Server) Serve() {
	// зададим контекст выполнения сервера
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// штатное завершение по сигналам: syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// определяем время для штатного завершения работы сервера
		shutdownCtx, cancel := context.WithTimeout(serverCtx, s.graceTimeout)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Пришёл сигнал завершить работу: штатно завершаем работу сервера не прерывая никаких активных подключений.
		// Завершение работы выполняется в порядке:
		// - закрытия всех открытых подключений;
		// - затем закрытия всех незанятых подключений;
		// - а затем бесконечного ожидания возврата подключений в режим ожидания;
		// - наконец, завершения работы.
		err := s.server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// запустим сервер
	log.Printf("starting HTTP-server at %s\n", s.server.Addr)
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// ожидаем сигнала остановки сервера через context
	<-serverCtx.Done()
}
