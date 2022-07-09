
var ticker = Scheduler.NewTicker(1*1000, function(){
	Console.Log("timer1")
}).Start()

Console.Log("timer started")


/*
var count1 = 0
timerId1 = setInterval(function(){
	Console.log("timer1="+count1)

	if (count1 > 5) {
		Console.log("clearInterval("+timerId1+")")
		clearInterval(timerId1)
	}
	               1/0
	count1++
},1000,500)

var count2 = 0
timerId2 = setInterval(function(){
	Console.log("timer2="+count2)

	if (count2 > 5) {
		Console.log("clearInterval("+timerId2+")")
		clearInterval(timerId2)
	}

	count2++
},600,500)

clearInterval(23)
*/
