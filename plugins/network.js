// plugin does not works correcly inside container. To obtain correct info from /proc/net/dev you have to start uexporter as node service.
// (container overwrites /proc/net/dev file with it self statistics)

var firstN = 30
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

  function getUptime(){
    return Number(new Date().getTime())
  }

  function groupBy(xs, key) {
    return xs.reduce(function(rv, x) {
      (rv[x[key]] = rv[x[key]] || []).push(x);
      return rv;
    }, {});
  };

  function parseInterfacesStat(content){
    var lines = content.split("\n").filter(line => line.indexOf(':') > -1)
    interfaces = []
    lines.forEach(line => {
        var values = line.substring(line.split(":")[0].length + 1, line.length).trim()
        var stats = values.match(/\S+/g)
        var interfaceName = line.split(":")[0].trim()
        var interface = {
          'interface' : interfaceName,
          'receive'   : Number(stats[0]),
          'transmit'  : Number(stats[8])
        }
        interfaces.push(interface)
    })
    return interfaces
  }

  var previousUptime = 0
  var previousStatsIndex = {}

  function getInterfaceNetStatByTime() {

    var currentStat = parseInterfacesStat(IOUtil.ReadAll(procPath + "/net/dev"))
    var curentUptime = getUptime()

    var stat = currentStat.map(interfaceStat => {

      try {
        var prev = (previousStatsIndex[interfaceStat.interface])[0]

        interfaceStat['interval']          = curentUptime - previousUptime
        interfaceStat['prev_receive']      = prev.receive
        interfaceStat['prev_transmit']     = prev.transmit
        interfaceStat['receive_rate']      = Number(1000 * (interfaceStat.receive   - interfaceStat.prev_receive ) / interfaceStat.interval)
        interfaceStat['transmit_rate']     = Number(1000 * (interfaceStat.transmit  - interfaceStat.prev_transmit) / interfaceStat.interval)
      } catch {}

      return interfaceStat
    })

    previousUptime = curentUptime
    previousStatsIndex = groupBy(currentStat, 'interface')

    return stat
  }

  var ticker = Scheduler.NewTicker(3091, function(){

    try {

      out = ""

      var interfacesStats = getInterfaceNetStatByTime()

      { // receive_rate
        var topLines = interfacesStats.map(interface =>{
          return 'net_received_kbps{host="' + hostname + '",interface="' + interface.interface + '"} ' + Math.round(interface.receive_rate / 1000)
        })
        out += topLines.join("\n")+"\n"
      }

      { // transmit_rate
        var topLines = interfacesStats.map(interface =>{
          return 'net_transferred_kbps{host="' + hostname + '",interface="' + interface.interface + '"} ' + Math.round(interface.transmit_rate / 1000)
        })
        out += topLines.join("\n")+"\n"
      }


      UExporter.Publish(out)

    } catch(e) {
      UExporter.Publish("# exception:"+e)
    }

  }).Start()
}
