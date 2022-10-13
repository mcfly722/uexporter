if (!String.format) {
  String.format = function(format) {
    var args = Array.prototype.slice.call(arguments, 1);
    return format.replace(/{(\d+)}/g, function(match, number) {
      return typeof args[number] != 'undefined'
        ? args[number]
        : match
      ;
    });
  };
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

function getFields(line) {
  const regexp = /\S+/gm
  return line.match(regexp)
}

function parseOutput(topCommandOutput) {
  var strings = topCommandOutput.split("\n")

  var headers = getFields(strings[6])

  var dataWithoutHeader = strings.slice(7,-1)


  return dataWithoutHeader.map(line => {
    var fields = getFields(line)
    var out = {}

    for (let i = 0;i<headers.length;i++) {
      out[headers[i]] = fields[i]
    }

    out[headers[headers.length-1]] = fields.slice(headers.length-1,fields.length).join(" ").replace(/(['"])/g,"'")

    return out
  })
}

function getCPUStat(data, firstN){
  var out = ""
  var sorted = data.sort((p1,p2) => (Number(p1['%CPU']) < Number(p2['%CPU']) ? 1:-1))
  var sortedAndFiltered = sorted.slice(0,firstN)
  var firstNSortedAndFiltered = sortedAndFiltered.map(process =>{
    return String.format('process_cpu_percentages{pid="{0}",cmd="{1}",user="{2}"} {3}',process.PID,process.COMMAND,process.USER,process['%CPU'])
  })

  out += firstNSortedAndFiltered.join("\n")

  return out
}


var ticker = Scheduler.NewTicker(7*1000, function(){

  var output = ""

  function onStdout(content){
    output += content+"\n"
  }

  function onDone(exitCode){

    var processData =  parseOutput(output)

    UExporter.Publish(getCPUStat(processData, 10))

    output = ""
  }

  Exec.NewCommand("top", ["-b","-n", "1"]).SetTimeoutMs(900).SetOnStdoutString(onStdout).SetOnDone(onDone).Start()

}).Start()
