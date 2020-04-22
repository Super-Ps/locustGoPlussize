
#!/bin/sh
# set var value
set +e
source /etc/profile
gorootpath=$GOROOT
gobinpath=go
myself=`which $0`
myname=`basename $myself`
mydir=`dirname $myself`
parentdir=`cd $mydir/.. ; pwd`
appdir=$parentdir/app
apppath=$appdir/httpmonitor
bulidresult=""
exitcode=0

# bulid
echo "build start"
if [ -f "$apppath" ]; then
    rm -rf $apppath
fi

export GOROOT=$gorootpath
export GOPATH=`cd $mydir/../../../.. ; pwd`
$gobinpath build -o $apppath falcon/server
err=$?

if [ $err != 0 ];then
    exitcode=1
    bulidresult="build fail"
else
    exitcode=0
    bulidresult="build success"
    chmod 755 $apppath
fi

echo $bulidresult
exit $exitcode