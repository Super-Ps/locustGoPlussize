
@echo off
REM Run as admin
fltmc>nul||cd/d %~dp0&&mshta vbscript:CreateObject("Shell.Application").ShellExecute("%~nx0","%1","","runas",1)(window.close)&&exit

REM setvar
set myself=%~dp0
set myroot=%myself:~,1%
cd %myself%..
%myroot%:
set parent=%cd%
set parentdir=%cd%\
set stopname=stop.bat
set piddir=%parent%\pid
set appdir=%parent%\app
set pidpath=%piddir%\httpmonitor.pid
set num=0
set exitcode=0
cd %appdir%
@echo on

@echo stop server
@echo off

if not exist "%pidpath%" (goto :exit)
setlocal enabledelayedexpansion
for /f "tokens=1,2 delims==" %%a in (%pidpath%) do (
    if not "%%a" == "" (
        taskkill /f /pid %%a
        if "%ERRORLEVEL%"=="0" (
            set /a num=!num!+1
            @echo on
            @echo server !num! stop, pid:%%a
            @echo off
        )
    )
)
del /f /q %pidpath%

:exit
@echo on
@echo stop over
@echo off
if "%1%"=="-nopause" (
    exit /b %exitcode%
)
pause