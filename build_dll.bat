@echo off
go build -ldflags "-s -w" -buildmode c-shared -o mbrtu.dll main.go & pause
