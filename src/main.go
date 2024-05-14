package main

import (
	"fmt"
	"github.com/elliotxx/healthcheck"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
	"golang.org/x/sys/unix"
	"log/slog"
	"os"
	"strconv"
	"syscall"
)

const (
	VolumePath1 = "VOLUME_PATH1"
	VolumePath2 = "VOLUME_PATH2"
	FileNumber  = "FILE_NUM"
	FileSize    = "FILE_SIZE"
	Role        = "ROLE"
)

func server() {
	r := gin.Default()

	r.GET("livez", healthcheck.NewHandler(healthcheck.NewDefaultHandlerConfig()))
	r.GET("readyz", healthcheck.NewHandler(healthcheck.NewDefaultHandlerConfigFor()))

	err := r.Run("0.0.0.0:8080")
	if err != nil {
		slog.Error("Failed to start server", "err", err)
	}
	slog.Info("Server started")
}

func scheduler(fn func(), crontab string) {
	s, err := gocron.NewScheduler()
	if err != nil {
		slog.Error("Failed to create scheduler", "err", err)
	}

	job, err := s.NewJob(
		gocron.CronJob(crontab, true),
		gocron.NewTask(fn),
	)
	if err != nil {
		slog.Error("Failed to create job", "err", err)
	}
	slog.Info("Job started", "id", job.ID())
	s.Start()
}

func fileError(err error, fileName string) {
	if err != nil {
		slog.Error("File error", "fileName", fileName, "err", err)
	}
}

func createFiles(num int64, size int64, path string) {
	data := make([]rune, size)
	one := []rune("1")
	for i := range data {
		data[i] = one[0]
	}

	for i := int64(0); i < num; i++ {
		fileName := fmt.Sprintf("%s/data%d.txt", path, i)

		fd, err := syscall.Open(fileName, syscall.O_CREAT|syscall.O_RDWR, 0644)
		fileError(err, fileName)

		_, err = syscall.Write(fd, []byte(string(data)))
		fileError(err, fileName)

		_, _, errLSeek := unix.Syscall(unix.SYS_LSEEK, uintptr(fd), 0, unix.SEEK_SET)
		if errLSeek != 0 {
			slog.Error("LSeek error", "fileName", fileName, "err", err)
		}

		_, _, errLSeek = unix.Syscall(unix.SYS_LSEEK, uintptr(fd), 0, unix.SEEK_END)
		if errLSeek != 0 {
			slog.Error("LSeek error", "fileName", fileName, "err", err)
		}

		_, _, errSyncFs := unix.Syscall(unix.SYS_SYNCFS, uintptr(fd), 0, 0)
		if errSyncFs != 0 {
			slog.Error("SyncFs error", "fileName", fileName, "err", err)
		}

		err = syscall.Close(fd)
		fileError(err, fileName)
	}
}

func producer() {
	volumePath1 := os.Getenv(VolumePath1)
	slog.Info("Volume path", VolumePath1, volumePath1)
	volumePath2 := os.Getenv(VolumePath2)
	slog.Info("Volume path", VolumePath2, volumePath2)
	fileNumberEnv := os.Getenv(FileNumber)
	fileNumber, err := strconv.ParseInt(fileNumberEnv, 0, 0)
	if err != nil {
		slog.Error("Failed to parse file number", "err", err)
		os.Exit(2)
	}
	slog.Info("FileNumber path", FileNumber, fileNumber)
	fileSizeEnv := os.Getenv(FileSize)
	fileSize, err := strconv.ParseInt(fileSizeEnv, 0, 0)
	if err != nil {
		slog.Error("Failed to parse file size", "err", err)
		os.Exit(2)
	}
	slog.Info("FileSize path", FileSize, fileSize)

	scheduler(
		func() {
			slog.Info("Producing")
			createFiles(fileNumber, fileSize, volumePath1)
			createFiles(fileNumber, fileSize, volumePath2)
		},
		"*/10 * * * * *")

	server()
}

func consumeVolume(path string) {
	slog.Info("Consuming", "path", path)
	entries, err := os.ReadDir(path)
	if err != nil {
		slog.Error("Failed to read volume", "err", err)
		os.Exit(2)
	}

	for _, entry := range entries {
		fileName := fmt.Sprintf("%s/%s", path, entry.Name())
		slog.Info(fileName)

		stat := &syscall.Stat_t{}
		err := syscall.Stat(fileName, stat)
		fileError(err, fileName)

		fd, err := syscall.Open(fileName, syscall.O_RDONLY, 0644)
		fileError(err, fileName)

		data := make([]byte, stat.Size)

		_, err = syscall.Read(fd, data)
		fileError(err, fileName)

		_, _, errLSeek := unix.Syscall(unix.SYS_LSEEK, uintptr(fd), 0, unix.SEEK_SET)
		if errLSeek != 0 {
			slog.Error("LSeek error", "fileName", fileName, "err", err)
		}

		_, _, errLSeek = unix.Syscall(unix.SYS_LSEEK, uintptr(fd), 0, unix.SEEK_END)
		if errLSeek != 0 {
			slog.Error("LSeek error", "fileName", fileName, "err", err)
		}

		_, _, errSyncFs := unix.Syscall(unix.SYS_SYNCFS, uintptr(fd), 0, 0)
		if errSyncFs != 0 {
			slog.Error("SyncFs error", "fileName", fileName, "err", err)
		}

		if stat.Size > 1 {
			fileValue := data[stat.Size-1]
			slog.Info("Value from file", "fileValue", string(fileValue))
		}

		err = syscall.Close(fd)
		fileError(err, fileName)
	}
}

func consumer() {
	volumePath1 := os.Getenv(VolumePath1)
	slog.Info("Volume path", VolumePath1, volumePath1)
	volumePath2 := os.Getenv(VolumePath2)
	slog.Info("Volume path", VolumePath2, volumePath2)

	scheduler(func() {
		slog.Info("Consuming")
		consumeVolume(volumePath1)
		consumeVolume(volumePath2)
	}, "*/1 * * * * *")

	server()
}

func main() {
	applicationRole := os.Getenv(Role)
	slog.Info("Application role", Role, applicationRole)

	if applicationRole == "producer" {
		producer()
	} else if applicationRole == "consumer" {
		consumer()
	} else {
		slog.Error("Application role not recognized", "applicationRole", applicationRole)
		os.Exit(1)
	}
}
