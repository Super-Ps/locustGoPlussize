
@echo off
set gorootpath=
set gobinpath=go
set myself=%~dp0
set appdir=%myself%..\app
set apppath=%appdir%\slave.exe
set GOPATH=%myself%..\..\..\..
set exitcode=0
if not exist "%appdir%" (mkdir "%appdir%")
if exist "%apppath%" (del /f /q "%apppath%")
if not "%gorootpath%"=="" (
    set GOROOT=%gorootpath%
    set gobinpath=%gorootpath%\bin\go
)
@echo on

@echo build start
call %gobinpath% build -o "%apppath%" falcon\slave
@echo off

if "%ERRORLEVEL%"=="0" (
    set exitcode=0
    set bulidresult=build success
    goto :exit
) else (
    set exitcode=1
    set bulidresult=build fail
    goto :exit
)

:exit
@echo on
@echo %bulidresult%
@echo off
if "%1%"=="-nopause" (
    exit /b %exitcode%
)
pause