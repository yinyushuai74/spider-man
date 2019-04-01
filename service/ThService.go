// Copyright (c) 2012-2019 Grabtaxi Holdings PTE LTD (GRAB), All Rights Reserved. NOTICE: All information contained herein
// is, and remains the property of GRAB. The intellectual and technical concepts contained herein are confidential, proprietary
// and controlled by GRAB and may be covered by patents, patents in process, and are protected by trade secret or copyright law.
//
// You are strictly forbidden to copy, download, store (in any medium), transmit, disseminate, adapt or change this material
// in any way unless prior written permission is obtained from GRAB. Access to the source code contained herein is hereby
// forbidden to anyone except current GRAB employees or contractors with binding Confidentiality and Non-disclosure agreements
// explicitly covering such access.
//
// The copyright notice above does not evidence any actual or intended publication or disclosure of this source code,
// which includes information that is confidential and/or proprietary, and is a trade secret, of GRAB.
//
// ANY REPRODUCTION, MODIFICATION, DISTRIBUTION, PUBLIC PERFORMANCE, OR PUBLIC DISPLAY OF OR THROUGH USE OF THIS SOURCE
// CODE WITHOUT THE EXPRESS WRITTEN CONSENT OF GRAB IS STRICTLY PROHIBITED, AND IN VIOLATION OF APPLICABLE LAWS AND
// INTERNATIONAL TREATIES. THE RECEIPT OR POSSESSION OF THIS SOURCE CODE AND/OR RELATED INFORMATION DOES NOT CONVEY
// OR IMPLY ANY RIGHTS TO REPRODUCE, DISCLOSE OR DISTRIBUTE ITS CONTENTS, OR TO MANUFACTURE, USE, OR SELL ANYTHING
// THAT IT MAY DESCRIBE, IN WHOLE OR IN PART.

package service

import (
	"context"
	"spider-man/dto"
)

const name = "https://www.wongnai.com/"

func Scrap() {
	Queryctx, cancel := context.WithCancel(context.Background())
	ConsumeCtx, cancelConsume := context.WithCancel(context.Background())
	body, err := MakeRequest(name, "_api/regions.json", "_v", "5.056", "locale", "th", "knownLocation", "false")
	if err != nil {
		panic(err.Error())
	}
	regionResp := &dto.RegionResp{}
	UnmarshallRegion(body, regionResp)
	println(regionResp)
	//var regionChan chan *dto.City
	//for _,city := range regionResp.Cities {
	//	regionChan <- city
	//
	//}
	//entities :=QueryCityAll(regionResp.Cities[77])
	ch := make(chan *dto.THEntity, 20)
	rowCh := make(chan []string, 20)
	go QueryCityAllCh(regionResp.Cities, ch)

	for i := 0; i < 16; i++ {
		go ConsumeMerchant(ch, Queryctx, rowCh, cancel)

	}

	go WriteToData(rowCh, cancelConsume)
	//f, data := MakeCVS("0")
	//defer f.Close()
	//for {
	//	select {
	//	case entity := <-ch:
	//		{
	//			data = CombineInfo(data, entity)
	//		}
	//	case <-ctx.Done():
	//		{
	//			fmt.Println("region finished")
	//			break
	//		}
	//	}
	//}

	//for _,entity := range entities {
	//	data =combineInfo(data,entity)
	//}
	select {
	case <-ConsumeCtx.Done():
		{
			break
		}
	}
}
