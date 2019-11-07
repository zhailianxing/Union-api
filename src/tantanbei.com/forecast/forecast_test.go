package forecast

import (
	"fmt"
	"testing"
	"time"
)

func TestForecastTransactionPrice(t *testing.T) {
	var serverTime int64 = time.Now().Unix()
	var overTime int64 = serverTime + 80

	for overTime > serverTime {
		serverTime = time.Now().Unix()

		fmt.Println("timeDivider:", (overTime-serverTime)/1000)

		forecastTransactionPrice := ForecastPrice(overTime-serverTime, 85200)
		fmt.Println(forecastTransactionPrice)
		<-time.After(1 * time.Second)
	}
}
