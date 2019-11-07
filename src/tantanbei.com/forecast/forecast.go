package forecast

var frequentLevel int = 0

//time unit is second , it mean is how long the auction stared
func ForecastPrice(timeNow int64, currentPrice int) (forecastPrice int) {

	//timeNow += 3540//test code
	if timeNow < 3540 {
		return currentPrice + 1500
	}

	frequentLevel = CalculateFrequent() - standardLevel[timeNow]
	//fmt.Println("timeNow:", timeNow, "frequentLevel:", frequentLevel, "frequent dev:", CalculateFrequent()-standardLevel[timeNow])
	if timeNow < 3580 {
		return baseDeviation[timeNow] + currentPrice + frequentLevel*(3600-int(timeNow))/120/100*100
	} else {
		return baseDeviation[timeNow] + currentPrice
	}
}

func ForecastPriceAverage(timeNow int64, currentPrice int) (forecastPrice int) {
	if timeNow < 3540 {
		return currentPrice + 1500
	}

	return currentPrice + forecastAverage[timeNow]
}
