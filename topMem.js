var counter = 1

var ticker = Scheduler.NewTicker(3*1000, function(){

  counter++
  Console.Log(counter)
  UExporter.Publish("counter:"+counter)
}).Start()
