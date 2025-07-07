@echo off
REM Script para criar estrutura de diretórios e arquivos de um projeto Go MVC

REM Cria diretórios
md project\cmd\app
md project\config
md project\database
md project\models
md project\repository
md project\service
md project\controllers
md project\routes

REM Cria arquivos em branco
cd project

type nul > .env
type nul > go.mod

type nul > cmd\app\main.go
type nul > config\config.go
type nul > database\connection.go
type nul > models\user.go
type nul > repository\user_repository.go
type nul > service\user_service.go
type nul > controllers\user_controller.go
type nul > routes\routes.go