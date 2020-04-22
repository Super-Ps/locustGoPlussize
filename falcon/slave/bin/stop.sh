
#!/bin/sh
# set var value
set +e
source /etc/profile
myself=`which $0`
myname=`basename $myself`
mydir=`dirname $myself`
parentdir=`cd $mydir/.. ; pwd`
piddir=$parentdir/pid
pidpath=$piddir/slave.pid
exitcode=0
num=0
echo "stop slave"

# kill old slave
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
            echo "slave $num stop, pid:$pid"
            kill $pid
        fi
    fi
done < $pidpath

rm -rf $pidpath
echo "stop over"
exit $exitcode