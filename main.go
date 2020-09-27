package main

import (
	"ecret-api/metrics"
	"ecret-api/options"
	"ecret-api/secret_api"
	"fmt"
	"net/http"
)

// about
var (
	GitCommit string
	GitBranch string
	BuildDate string
	Version   string
)

func main() {
	// инициализация приложения
	application.LogAbout(GitCommit, GitBranch, BuildDate, Version)

	// загрузка опций
	opt := options.Load("mt-secret-api")
	application.InitLog(&opt.Log)

	// инициализация сервиса
	metrics.Init(opt.ApplicationName)
	secret_api.Init(&opt.SecretAPI)
	grpc.Init(&opt.GRPC)
	http.Init(&opt.HTTP, fmt.Sprintf("%s %s/%s/%s/%s", opt.ApplicationName, GitCommit, GitBranch, BuildDate, Version))

	// запуск обработчиков
	go grpc.Start()
	go http.Start()

	// ожидание завершения приложения
	application.Wait()

	// финализация сервиса
	http.Stop()
	grpc.Stop()
	application.LogDone()
}
