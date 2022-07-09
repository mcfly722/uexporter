
var ticker = Scheduler.NewTicker(1*1000, function(){


	Console.Log("exec!")

	var output = ""

	function onStdout(content){
		  output+=content+'\r\n'
	}

	function onDone(content){
		Console.Log("content="+output)
	}

	Exec.NewCommand("ping.exe", ["-n","2", "localhost"]).SetTimeoutMs(900).SetOnStdoutString(onStdout).SetOnDone(onDone).Start()

}).Start()

Console.Log("timer started")
