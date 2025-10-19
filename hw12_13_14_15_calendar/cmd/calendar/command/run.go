package command

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/silversbro/v-v-serebryakov/hw12_13_14_15_calendar/internal/app"
	"github.com/silversbro/v-v-serebryakov/hw12_13_14_15_calendar/internal/config"
	"github.com/silversbro/v-v-serebryakov/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/silversbro/v-v-serebryakov/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/silversbro/v-v-serebryakov/hw12_13_14_15_calendar/internal/storage/memory"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the calendar server",
	Run:   run,
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func run(cmd *cobra.Command, args []string) {
	// Загрузка конфигурации
	cfg, err := config.Load(configFile)
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Инициализация логгера
	logg := logger.New(cfg.Logger.Level)

	// Инициализация хранилища (пока только memory)
	storage := memorystorage.New()

	// Создание приложения
	calendar := app.New(logg, storage)

	// Создание HTTP сервера
	server := internalhttp.NewServer(logg, calendar)

	// Настройка graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	// Запуск сервера в горутине
	go func() {
		logg.Info("Starting calendar server...")
		if err := server.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
		}
	}()

	// Ожидание сигнала завершения
	<-ctx.Done()

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	logg.Info("Shutting down calendar server...")

	if err := server.Stop(shutdownCtx); err != nil {
		logg.Error("failed to stop http server: " + err.Error())
	}

	logg.Info("Calendar server stopped")
}
