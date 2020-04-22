
#!/bin/sh
# set var value
set +e
source /etc/profile
myself=`which $0`
myname=`basename $myself`
mydir=`dirname $myself`
parentdir=`cd $mydir/.. ; pwd`
piddir=$parentdir/pid
pidpath=$piddir/httpmonitor.pid
exitcode=0
num=0
echo "stop server"

# kill old server
if [ ! -f "$pidpath" ]; then
    echo "stop over"
    exit $exitcode
fi

while read pid
do
    if [ "$pid" != "" ];then
        pidstatus=`lsof -p $pid` || echo ""
        if [ "$pidstatus" != "" ];then
            num=`expr $num + 1`
            echo "server $num stop, pid:$pid"
            kill -9 $pid
        fi
    fi
done < $pidpath

if [ -f "$pidpath" ]; then
    rm -rf $pidpath
fi
echo "stop over"
exit $exitcode