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
     return Math.round(Number((IOUtil.ReadAll(procPath+"/uptime").split(" "))[0].trim()))
  }

  var ticker = Scheduler.NewTicker(3109, function(){

    try {

      out = 'uptime{host="' + hostname + '"} ' + getUptime() + "\n"

      UExporter.Publish(out)

    } catch(e) {
      UExporter.Publish("# exception:"+e)
    }

  }).Start()
}
