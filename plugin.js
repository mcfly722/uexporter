
var output = ""

var ticker = Scheduler.NewTicker(3*1000, function(){


	function onStdout(content){
		output += content+"\n"
		UExporter.Publish(output)
	}

	function onDone(exitCode){
	}

	Exec.NewCommand("ping.exe", ["-n","2", "localhost"]).SetTimeoutMs(900).SetOnStdoutString(onStdout).SetOnDone(onDone).Start()


}).Start()

Console.Log("timer started")
