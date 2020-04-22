
@echo off

set allavg=
:begin
set avgvalue=
set avgvalue=%1%
if "%avgvalue%"=="" (goto :setvar)
if "%avgvalue%"=="-masterhost" (set masterhost=%2%)
if "%avgvalue%"=="-masterport" (set masterport=%2%)
if "%avgvalue%"=="-conf" (set conf=%2%)
if "%avgvalue%"=="-after" (set after=%2%)
if "%avgvalue%"=="-slaves" (set slaves=%2%)
if "%avgvalue%"=="-gctime" (set gctime=%2%)
if "%avgvalue%"=="-local" (set local=--local)
if "%avgvalue%"=="-random" (set israndom=--random)
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
set bindir=%parent%\bin
set appdir=%parent%\app
set outdir=%parent%\out
set confdir=%parent%\conf
set apppath=%appdir%\slave.exe
set confpath=%confdir%\main.yml
set outpath=%outdir%\slave.out
set start_time=%date:~0,4%-%date:~5,2%-%date:~8,2% %time:~0,2%:%time:~3,2%:%time:~6,2%
set israndom=--random
cd %appdir%

if "%masterhost%"=="" (set masterhost=127.0.0.1)
if "%masterport%"=="" (set masterport=5557)
if "%conf%"=="" (set conf=%confpath%)
if "%after%"=="" (set after=0)
if "%gctime%"=="" (set gctime=60000)
if "%slaves%"=="" (set slaves=1)
if not exist "%outdir%" (mkdir "%outdir%")
if exist "%outpath%" (del /f /q "%outpath%")
echo %start_time% >> "%outpath%"
@echo on

@echo start slave
%apppath% --slave --master-host=%masterhost% --master-port=%masterport% --config="%conf%" --after=%after% --gctime-interval=%gctime% %local% %israndom%
@echo start over
pause

:exit
@echo off
exit
