
@echo off
set myself=%~dp0
set myroot=%myself:~,1%
cd %myself%..
%myroot%:
set parent=%cd%
set parentdir=%cd%\
set piddir=%parent%\pid
set pidpath=%piddir%\slave.pid
set num=0
set exitcode=0
@echo on

@echo stop slave
@echo off

if not exist "%pidpath%" (goto :exit)
setlocal enabledelayedexpansion
for /f "tokens=1,2 delims==" %%a in (%pidpath%) do (
    if not "%%a" == "" (
        taskkill /pid %%a
        if "%ERRORLEVEL%"=="0" (
            set /a num=!num!+1
            @echo on
            @echo slave !num! stop, pid:%%a
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