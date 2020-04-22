
#!/bin/sh
# set var value
set +e
source /etc/profile
myself=`which $0`
myname=`basename $myself`
mydir=`dirname $myself`
parentdir=`cd $mydir/.. ; pwd`
appname=httpmonitor
outname=$appname.out
outdir=$parentdir/out
stoppath=$mydir/stop.sh
appdir=$parentdir/app
apppath=$appdir/httpmonitor
piddir=$parentdir/pid
pidpath=$piddir/httpmonitor.pid
chmod 755 $apppath
exitcode=0

# get cmds value
cmds=$#
for((i=1;i<=$cmds;i++))
do 
    if [ $1 == "-showname" ];then
        showname=$2
    fi
    if [ $1 == "-serverport" ];then
        serverport=$2
    fi
    if [ $1 == "-gctime" ];then
        gctime=$2
    fi
    if [ $1 == "-pausetime" ];then
        pausetime=$2
    fi
    if [ $1 == "-mode" ];then
        mode=$2
    fi
    if [ $1 == "-pidfile" ];then
        pidfile=$2
    fi
    if [ $1 == "--route-root" ];then
        routeroot=$2
    fi
    if [ $1 == "--route-monitor" ];then
        routemonitor=$2
    fi
    if [ $1 == "--route-address" ];then
        routeaddress=$2
    fi
    if [ $1 == "--route-getdemo" ];then
        routegetdemo=$2
    fi
    if [ $1 == "--route-postdemo" ];then
        routepostdemo=$2
    fi
    if [ $1 == "--route-restart" ];then
        routerestart=$2
    fi
    if [ $1 == "--route-quit" ];then
        routequit=$2
    fi
    if [ $1 == "--route-stop" ];then
        routestop=$2
    fi
    if [ $1 == "-outpath" ];then
        outpath=$2
    fi
    if [ $1 == "-pidpath" ];then
        pidpath=$2
    fi
    if [ $1 == "-nohup" ];then
        isnohup=true
    fi
    if [ $1 == "--noout" ];then
        noout=--noout
    fi
    shift
done

# set cmds defalut value
if [ "$showname" == "" ];then
    showname=`hostname`
fi

if [ "$serverport" == "" ];then
    serverport=5555
fi

if [ "$gctime" == "" ];then
    gctime=60
fi

if [ "$pausetime" == "" ];then
    pausetime=60
fi

if [ "$mode" == "" ];then
    mode="release"
fi

if [ "$pidfile" == "" ];then
    pidfile=$pidpath
fi

if [ "$routeroot" == "" ];then
    routeroot="/"
fi

if [ "$routemonitor" == "" ];then
    routemonitor="/monitor"
fi

if [ "$routeaddress" == "" ];then
    routeaddress="/address"
fi

if [ "$routegetdemo" == "" ];then
    routegetdemo="/getdemo"
fi

if [ "$routepostdemo" == "" ];then
    routepostdemo="/postdemo"
fi

if [ "$routerestart" == "" ];then
    routerestart="/restart"
fi

if [ "$routequit" == "" ];then
    routequit="/quit"
fi

if [ "$routestop" == "" ];then
    routestop="/stop"
fi

if [ "$outpath" == "" ];then
    outpath=$outdir/$outname
fi

if [ ! -d "$outdir" ]; then
    mkdir -p "$outdir"
fi

if [ ! -d "$piddir" ]; then
    mkdir -p "$piddir"
fi

# kill old server
chmod 755 $stoppath
stty -echo
$stoppath
stty echo

# start
echo "start server"
if [ "$isnohup" == "true" ];then
    cd $mydir; nohup $apppath --port=$serverport --name="$showname" --mode=$mode --pid=$pidfile --gctime=$gctime --pausetime=$pausetime \
--route-root="$routeroot" --route-monitor="$routemonitor" --route-address="$routeaddress" --route-getdemo="$routegetdemo" \
--route-postdemo="$routepostdemo" --route-restart="$routerestart" --route-quit="$routequit" --route-stop="$routestop" --noout >$outpath &
else
    cd $mydir; $apppath --port=$serverport --name="$showname" --mode=$mode --pid=$pidfile --gctime=$gctime --pausetime=$pausetime \
--route-root="$routeroot" --route-monitor="$routemonitor" --route-address="$routeaddress" --route-getdemo="$routegetdemo" \
--route-postdemo="$routepostdemo" --route-restart="$routerestart" --route-quit="$routequit" --route-stop="$routestop" $noout &
fi

exit $exitcode