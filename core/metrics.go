package core

func CalculateSendingTime(message *PingMessage) float64{
	elapsed := message.ReceivedTime.Sub(message.SentTime)
	return elapsed.Seconds()
}


func CalculateMinSendingTime(){

}

func CalculateMaxSendingTime(){

}