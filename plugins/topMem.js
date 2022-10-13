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

function ToNumber(str) {
  if (str.slice(-1) === 'g') {
    return Math.round(Number(str.replace(/.$/,"")) * 1024 * 1000)
  }
  return Math.round(Number(str))
}


function getResidentMemorySizeStat(data, firstN){
  out = "#  RES - Resident Memory Size\n"
  var sorted = data.sort((p1,p2) => (ToNumber(p1.RES) < ToNumber(p2.RES) ? 1:-1))
  var sortedAndFiltered = sorted.slice(0,firstN)
  var firstNSortedAndFiltered = sortedAndFiltered.map(process =>{
    return String.format('process_mem_res_kb{pid="{0}",cmd="{1}",user="{2}"} {3}',process.PID,process.COMMAND,process.USER,ToNumber(process.RES))
  })

  out += firstNSortedAndFiltered.join("\n")+"\n"

  var allOthers = sorted.slice(firstN,-1).map(process => ToNumber(process.RES))

  out += String.format('process_mem_res_kb{pid="-1",cmd="others",user="others"} {0}', sum(allOthers))

  return out
}

var ticker = Scheduler.NewTicker(11*1000, function(){

  var output = ""

  function onStdout(content){
    output += content+"\n"
  }

  function onDone(exitCode){

    var processData =  parseOutput(output)

    UExporter.Publish(getResidentMemorySizeStat(processData, 20))

    output = ""
  }

  Exec.NewCommand("top", ["-b","-n", "1"]).SetTimeoutMs(900).SetOnStdoutString(onStdout).SetOnDone(onDone).Start()

}).Start()
