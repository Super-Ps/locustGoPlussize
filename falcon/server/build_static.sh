
#!/bin/sh
# set var value
set +e
source /etc/profile
myself=`which $0`
myname=`basename $myself`
mydir=`dirname $myself`
cd $mydir
sourcepath=./index.html ./static/...
descpath=./service/static.go
packagename=service
go-bindata -o=$descpath -pkg=$packagename $sourcepath