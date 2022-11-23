// plugin could be used from container to obtain CPU info from it hosting k8s node

var supportedOS = ['linux']
var firstN = 20
var procPath = "/proc"

if (!(supportedOS.includes(OS.OS()))) {
  var msg = "OS="+OS.OS()+ " is not supported. List of supported OS:["+supportedOS+"]"
  Console.Log(msg)
  UExporter.Publish(msg)
} else {

  if (OS.Getenv("UEXPORTER_PROC_PATH") !== "") {
    procPath = OS.Getenv("UEXPORTER_PROC_PATH")
  }

  var hostname = OS.Getenv("UEXPORTER_HOST_NAME")
  if (!hostname) {
    hostname = IOUtil.ReadAll(procPath + "/sys/kernel/hostname").trim()
  }

  function sum(array) {
    var total = 0;
    for (var i = 0; i < array.length; i++) {
      if (Number.isInteger(array[i])) {
        total += array[i];
      }
    }
    return total;
  }

  function getUptime(){
    return Number(new Date().getTime())
  }

  // https://man7.org/linux/man-pages/man5/proc.5.html
  function parseProcStat(line){
    var values = line.split(" ")
    return {
      'pid'  : values[0],
      'comm' : values[1].replace(/^\(/,"").replace(/\)$/,""),
      'utime': Number(values[13]),
      'stime': Number(values[14])
    }
  }

  function getProcessesStats(){
    var files = IOUtil.ReadDir(procPath)

    var processes = []

    files.forEach(file => {
      if (file.IsDir()) {
        try {
          var stat = parseProcStat(IOUtil.ReadAll(procPath + "/" + file.Name() + "/stat"))

          processes.push(stat)
        } catch {}
      }
    })
    return processes
  }


  function groupBy(xs, key) {
    return xs.reduce(function(rv, x) {
      (rv[x[key]] = rv[x[key]] || []).push(x);
      return rv;
    }, {});
  };

  var previousUptime = 0
  var previousProcStatsIndex = {}

  function getCPUProcStat() {

    var currentProcStat = getProcessesStats()
    var curentUptime = getUptime()

    var cpuStat = currentProcStat.map(procStat => {

      try {
        procStat['prev_utime'] = (previousProcStatsIndex[procStat.pid])[0].utime
        procStat['prev_stime'] = (previousProcStatsIndex[procStat.pid])[0].stime
        procStat['prev_cpu']   = (previousProcStatsIndex[procStat.pid])[0].cpu
        procStat['interval'] = curentUptime - previousUptime

        procStat['cpu'] =  Number((100 * ((procStat.utime - procStat.prev_utime) + (procStat.stime - procStat.prev_stime))) / (procStat.interval/10))

        if (isNaN(procStat['cpu'])){
          procStat['cpu'] = Number(0)
        }

      } catch {}


      return procStat
    })

    previousUptime = curentUptime

    previousProcStatsIndex = groupBy(currentProcStat, 'pid')

    return cpuStat
  }


  var ticker = Scheduler.NewTicker(3191, function(){

    try {

      out = ""

      try {

        var stat = getCPUProcStat()

        var filtered =stat.filter(process => process.cpu !== undefined)
        var sorted = filtered.sort((p1,p2) => (Number(p1.cpu) < Number(p2.cpu) ? 1:-1))
        var topN = sorted.slice(0, firstN)
        var topNLines = topN.map(process =>{
          return 'process_cpu_percentages{pid="' + process.pid + '",name="'+process.comm + '",host="'+hostname+'"} ' + Math.round(100 * process.cpu) / 100
        })

        out += topNLines.join("\n")+"\n"

      } catch {}

      UExporter.Publish(out)

    } catch(e) {
      UExporter.Publish("# exception:"+e)
    }

  }).Start()
}
