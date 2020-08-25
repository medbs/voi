package core

type Metrics struct {
	MinLatency float64
	MaxLatency float64
	AvgLatency float64
}

func CalculateMinMaxLatency(pms []PingMessage) (float64, float64) {
	var minLatency float64 = 10000
	var maxLatency float64 = 0

	for i := 0; i < len(pms); i++ {
		st := calculateSendingTime(&pms[i])
		if st < minLatency {
			minLatency = st
		}

		if st > maxLatency {
			maxLatency = st
		}
	}

	return maxLatency, maxLatency
}

func CalculateAvgLatency(pms []PingMessage) float64 {

	var avg float64 = 0

	for i := 0; i < len(pms); i++ {
		st := calculateSendingTime(&pms[i])
		avg = avg + st
	}

	return avg
}

func calculateSendingTime(message *PingMessage) float64 {
	elapsed := message.ReceivedTime.Sub(message.SentTime)
	return elapsed.Seconds()
}
