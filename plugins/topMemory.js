var procPath = "/proc"

var firstN = 10

if (OS.Getenv("UEXPORTER_PROC_PATH") !== "") {
  procPath = OS.Getenv("UEXPORTER_PROC_PATH")
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

var hostname = OS.Getenv("UEXPORTER_HOST_NAME")
if (!hostname) {
  hostname = IOUtil.ReadAll(procPath + "/sys/kernel/hostname").trim()
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

function ToNumber(str) {
  var n = (str ?? "0").replace(/[^0-9\.]+/,"")
  return Math.round(Number(n))
}

function getAllProcesses(){

  var files = IOUtil.ReadDir(procPath)

  var processes = []

  files.forEach(file => {
    if (file.IsDir()) {
      try {
        var status = parseProcStatus(IOUtil.ReadAll(procPath + "/" + file.Name() + "/status"))
        processes.push(status)
      } catch {}
    }
  })

  return processes
}

var ticker = Scheduler.NewTicker(1000, function(){

  try {

    out = "#  RSS - Memory Resident Set Size\n"

    var processes = getAllProcesses()


    { // sort by RSS firstN
      var sorted = processes.sort((p1,p2) => (ToNumber(p1.VmRSS) < ToNumber(p2.VmRSS) ? 1:-1))
      var sortedAndFiltered = sorted.slice(0, firstN)
      var firstNSortedAndFiltered = sortedAndFiltered.map(process =>{
        return 'process_mem_rss_kb{pid="'+process.Pid+'",name="'+process.Name+'",host="'+hostname+'"} '+ToNumber(process.VmRSS)
      })
      out += firstNSortedAndFiltered.join("\n")+"\n"

    }

    { // get all others
      var allOthers = sorted.slice(firstN,-1).map(process => ToNumber(process.VmRSS))
      out += 'process_mem_res_kb{pid="0",name="allOthers",host="'+hostname+'"} ' + sum(allOthers)
    }

    //Console.Log("published")
    UExporter.Publish(out)

  } catch(e) {
    UExporter.Publish("# exception:"+e)
  }

}).Start()

Console.Log("topMemory.js v3 started")
Console.Log("procPath="+procPath+" hostname="+hostname)
