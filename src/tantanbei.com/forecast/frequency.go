package forecast

import (
	"fmt"
	"sync"
	"time"

	"project/realapi/data/auction"
	"strconv"

	"tantanbei.com/xfile"
)

const stage int64 = 5

//start time, end time and now time unit is second
var startTime int64
var endTime int64
var nowTime int64

var DataBuf map[int64]int = make(map[int64]int, 4096)
var lock sync.Mutex

func StartRecord(start, end int64, getData func() int) {
	var preTime int64 = 0
	var currentData int
	startTime = start
	endTime = end

	auction.SetForecastPrice(auction.GetCautionPrice() + 900 + 1300)

	for {
		<-time.After(100 * time.Millisecond)

		nowTime = time.Now().Unix()
		//		fmt.Println("nowTime", nowTime, "endTime", endTime)
		if nowTime < startTime {
			continue
		} else if nowTime > endTime {
			//after record end, write the result to file
			fmt.Println("time over...")
			filename := fmt.Sprint("./", nowTime, ".cache")
			dataCopy := make(map[string]int)

			lock.Lock()
			for key, value := range DataBuf {
				keyStr := strconv.FormatInt(key, 10)
				dataCopy[keyStr] = value
			}
			lock.Unlock()

			fmt.Println("write result")
			xfile.WriteJsonFile(filename, dataCopy)
			break
		}

		currentData = getData()

		if preTime == nowTime {
			continue
		}

		preTime = nowTime

		lock.Lock()
		DataBuf[nowTime-startTime] = currentData
		//fmt.Println("add data:", currentData)
		lock.Unlock()

		forecastPrice := ForecastPrice(nowTime-startTime, currentData)
		if forecastPrice != 0 {
			auction.SetForecastPrice(forecastPrice)
		}

		forecastPriceAverage := ForecastPriceAverage(nowTime-startTime, currentData)
		if forecastPriceAverage != 0 {
			auction.SetForecastPriceAverage(forecastPriceAverage)
		}

		if nowTime-startTime == 3570 {
			auction.SetSecond30(forecastPrice)
			auction.SetSecond30Average(forecastPriceAverage)
		}

		if nowTime-startTime == 3580 {
			auction.SetSecond40(forecastPrice)
			auction.SetSecond40Average(forecastPriceAverage)
		}

		if nowTime-startTime == 3585 {
			auction.SetSecond45(forecastPrice)
			auction.SetSecond45Average(forecastPriceAverage)
		}

		if nowTime-startTime == 3590 {
			auction.SetSecond50(forecastPrice)
			auction.SetSecond50Average(forecastPriceAverage)
		}
	}
}

func CalculateFrequent() (level int) {
	lock.Lock()
	defer lock.Unlock()

	now := time.Now().Unix()

	//fmt.Println(DataBuf)
	//fmt.Println("lastData:", now-startTime-1)
	lastData, ok := DataBuf[now-startTime-1]
	if !ok {
		return
	}

	//fmt.Println("stageStartData:", now-startTime-stage-1)
	stageStartData, ok := DataBuf[now-startTime-stage-1]
	if !ok {
		return
	}

	level = lastData - stageStartData
	return
}
