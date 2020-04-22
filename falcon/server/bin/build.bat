
@echo off
REM Set vars
set myself=%~dp0
set GOPATH=%myself%..\..\..\..
set myroot=%myself:~,1%
cd %myself%..
%myroot%:
set parent=%cd%
set parentdir=%cd%\
set bindir=%parent%\bin
set appdir=%parent%\app
set piddir=%parent%\pid
set apppath=%appdir%\httpmonitor.exe
set gorootpath=
set gobinpath=go
set exitcode=0
cd %appdir%

if not exist "%appdir%" (mkdir %appdir%)
if exist "%apppath%" (del /f /q %apppath%)
if not "%gorootpath%"=="" (
    set GOROOT=%gorootpath%
    set gobinpath=%gorootpath%\bin\go
)
@echo on

@echo build start
call %gobinpath% build -o %apppath% falcon\server
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
