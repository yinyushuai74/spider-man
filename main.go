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

package main

import (
	"flag"
	"grab-crawler/common"
)

var(
	cityID = flag.Int64("city", 217, "target city")
	mode = flag.Int64("mode", 0, "with menu or not")
	size = flag.Int("size", 5000, "total size per csv file")
	districtID = flag.Int64("district", -1, "target district")
	sortType = flag.Int64("sortType", 10, "sort type")

)

func main() {

	flag.Parse()

	switch *mode {
	case 0:
		common.ScrapeNowMerchantRate(*cityID, *size, *sortType, *districtID)
	case 1:
		common.ScrapeNowMerchant(*cityID, *size, *sortType, *districtID)
	}

	//outFunc()
	//time.Sleep(30*time.Second)
	//fmt.Println("main finished")
	//j :=0
	//	for i := 0; i < 10; i++ {
	//
	//		va,result := randomHit(100)
	//		if result {
	//			fmt.Println(va,j)
	//			j++
	//		}
	//	}
	//common.SrapNOW()
	//a,_:= common.MatchItem("https://www.wongnai.com/","restaurants/isao")
	//common.SrapNOW()
	//fmt.Println(a)
	//workHours := common.FindWorkTime("restaurants/isao")
	//fmt.Println(workHours)
	//menu := common.FindMenu("https://www.wongnai.com/","restaurants/montnomsod-dinsor")
	//fmt.Println(menu.ConvertMenu())
	//common.Vnscrap()
	//service.Scrap()
	//data := ""

	//for _, merchantDetail := range merchantDetails {
	//	rating := strconv.FormatFloat(merchantDetail.Rating, 'E', -1, 64)
	//	merchant := merchantDetail.MerchantName + ", " + merchantDetail.MerchantPhone + ", " + merchantDetail.MerchantAddress + ", " + merchantDetail.District + ", " + merchantDetail.City + ", " + merchantDetail.MerchantCategory + ", " + merchantDetail.OperatingHours + ", "+rating+", "
	//	for _,merchantDetail := range merchantDetail.MenuCategories {
	//		category := merchant + merchantDetail.MenuTitle + ", "
	//		for _,merchantDetail := range merchantDetail.Item {
	//			row := category + merchantDetail.Name + ", "+ merchantDetail.Price + ", "+ merchantDetail.FoodOrderTimes + "\n"
	//			data = data + row
	//		}
	//	}
	//}
	///////////////++++++++++++++++++/////////////////////////////////////
	//district := []string{
	//	"quan-5",
	//	"quan-6",
	//	"quan-7",
	//	"quan-8",
	//	"quan-9",
	//	"quan-10",
	//	"quan-11",
	//	"quan-12",
	//	"an-phu-thao-dien",
	//	"binh-chanh",
	//	"phu-my-hung",
	//	"quan-binh-thanh",
	//	"quan-binh-tan",
	//	"quan-go-vap",
	//	"quan-phu-nhuan",
	//	"quan-thu-duc",
	//	"quan-tan-binh",
	//	"quan-tan-phu"}
	//for _, zone := range district {
	//	vietnammm.ScrapeVN(context.Background(),zone)
	//}

}
