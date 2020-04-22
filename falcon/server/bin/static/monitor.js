
defaultButton();
var keyWord = "";
var pidList = new Array();
var isMonitor = 0;
var host = ".";
var monitorUrl = host.trim() + $("#box_custom a.start-button").attr("href").trim();
var stopUrl = host.trim() + $("#box_custom a.stop-button").attr("href").trim();
var restartUrl = host.trim() + $("#box_custom a.restart-button").attr("href").trim();
var exitUrl = host.trim() + $("#box_custom a.exit-button").attr("href").trim();
var monitor_tpl = $("#monitor-template");
var info_tpl = $("#info-template");
var process_tpl = $("#pro-monitor-template");
var alternate = false;
var desc = false;
var processSortName = "name";
var processMap = {};
var svrCpuChart = new MonitorLineChart($("#server_charts_container"), "Cpu percent (%/s)", ["%"], "%");
var svrMemoryChart = new MonitorLineChart($("#server_charts_container"), "Memory used (MB)", ["Rss", "Vms"], "MB");
var svrNetIoSpeedChart = new MonitorLineChart($("#server_charts_container"), "Net io speed (Mb/s)", ["Recv", "Send"], "Mb/s");
var SvrNetIoCountChart = new MonitorLineChart($("#server_charts_container"), "Net io count (Packets/s)", ["Recv", "Send"], "Packets/s");
var svrDiskIoSpeedChart = new MonitorLineChart($("#server_charts_container"), "Disk io speed (MB/s)", ["Read", "Write"], "MB/s");
var svrDiskIoCountChart = new MonitorLineChart($("#server_charts_container"), "Disk io count (Count/s)", ["Read", "Write"], "Count/s");

var proCpuChart = new MonitorLineChart($("#process_charts_container"), "Cpu percent (%/s)", [], "%");
var proMemoryRssChart = new MonitorLineChart($("#process_charts_container"), "Memory rss used (MB)", [], "MB");
var proMemoryVmsChart = new MonitorLineChart($("#process_charts_container"), "Memory vms used (MB)", [], "MB");
var proDiskIoReadSpeedChart = new MonitorLineChart($("#process_charts_container"), "Disk io read speed (MB/s)", [], "MB/s");
var proDiskIoWriteSpeedChart = new MonitorLineChart($("#process_charts_container"), "Disk io write speed (MB/s)", [], "MB/s");
var proDiskIoReadCountChart = new MonitorLineChart($("#process_charts_container"), "Disk io read count (Count/s)", [], "Count/s");
var proDiskIoWriteCountChart = new MonitorLineChart($("#process_charts_container"), "Disk io write count (Count/s)", [], "Count/s");

var sortBy = function(field, reverse, primer){
    reverse = (reverse) ? -1 : 1;
    return function(a,b){
        a = a[field];
        b = b[field];
       if (typeof(primer) != "undefined"){
           a = primer(a);
           b = primer(b);
       }
       if (a<b) return reverse * -1;
       if (a>b) return reverse * 1;
       return 0;
    }
}

$("#box_custom a.start-button").click(function(event) {
    event.preventDefault();
    $("#box_custom a.start-button").hide();
    $("#box_custom a.stop-button").show();
    defaultServerCharts();
    defaultProcessCharts();
    isMonitor = 1;
    updateMonitor();
});

$("#box_custom a.stop-button").click(function(event) {
    event.preventDefault();
    $("#box_custom a.stop-button").hide();
    $("#box_custom a.start-button").show();
    isMonitor = 0;
    stopMonitor();
});

$("#box_custom a.restart-button").click(function(event) {
    event.preventDefault();
    $("#box_custom a.restart-button").hide();
    $("#box_custom a.stop-button").hide();
    $("#box_custom a.start-button").hide();
    $("#box_custom a.exit-button").hide();
    $("#box_custom a.message-text").show();
    $("#box_custom a.message-text").html("Please wait while restarting");
    isMonitor = 0;
    restartMonitor();
    setTimeout(defaultButton, 5000);
});

$("#box_custom a.exit-button").click(function(event) {
    event.preventDefault();
    $("#box_custom a.restart-button").hide();
    $("#box_custom a.stop-button").hide();
    $("#box_custom a.start-button").hide();
    $("#box_custom a.exit-button").hide();
    $("#box_custom a.message-text").show();
    $("#box_custom a.message-text").html("Service exited");
    isMonitor = 0;
    exitMonitor();
});

$("#keyword_search_button").click(function(event) {
    event.preventDefault();
    keyWord = $("#key_word").val().trim();
});

$("#pid_search_button").click(function(event) {
    event.preventDefault();
    var processPid = $("#process_pid").val().trim();
    if (processPid == "") {
        pidList = new Array()
    } else {
        pidList = processPid.split(",");
    }
    defaultProcessCharts();
});

$("ul.tabs").tabs("div.panes > div").on("onClick", function(event) {
    event.preventDefault();
    if (event.target.id == "server_charts_tab") {
        svrCpuChart.resize();
        svrMemoryChart.resize();
        svrNetIoSpeedChart.resize();
        SvrNetIoCountChart.resize();
        svrDiskIoSpeedChart.resize();
        svrDiskIoCountChart.resize();
    }

    if (event.target.id == "process_charts_tab") {
        proCpuChart.resize();
        proMemoryRssChart.resize();
        proMemoryVmsChart.resize();
        proDiskIoReadSpeedChart.resize();
        proDiskIoWriteSpeedChart.resize();
        proDiskIoReadCountChart.resize();
        proDiskIoWriteCountChart.resize();
    }
});

$("#process_monitor th.stats_label").click(function(event) {
    event.preventDefault();
    processSortName = $(this).attr("data-sortkey");
    desc = !desc;
});

function defaultButton() {
    $("#box_custom a.restart-button").show();
    $("#box_custom a.start-button").show();
    $("#box_custom a.exit-button").show();
    $("#box_custom a.stop-button").hide();
    $("#box_custom a.message-text").hide();
    $("#box_custom a.message-text").html("");
}

function defaultServerCharts() {
    $("#server_charts_container").html("");
    svrCpuChart = new MonitorLineChart($("#server_charts_container"), "Cpu percent (%/s)", ["%"], "%");
    svrMemoryChart = new MonitorLineChart($("#server_charts_container"), "Memory used (MB)", ["Rss", "Vms"], "MB");
    svrNetIoSpeedChart = new MonitorLineChart($("#server_charts_container"), "Net io speed (Mb/s)", ["Recv", "Send"], "Mb/s");
    SvrNetIoCountChart = new MonitorLineChart($("#server_charts_container"), "Net io count (Packets/s)", ["Recv", "Send"], "Packets/s");
    svrDiskIoSpeedChart = new MonitorLineChart($("#server_charts_container"), "Disk io speed (MB/s)", ["Read", "Write"], "MB/s");
    svrDiskIoCountChart = new MonitorLineChart($("#server_charts_container"), "Disk io count (Count/s)", ["Read", "Write"], "Count/s");
}

function defaultProcessCharts() {
    $("#process_charts_container").html("");
    if (pidList.length == 0) {
        proCpuChart = new MonitorLineChart($("#process_charts_container"), "Cpu percent (%/s)", [], "%");
        proMemoryRssChart = new MonitorLineChart($("#process_charts_container"), "Memory rss used (MB)", [], "MB");
        proMemoryVmsChart = new MonitorLineChart($("#process_charts_container"), "Memory vms used (MB)", [], "MB");
        proDiskIoReadSpeedChart = new MonitorLineChart($("#process_charts_container"), "Disk io read speed (MB/s)", [], "MB/s");
        proDiskIoWriteSpeedChart = new MonitorLineChart($("#process_charts_container"), "Disk io write speed (MB/s)", [], "MB/s");
        proDiskIoReadCountChart = new MonitorLineChart($("#process_charts_container"), "Disk io read count (Count/s)", [], "Count/s");
        proDiskIoWriteCountChart = new MonitorLineChart($("#process_charts_container"), "Disk io write count (Count/s)", [], "Count/s");
    } else {
        var pidNameList = new Array();
        for (i=0; i<pidList.length ;i++ ) {
            var pid = pidList[i];
            if (processMap.hasOwnProperty(pid) != false) {
                pidNameList.push(pid.toString() + "(" + processMap[pid].name + ")");
            } else {
                pidNameList.push(pid);
            }
        }
        proCpuChart = new MonitorLineChart($("#process_charts_container"), "Cpu percent (%/s)", pidNameList, "%");
        proMemoryRssChart = new MonitorLineChart($("#process_charts_container"), "Memory rss used (MB)", pidNameList, "MB");
        proMemoryVmsChart = new MonitorLineChart($("#process_charts_container"), "Memory vms used (MB)", pidNameList, "MB");
        proDiskIoReadSpeedChart = new MonitorLineChart($("#process_charts_container"), "Disk io read speed (MB/s)", pidNameList, "MB/s");
        proDiskIoWriteSpeedChart = new MonitorLineChart($("#process_charts_container"), "Disk io write speed (MB/s)", pidNameList, "MB/s");
        proDiskIoReadCountChart = new MonitorLineChart($("#process_charts_container"), "Disk io read count (Count/s)", pidNameList, "Count/s");
        proDiskIoWriteCountChart = new MonitorLineChart($("#process_charts_container"), "Disk io write count (Count/s)", pidNameList, "Count/s");
    }
}

function lengthList(jsonData){
    var jsonLength = 0;
    for(var item in jsonData){
        jsonLength++;
    }
    return jsonLength;
}

function stopMonitor() {
    $.post(stopUrl)
}

function restartMonitor() {
    $.post(restartUrl)
}

function exitMonitor() {
    $.post(exitUrl)
}

function updateMonitor() {
    if (isMonitor == 0) {
        return
    }
    $.get((keyWord == "") ? monitorUrl : monitorUrl + "?keyword=" + keyWord, function (body) {
        if (body.hasOwnProperty("data") == true){
            $("#machine_info tbody").empty();
            $('#machine_info tbody').jqoteapp(info_tpl, body.data);
            $('#machine_info tbody').jqoteapp(info_tpl, {"host_id": ""});

            $("#server_monitor tbody").empty();
            $('#server_monitor tbody').jqoteapp(monitor_tpl, body.data.monitor.system);
            $('#server_monitor tbody').jqoteapp(monitor_tpl, {"last_time": ""});

            svrCpuChart.addValue([body.data.monitor.system.cpu_percent]);
            svrMemoryChart.addValue([((body.data.monitor.system.rss)/1024/1024).toFixed(2), ((body.data.monitor.system.vms)/1024/1024).toFixed(2)]);
            svrNetIoSpeedChart.addValue([((body.data.monitor.system.net_bytes_recv_diff)/1024/1024).toFixed(2), ((body.data.monitor.system.net_bytes_sent_diff)/1024/1024).toFixed(2)]);
            SvrNetIoCountChart.addValue([body.data.monitor.system.net_packets_recv_diff, body.data.monitor.system.net_packets_sent_diff]);
            svrDiskIoSpeedChart.addValue([((body.data.monitor.system.io_read_bytes_diff)/1024/1024).toFixed(2), ((body.data.monitor.system.io_write_bytes_diff)/1024/1024).toFixed(2)]);
            svrDiskIoCountChart.addValue([body.data.monitor.system.io_read_count_diff, body.data.monitor.system.io_write_count_diff]);

            if (body.data.monitor.process != null){
                processMap = null;
                processMap = {};
                for(i = 0, len=body.data.monitor.process.length; i < len; i++) {
                    processMap[body.data.monitor.process[i].pid] = body.data.monitor.process[i];
                }

                totalRow = body.data.monitor.process.pop();
                sortedStats = (body.data.monitor.process).sort(sortBy(processSortName, desc));
                sortedStats.push(totalRow);
                $("#process_monitor tbody").empty();
                $("#process_monitor tbody").jqoteapp(process_tpl, sortedStats);
                $('#process_monitor tbody').jqoteapp(process_tpl, {"pid": ""});

                if (pidList.length != 0) {
                    var proCpuList = new Array();
                    var proMemRssList = new Array();
                    var proMemVmsList = new Array();
                    var proDiskIoReadSpeedList = new Array();
                    var proDiskIoWriteSpeedList = new Array();
                    var proDiskIoReadCountList = new Array();
                    var proDiskIoWriteCountList = new Array();
                    for (i=0; i<pidList.length ;i++ )
                    {
                        var pid = pidList[i];
                        if (processMap.hasOwnProperty(pid) == false) {
                            proCpuList.push(0);
                            proMemRssList.push(0);
                            proMemVmsList.push(0);
                            proDiskIoReadSpeedList.push(0);
                            proDiskIoWriteSpeedList.push(0);
                            proDiskIoReadCountList.push(0);
                            proDiskIoWriteCountList.push(0);
                        } else {
                            proCpuList.push(processMap[pid].cpu_percent);
                            proMemRssList.push(((processMap[pid].rss)/1024/1024).toFixed(2));
                            proMemVmsList.push(((processMap[pid].vms)/1024/1024).toFixed(2));
                            proDiskIoReadSpeedList.push(((processMap[pid].io_read_bytes_diff)/1024/1024).toFixed(2));
                            proDiskIoWriteSpeedList.push(((processMap[pid].io_write_bytes_diff)/1024/1024).toFixed(2));
                            proDiskIoReadCountList.push(processMap[pid].io_read_count_diff);
                            proDiskIoWriteCountList.push(processMap[pid].io_write_count_diff);
                        }
                    }
                    proCpuChart.addValue(proCpuList);
                    proMemoryRssChart.addValue(proMemRssList);
                    proMemoryVmsChart.addValue(proMemVmsList);
                    proDiskIoReadSpeedChart.addValue(proDiskIoReadSpeedList);
                    proDiskIoWriteSpeedChart.addValue(proDiskIoWriteSpeedList);
                    proDiskIoReadCountChart.addValue(proDiskIoReadCountList);
                    proDiskIoWriteCountChart.addValue(proDiskIoWriteCountList);
                }
            }
        }
    });
    setTimeout(updateMonitor, 2000);
}

function formatSeconds(value) {
    if (value == undefined){
        value = "00:00:00";
    }    
    var second = parseInt(value)
    var min = 0
    var hour = 0
    if(second > 60) {
        min = parseInt(second/60);
        second = parseInt(second%60);
        if(min > 60) {
            hour = parseInt(min/60);
            min = parseInt(min%60);
            }
        }
    
    second = second.toString()
    min =  min.toString()
    hour =  hour.toString()
    if (second.length == 1) {
        second = "0" + second;
    }
    if (min.length == 1) {
        min = "0" + min;
    }
    if (hour.length == 1) {
        hour = "0" + hour;
    }
    return (hour + ":" + min + ":" + second);
}