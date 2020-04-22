
#!/bin/sh
# set var value
set +e
source /etc/profile
myself=`which $0`
myname=`basename $myself`
mydir=`dirname $myself`
parentdir=`cd $mydir/.. ; pwd`
stoppath=$mydir/stop.sh
appdir=$parentdir/app
apppath=$appdir/slave
confdir=$parentdir/conf
confpath=$confdir/main.yml
outdir=$parentdir/out
outpath=$outdir/slave.out
random=--random
chmod 755 $apppath
exitcode=0

# get cmds value
cmds=$#
for((i=1;i<=$cmds;i++))
do 
    if [ $1 == "-masterhost" ];then
        masterhost=$2
    fi
    if [ $1 == "-masterport" ];then
        masterport=$2
    fi
    if [ $1 == "-conf" ];then
        conf=$2
    fi
    if [ $1 == "-after" ];then
        after=$2
    fi
    if [ $1 == "-slaves" ];then
        slaves=$2
    fi
    if [ $1 == "-gctime" ];then
        gctime=$2
    fi
    if [ $1 == "-local" ];then
        local="--local"
    fi
    if [ $1 == "-random" ];then
        random="--random"
    fi
    shift
done

# set cmds defalut value
if [ "$masterhost" == "" ];then
    masterhost="127.0.0.1"
fi

if [ "$masterport" == "" ];then
    masterport=5557
fi

if [ "$conf" == "" ];then
    conf=$confpath
fi

if [ "$after" == "" ];then
    after=0
fi

if [ "$slaves" == "" ];then
    slaves=1
fi

if [ "$gctime" == "" ];then
    gctime=60000
fi

# create dir
if [ ! -f "$outdir" ]; then
    mkdir -p $outdir
fi

if [ -f "$outpath" ]; then
    rm -rf $outpath
fi

# kill old slave
chmod 755 $stoppath
# $stoppath

# start
echo "start slave"
for((i=1;i<=$slaves;i++))
do
    nohup $apppath --master-host=$masterhost --master-port=$masterport --config=$conf --after=$after --gctime-interval=$gctime $local $random >>$outpath 2>&1 &
    echo "slave $i start, pid:$!"
done
echo "start over"
exit $exitcode