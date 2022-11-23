
//  This plugin does not works as I expected.
//  /proc/*/net/netstat and /proc/*/net/dev contains network counters for each process, but it is only for IP protocol, not for sockets/TCP/UDP.
//  Thats why if you will start this plugin it would show same values for different processes even if only one downloads some content.
//  Unfortunatelly, I didnt find correct way to obtain TCP+UDP traffic for each process to separate it from other processes. If you know how to do it from /proc, give me, please, to know :)


var firstN = 30
var procPath = "/proc"
var supportedOS = ['linux']

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

  function getUptime(){
    return Number(new Date().getTime())
  }

  function groupBy(xs, key) {
    return xs.reduce(function(rv, x) {
      (rv[x[key]] = rv[x[key]] || []).push(x);
      return rv;
    }, {});
  };

  function sum(array) {
    var total = 0;
    for (var i = 0; i < array.length; i++) {
      if (Number.isInteger(array[i])) {
        total += array[i];
      }
    }
    return total;
  }

  function parseProcStatus(content){
    var result = {}

    content.split("\n").forEach(line => {
      var key = line.split(":")[0]
      var value = line.substring(key.length+1).trim()
      result[key] = value
    })

    return result
  }

  function parseNetStatus(content, pid, name){
    var lines = content.split("\n").filter(line => line.indexOf(':') > -1)

    interfaces = []

    lines.forEach(line => {
        var values = line.substring(line.split(":")[0].length + 1, line.length).trim()

        var stats = values.match(/\S+/g)
        var interfaceName = line.split(":")[0].trim()

        var interface = {
          'id'        : pid+':'+interfaceName,
          'pid'       : pid,
          'pname'     : name,
          'interface' : interfaceName,
          'receive'   : Number(stats[0]),
          'transmit'  : Number(stats[8])
        }

        interfaces.push(interface)
    })

    return interfaces
  }


  function getInterfaceNetStat() {

    var files = IOUtil.ReadDir(procPath)

    var interfaces = []

    files.forEach(file => {
      if (file.IsDir()) {
        try {
          var status = parseProcStatus(IOUtil.ReadAll(procPath + "/" + file.Name() + "/status"))
          var nets   = parseNetStatus(IOUtil.ReadAll(procPath + "/" + file.Name() + "/net/dev"), status.Pid, status.Name)
          nets.forEach(net => {interfaces.push(net)})
        } catch {}
      }
    })

    return interfaces
  }

  var previousUptime = 0
  var previousProcNetStatsIndex = {}

  function getInterfaceNetStatByTime() {

    var currentProcNetStat = getInterfaceNetStat()
    var curentUptime = getUptime()

    var stat = currentProcNetStat.map(procStat => {

      try {
        var prev = (previousProcNetStatsIndex[procStat.id])[0]
        procStat['interval']          = curentUptime - previousUptime
        procStat['prev_receive']      = prev.receive
        procStat['prev_transmit']     = prev.transmit
        procStat['recv']              = Number(1000 * (procStat.receive   - procStat.prev_receive ) / procStat.interval)
        procStat['send']              = Number(1000 * (procStat.transmit  - procStat.prev_transmit) / procStat.interval)
      } catch {}

      return procStat
    })

    previousUptime = curentUptime
    previousProcNetStatsIndex = groupBy(currentProcNetStat, 'id')

    return stat
  }

  var ticker = Scheduler.NewTicker(3081, function(){

    try {

      out = ""

      var stat = getInterfaceNetStatByTime()

      { // top of receivers
        var filtered = stat.filter(p => (p.recv > 0))
        var sorted   = filtered.sort((p1,p2) => (p1.recv < p2.recv ? 1:-1))
        var top      = sorted.slice(0, firstN)

        var topLines = top.map(process =>{
          return 'process_net_received_kb_per_sec{host="'+hostname+'",pid="'+process.pid+'",pname="'+process.pname+'",interface="'+process.interface+'"} ' + Math.round(process.recv / 1000)
        })

        out += topLines.join("\n")+"\n"
      }

      { // top of senders
        var filtered = stat.filter(p => (p.send > 0))
        var sorted   = filtered.sort((p1,p2) => (p1.send < p2.send ? 1:-1))
        var top      = sorted.slice(0, firstN)

        var topLines = top.map(process =>{
          return 'process_net_sended_kb_per_sec{host="'+hostname+'",pid="'+process.pid+'",pname="'+process.pname+'",interface="'+process.interface+'"} ' + Math.round(process.send / 1000)
        })

        out += topLines.join("\n")+"\n"
      }

      UExporter.Publish(out)

    } catch(e) {
      UExporter.Publish("# exception:"+e)
    }

  }).Start()

}
