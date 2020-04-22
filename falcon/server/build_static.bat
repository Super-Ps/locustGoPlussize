
@echo off
set myself=%~dp0
set myroot=%myself:~,1%
cd %myself%
%myroot%:
set sourcepath=./index.html ./static/...
set descpath=./service/static.go
set packagename=service
go-bindata -o=%descpath% -pkg=%packagename% %sourcepath%