
@echo off
set myself=%~dp0
set myroot=%myself:~,1%
cd %myself%
%myroot%:
set sourcepath=./winapi.dll
set descpath=./windll_windows.go
set packagename=monitor
go-bindata -o=%descpath% -pkg=%packagename% %sourcepath%