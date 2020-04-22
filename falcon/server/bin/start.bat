
@echo off
REM Run as admin
fltmc>nul||cd/d %~dp0&&mshta vbscript:CreateObject("Shell.Application").ShellExecute("%~nx0","%1","","runas",1)(window.close)&&exit

set allavg=
:begin
set avgvalue=
set avgvalue=%1%
if "%avgvalue%"=="" (goto :setvar)
if "%avgvalue%"=="-showname" (set showname=%2%)
if "%avgvalue%"=="-serverport" (set serverport=%2%)
if "%avgvalue%"=="-gctime" (set gctime=%2%)
if "%avgvalue%"=="-pausetime" (set pausetime=%2%)
if "%avgvalue%"=="-mode" (set mode=%2%)
if "%avgvalue%"=="-pidfile" (set pidfile=%2%)
if "%avgvalue%"=="--route-root" (set routeroot=%2%)
if "%avgvalue%"=="--route-monitor" (set routemonitor=%2%)
if "%avgvalue%"=="--route-address" (set routeaddress=%2%)
if "%avgvalue%"=="--route-getdemo" (set routegetdemo=%2%)
if "%avgvalue%"=="--route-postdemo" (set routepostdemo=%2%)
if "%avgvalue%"=="--route-restart" (set routerestart=%2%)
if "%avgvalue%"=="--route-quit" (set routequit=%2%)
if "%avgvalue%"=="--route-stop" (set routestop=%2%)
if "%avgvalue%"=="-nopause" (set nopause=true)
if "%avgvalue%"=="-service" (set iservice=true)
if "%avgvalue%"=="-outpath" (set outpath=%2%)
if "%avgvalue%"=="--noout" (set noout=--noout)
set allavg=%allavg% %avgvalue%
SHIFT
goto :begin

:setvar
set myself=%~dp0
set myroot=%myself:~,1%
cd %myself%..
%myroot%:
set parent=%cd%
set parentdir=%cd%\
set stopname=stop.bat
set outname=httpmonitor.out
set bindir=%parent%\bin
set appdir=%parent%\app
set piddir=%parent%\pid
set outdir=%parent%\out
set apppath=%appdir%\httpmonitor.exe
set pidpath=%piddir%\httpmonitor.pid
set stoppath=%bindir%\%stopname%
set exitcode=0
cd %appdir%

if "%showname%"=="" (set showname=%COMPUTERNAME%)
if "%serverport%"=="" (set serverport=5555)
if "%pidfile%"=="" (set pidfile=%pidpath%)
if "%gctime%"=="" (set gctime=60)
if "%pausetime%"=="" (set pausetime=60)
if "%mode%"=="" (set mode=release)
if "%routeroot%"=="" (set routeroot=/)
if "%routemonitor%"=="" (set routemonitor=/monitor)
if "%routeaddress%"=="" (set routeaddress=/address)
if "%routegetdemo%"=="" (set routegetdemo=/getdemo)
if "%routepostdemo%"=="" (set routepostdemo=/postdemo)
if "%routerestart%"=="" (set routerestart=/restart)
if "%routequit%"=="" (set routequit=/quit)
if "%routestop%"=="" (set routestop=/stop)
if "%outpath%"=="" (set outpath=%outdir%/%outname%)

if not exist "%piddir%" (mkdir "%piddir%")
if not exist "%outdir%" (mkdir "%outdir%")
@echo on

call %stoppath% -nopause
@echo start server
if "%iservice%"=="true" (
    call "%apppath%" --port=%serverport% --name="%showname%" --mode=%mode% --pid=" " --gctime=%gctime% --pausetime=%pausetime% ^
    --route-root="%routeroot%" --route-monitor="%routemonitor%" --route-address="%routeaddress%" --route-getdemo="%routegetdemo%" ^
    --route-postdemo="%routepostdemo%" --route-restart="%routerestart%" --route-quit="%routequit%" --route-stop="%routestop%" --noout > "%outpath%"
) else (
    call "%apppath%" --port=%serverport% --name="%showname%" --mode=%mode% --pid=" " --gctime=%gctime% --pausetime=%pausetime% ^
    --route-root="%routeroot%" --route-monitor="%routemonitor%" --route-address="%routeaddress%" --route-getdemo="%routegetdemo%" ^
    --route-postdemo="%routepostdemo%" --route-restart="%routerestart%" --route-quit="%routequit%" --route-stop="%routestop%" %noout%
)

if not "%ERRORLEVEL%"=="0" (
    @echo off
    set exitcode=1
    goto :exit
)

:exit
@echo on
@echo start over
@echo off
if "%nopause%"=="true" (
    exit /b %exitcode%
)
pause