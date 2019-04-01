package service

import (
	"context"
	"grab-crawler/common"
	"grab-crawler/dto"
)

const name = "https://www.wongnai.com/"

func Scrap() {
	Queryctx, cancel := context.WithCancel(context.Background())
	ConsumeCtx, cancelConsume := context.WithCancel(context.Background())
	body, err := common.MakeRequest(name, "_api/regions.json", "_v", "5.056", "locale", "th", "knownLocation", "false")
	if err != nil {
		panic(err.Error())
	}
	regionResp := &dto.RegionResp{}
	common.UnmarshallRegion(body, regionResp)
	println(regionResp)
	//var regionChan chan *dto.City
	//for _,city := range regionResp.Cities {
	//	regionChan <- city
	//
	//}
	//entities :=common.QueryCityAll(regionResp.Cities[77])
	ch := make(chan *dto.THEntity, 20)
	rowCh := make(chan []string, 20)
	go common.QueryCityAllCh(regionResp.Cities, ch)

	for i := 0; i < 16; i++ {
		go common.ConsumeMerchant(ch, Queryctx, rowCh, cancel)

	}

	go common.WriteToData(rowCh, cancelConsume)
	//f, data := common.MakeCVS("0")
	//defer f.Close()
	//for {
	//	select {
	//	case entity := <-ch:
	//		{
	//			data = common.CombineInfo(data, entity)
	//		}
	//	case <-ctx.Done():
	//		{
	//			fmt.Println("region finished")
	//			break
	//		}
	//	}
	//}

	//for _,entity := range entities {
	//	data =common.combineInfo(data,entity)
	//}
	select {
	case <-ConsumeCtx.Done():
		{
			break
		}
	}
}
