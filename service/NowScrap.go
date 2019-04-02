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
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"spider-man/dto"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

const (
	NOW         = "https://gappapi.deliverynow.vn"
	metadata    = "/api/meta/get_metadata"
	deliveryIDs = "/api/delivery/search_delivery_ids"
	getDetail   = "/api/delivery/get_detail?id_type=2&request_id=%d"
	getDishes   = "/api/dish/get_delivery_dishes?id_type=2&request_id=%d"
)

//func doRequestWithHeader(district string) ([]*dto.MerchantDetailTotal, error) {
//	web := "https://www.now.in.th"
//	url := web + "/Offer/LoadMore?category=0&pageIndex=0&pageSize=400&district=" + district
//	client := &http.Client{}
//	request, _ := http.NewRequest("GET", url, nil)
//	request.Header.Add("Referer", "https://www.now.in.th/nonthaburi")
//	request.Header.Add("Cookie", "ASP.NET_SessionId=xwqgttdeosnt124ncmdafevq; mroot=1; _ga=GA1.3.1857299294.1527145648; view=grid; flg=en; ilg=0; _gid=GA1.3.1063497128.1527499659; floc=716; _gat=1")
//	request.Header.Add("Host", "www.now.in.th")
//	response, _ := client.Do(request)
//	body, err := ioutil.ReadAll(response.Body)
//	//fmt.Println(string(body))
//	if err != nil {
//		return nil, err
//	}
//	merchantResp := &dto.TotalPageResp{}
//	fmt.Println(string(body))
//	json.Unmarshal(body, merchantResp)
//	//fmt.Println(merchantResp)
//	merchantDetailList, err := fetchMerchantItem(merchantResp, web)
//	if err != nil {
//		return nil, err
//	}
//	fmt.Println(merchantDetailList)
//	return merchantDetailList, nil
//
//}

func generateNOWRequest(method string, action string, body interface{}) []byte {
	url := NOW + action
	client := &http.Client{}
	reqJSON, err := json.Marshal(body)
	if err != nil {
		return nil
	}
	reqReader := strings.NewReader(string(reqJSON))
	request, _ := http.NewRequest(method, url, reqReader)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("x-foody-api-version", "1")
	request.Header.Add("x-foody-app-type", "1004")
	request.Header.Add("x-foody-client-language", "en")
	request.Header.Add("x-foody-client-type", "1")
	request.Header.Add("x-foody-client-version", "1.8.3")
	request.Header.Add("x-foody-client-id", "")
	resp, err := client.Do(request)
	if err != nil {
		return nil
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	return respBody

}

func getMetaData() []*dto.MetaCity {
	url := NOW + metadata
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("x-foody-api-version", "1")
	request.Header.Add("x-foody-app-type", "1004")
	request.Header.Add("x-foody-client-language", "vi")
	request.Header.Add("x-foody-client-type", "1")
	request.Header.Add("x-foody-client-version", "1.8.3")
	request.Header.Add("x-foody-client-id", "")
	resp, err := client.Do(request)
	if err != nil {
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	MetaResp := &dto.MeteData{}
	err = json.Unmarshal(body, MetaResp)
	if err != nil {
		fmt.Println(err)
	}
	return MetaResp.Reply.Country.Cities
}

func getMerchantIDs(cityID int64, districts []*dto.MetaDistrict, sortType int64) []int64 {
	merchantIDs := make([]int64, 0)

	for _, district := range districts {
		req := &dto.DeliveryIDsRequest{
			CategoryGroup: 1,
			CityID:        cityID,
			DeliveryOnly:  true,
			FoodyServices: []int64{1},
			Keyword:       "",
			SortType:      sortType,
			DistrictIds:   []int64{district.DistrictID},
		}
		respReader := generateNOWRequest("POST", deliveryIDs, req)
		deliveryIDs := &dto.DeliveryIDsResp{}
		json.Unmarshal(respReader, deliveryIDs)
		merchantIDs = append(merchantIDs, deliveryIDs.Reply.DeliveryIds...)
	}

	return merchantIDs

}

func PutMerchantID(cityID int64, idChan chan int64, sortType int64, districtID int64, cancel context.CancelFunc) (*dto.MetaCity, int) {
	citys := getMetaData()
	var meteCity *dto.MetaCity
	var districts []*dto.MetaDistrict
	for _, city := range citys {
		if city.Id == cityID {
			if districtID == -1 {
				districts = city.Districts
			} else {
				for _, v := range city.Districts {
					if v.DistrictID == districtID {
						districts = []*dto.MetaDistrict{v}
						break
					}
				}
				if districts == nil {
					return &dto.MetaCity{}, 0
				}
			}
			meteCity = city
			break
		}
	}
	//merchantIDs := getMerchantLocal()
	merchantIDs := getMerchantIDs(cityID, districts, sortType)
	fmt.Println("get MerchantIDs done")
	size := len(merchantIDs)
	go func(ids []int64) {
		for _, id := range merchantIDs {
			idChan <- id
		}
		fmt.Println("put merchantID done")
		cancel()
	}(merchantIDs)
	return meteCity, size
}

func GetMerchantDetail(id int64, ctx context.Context, detailChan chan dto.DetailTotal) {
	restDishes := &dto.DishResp{}
	restDetail := &dto.DeliveryDetailResp{}
	respBody := generateNOWRequest("GET", fmt.Sprintf(getDetail, id), nil)
	json.Unmarshal(respBody, restDetail)
	respBody = generateNOWRequest("GET", fmt.Sprintf(getDishes, id), nil)
	json.Unmarshal(respBody, restDishes)
	if restDetail.Reply == nil || restDishes.Reply == nil {
		return
	}
	if restDishes.Reply.MenuInfos == nil || restDetail.Reply.DeliveryDetail == nil {
		return
	}
	DetailTotal := dto.DetailTotal{
		RestaurantMenu: restDishes.Reply.MenuInfos,
		DeliveryDetail: restDetail.Reply.DeliveryDetail,
	}
	detailChan <- DetailTotal
}

func GetMerchantDetailForRating(id int64, detailChan chan dto.DetailTotal) {
	restDetail := &dto.DeliveryDetailResp{}
	respBody := generateNOWRequest("GET", fmt.Sprintf(getDetail, id), nil)
	json.Unmarshal(respBody, restDetail)
	if restDetail.Reply == nil || restDetail.Reply.DeliveryDetail == nil {
		return
	}
	DetailTotal := dto.DetailTotal{
		DeliveryDetail: restDetail.Reply.DeliveryDetail,
	}
	detailChan <- DetailTotal
}
func convertInfoWithRating(detaiChan chan dto.DetailTotal, provinceName string, ctx context.Context, merchantSize int, itemSize int) {
	var err error
	var num = 0
	var itemNum = 0
	var f *os.File
	var data [][]string
	var w *csv.Writer
	var rate = "0"
	var reviewSize = "0"
	for {
		select {
		case <-ctx.Done():
			w = csv.NewWriter(f) //创建一个新的写入文件流
			w.WriteAll(data)     //写入数据
			if err != nil {
				fmt.Println(err)
			}
			w.Flush()
			f.Close()
			return
		case detail := <-detaiChan:
			num++
			if detail.DeliveryDetail != nil && detail.DeliveryDetail.Rating != nil {
				rate = fmt.Sprintf("%f", detail.DeliveryDetail.Rating.Avg)
				reviewSize = formatInt64(detail.DeliveryDetail.Rating.TotalReview)
			}
			promoByte, err := json.Marshal(detail.DeliveryDetail.Delivery.Promotions)
			if err != nil {
				fmt.Println("marshal promotion err")
			}
			deliveryFee, err := json.Marshal(detail.DeliveryDetail.Delivery.ShippingFee)
			if err != nil {
				fmt.Println("marshal promotion err")
			}
			fmt.Println(string(deliveryFee))
			row := []string{
				formatInt64(detail.DeliveryDetail.DeliveryId),
				provinceName,
				detail.DeliveryDetail.Name,
				detail.DeliveryDetail.Address,
				formatFloat64(detail.DeliveryDetail.Position.Latitude),
				formatFloat64(detail.DeliveryDetail.Position.Longitude),
				"",
				rate,
				reviewSize,
				strconv.FormatBool(detail.DeliveryDetail.Position.IsVerified),
				string(promoByte),
				string(deliveryFee),

			}
			if itemNum == 0 {
				f, data = MakeCVSNewRate(provinceName + "_" + strconv.Itoa(itemNum/itemSize))
			}
			if itemNum != 0 && itemNum%itemSize == 0 {
				w = csv.NewWriter(f) //创建一个新的写入文件流
				w.WriteAll(data)     //写入数据
				w.Flush()
				f.Close()
				f, data = MakeCVSNewRate(provinceName + "_" + strconv.Itoa(itemNum/itemSize))
			}
			itemNum++
			fmt.Println(row)
			data = append(data, row)
		}
		fmt.Println("merchant : %d", num)
		if num == merchantSize {
			w = csv.NewWriter(f) //创建一个新的写入文件流
			w.WriteAll(data)     //写入数据
			if err != nil {
				fmt.Println(err)
			}
			w.Flush()
			f.Close()
			return
		}
	}

}

func MakeCVSNewRate(fileKey string) (*os.File, [][]string) {
	fileName := "merchant_" + fileKey + ".csv"
	f, err := os.Create(fileName) //创建文件
	if err != nil {
		panic(err.Error())
	}
	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	data := make([][]string, 0)
	title := []string{"MerchantID", "city", "name", "address", "latitude", "logitude", "verified", "rate", "reviewSize", "isVerified","Promotion", "delivery"}
	data = append(data, title)
	return f, data
}

func ScrapeNowMerchant(cityId int64, rows int, sortType int64, districtID int64) {
	idChan := make(chan int64)
	//threadSize := make(chan bool, 30)
	detailChan := make(chan dto.DetailTotal, 50)
	ctx, cancel := context.WithCancel(context.Background())
	city, size := PutMerchantID(cityId, idChan, sortType, districtID, cancel)
	fmt.Println("get MeteData done")
	go func(idChan chan int64) {
		for i := 0; i < 40; i++ {
			for {
				id := <-idChan
				GetMerchantDetail(id, context.Background(), detailChan)
			}
		}
	}(idChan)
	convertInfo(detailChan, city.Name, ctx, size, rows)
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func ScrapeNowMerchantRate(cityId int64, rows int, sortType int64, districtID int64) {
	idChan := make(chan int64)
	//threadSize := make(chan bool, 30)
	detailChan := make(chan dto.DetailTotal, 50)
	ctx, cancel := context.WithCancel(context.Background())
	city, size := PutMerchantID(cityId, idChan, sortType, districtID, cancel)
	fmt.Println("get MeteData done")
	wg := sync.WaitGroup{}
	go func(idChan chan int64) {
		for i := 0; i < 50; i++ {
			wg.Add(1)
			for {
				id := <-idChan
				GetMerchantDetailForRating(id, detailChan)
			}
		}
	}(idChan)
	convertInfoWithRating(detailChan, city.Name, ctx, size, rows)
	wg.Wait()
}

func convertInfo(detaiChan chan dto.DetailTotal, provinceName string, ctx context.Context, merchantSize int, itemSize int) {
	var err error
	var num = 0
	var itemNum = 0
	var f *os.File
	var data [][]string
	var w *csv.Writer
	for {
		select {
		case <-ctx.Done():
			w = csv.NewWriter(f) //创建一个新的写入文件流
			w.WriteAll(data)     //写入数据
			if err != nil {
				fmt.Println(err)
			}
			w.Flush()
			f.Close()
			return
		case detail := <-detaiChan:
			num++
			for _, category := range detail.RestaurantMenu {
				for _, dish := range category.Dishes {
					timeJSON, err := json.Marshal(detail.DeliveryDetail.Delivery.Time)
					if err != nil {
						fmt.Println(err)
					}
					restCategory, err := json.Marshal(detail.DeliveryDetail.Category)
					if err != nil {
						fmt.Println(err)
					}
					cuisines, err := json.Marshal(detail.DeliveryDetail.Cuisines)
					if err != nil {
						fmt.Println(err)
					}
					row := []string{
						formatInt64(detail.DeliveryDetail.DeliveryId),
						detail.DeliveryDetail.Name,
						detail.DeliveryDetail.Url,
						provinceName,
						detail.DeliveryDetail.Address,
						string(restCategory),
						string(cuisines),
						fmt.Sprintf("{%f,%f}", detail.DeliveryDetail.Position.Latitude, detail.DeliveryDetail.Position.Longitude),
						fmt.Sprintf("{%f~%f}", detail.DeliveryDetail.PriceRange.MinPrice, detail.DeliveryDetail.PriceRange.MaxPrice),
						"",
						string(timeJSON),
						category.DishTypeName,
						dish.Name,
						dish.Price.Text,
						formatInt64(dish.TotalOrder),
						dish.Description,
					}

					if detail.DeliveryDetail != nil && detail.DeliveryDetail.Rating != nil {
						row[9] = formatFloat64(detail.DeliveryDetail.Rating.Avg)
					}
					if itemNum == 0 {
						f, data = MakeCVSNew(provinceName + "_" + strconv.Itoa(itemNum/itemSize))
					}
					if itemNum != 0 && itemNum%itemSize == 0 {
						w = csv.NewWriter(f) //创建一个新的写入文件流
						w.WriteAll(data)     //写入数据
						w.Flush()
						f.Close()
						f, data = MakeCVSNew(provinceName + "_" + strconv.Itoa(itemNum/itemSize))
					}
					itemNum++
					fmt.Println(row)
					data = append(data, row)
				}
			}
		}
		fmt.Println("merchant : %d", num)
	}

}

func MakeCVSNew(fileKey string) (*os.File, [][]string) {
	fileName := "merchant_" + fileKey + ".csv"
	f, err := os.Create(fileName) //创建文件
	if err != nil {
		panic(err.Error())
	}
	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	data := make([][]string, 0)
	title := []string{"MerchantID", "Name", "Url", "cityName", "address", "Merchant_Category", "Merchant_Cuisines", "Position", "Price_Range",
		"Rating", "Operating_Hour", "Dish_Type", "Dish_Name", "Dish_Price", "Order_Times", "Desc"}
	data = append(data, title)
	return f, data
}

func formatInt64(variable int64) string {
	return strconv.FormatInt(variable, 10)
}

func formatFloat64(variable float64) string {
	return fmt.Sprintf("%f", variable)
}

func doPostRequest(provinceId int, ch chan *dto.MerchantTotal, district []*dto.DistrictResult) {

	for _, district := range district {
		i := 1
		for {
			web := NOW
			url := web + "/List/GetListDeliveryByFilter"

			client := &http.Client{}
			detailCh := ch
			requestBody := fmt.Sprintf("{\"filters\":{\"Keyword\":null,\"CategoryIds\":null,\"DistrictIds\":[%d],\"ProvinceId\":%d,\"CuisineIds\":null,\"SortType\":11,\"PageSize\":30,\"Lat\":39.9834503,\"Long\":116.32064899999999,\"PageIndex\":%d,\"MasterCategories\":[1]}}", district.Id, provinceId, i)
			postBody := strings.NewReader(requestBody)
			request, _ := http.NewRequest("POST", url, postBody)
			request.Header.Add("Content-Type", "application/json")
			response, _ := client.Do(request)
			body, err := ioutil.ReadAll(response.Body)
			//fmt.Println(string(body))
			if err != nil {
				return
			}
			merchantResp := &dto.TotalPageResp{}
			fmt.Println(string(body))
			json.Unmarshal(body, merchantResp)
			fetchMerchantItem(merchantResp, detailCh, district)
			if len(merchantResp.Result.MerchantTotalList) < 30 {
				break
			}
			i++
		}
	}
	ch <- nil

	//fmt.Println(merchantResp)
}

func getDistrict() *dto.NOWDistrict {
	web := NOW + "/List/GetDistrictForFilter"
	client := &http.Client{}
	request, _ := http.NewRequest("POST", web, nil)
	request.Header.Add("Cookie", "ASP.NET_SessionId=ffv4sgl0oc3ke3wslm2i1hkt; mroot=1; _ga=GA1.2.448560280.1527672709; view=grid; _gid=GA1.2.1482463936.1529379066; DELIVERY.AUTH.UDID=14964d8c-6793-417c-a1b1-054d34cc98a4; DELIVERY.AUTH=F69FF82B7CA63DF21CA08BB58D0588BDD871D23017FC27CB6D95A3320FCE52E668DDC4B1EBEB00D66E7E13F078B949BED1ED5DD3E9ECB3FE5A70B350A325808372D5F956A3AE4567ECD1940912345F23C4F9360A355E9F6ED8519F44C1D86DE7431ED56D30F2140DF5459B134DFF1552426446B6201546FA5504D44AA925333EDA3EBCCF0A4C028BB2A2A40F934FB47B07F386D05B7B53C8F7B7C908F35687E1B1D30F6228AAA30903DEC4964B4060F0476095A786AF71DC11B4BDE5542CA172E285AFCD5544D46D5223655DBBDFD8DD153798194E11C4AE7B11BBE399A67DD2E609932D38CE243CDC46BDEB3314830762367AF8E08CF388F9A1EAF9806C6D30; ilg=1; flg=en; floc=217; _gat=1")
	resp, _ := client.Do(request)
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	if err != nil {
		return nil
	}
	districtresp := &dto.NOWDistrict{}
	json.Unmarshal(body, districtresp)
	return districtresp

}

//deprecated
func SrapNOW() {
	districtmap := make(map[int]string, 50)
	districtresp := getDistrict()
	districts := districtresp.Result
	for _, district := range districts {
		districtmap[district.Id] = district.AsciiName

	}
	web := NOW
	province := make(map[string]string, 2)
	provinceId := 217
	//province["217"] = "ho-chi-minh"
	//province["218"] = "ha-noi"
	//province[strconv.Itoa(provinceId)] = "da-nang"
	//province["6"] = "District 4"
	//province["7"] = "District 5"
	//province["8"] = "District 6"
	//province["9"] = "District 7"
	//province["10"] = "District 8"
	//province["11"] = "District 9"
	//province["12"] = "District 10"
	//province["13"] = "District 11"
	//province["15"] = "Bình Thạnh District"
	//province["16"] = "Tân Bình District"
	//province["17"] = "Phú Nhuận District"
	//province["19"] = "Tân Phú District"
	//province["2"] = "Gò Vấp District"
	//province["18"] = "Binh Tan District"
	//province["693"] = "Thu Duc District"
	//province["696"] = "Binh Chanh"
	dataCh := make(chan []string, 30)
	cont, cancel := context.WithCancel(context.Background())
	merchantCh := make(chan *dto.MerchantDetailTotal, 30)

	ch := make(chan *dto.MerchantTotal, 30)
	go doPostRequest(provinceId, ch, districts)
	for i := 0; i < 30; i++ {

		go QueryMore(ch, merchantCh, web)
	}
	go QueryAll2Ch(merchantCh, dataCh, strconv.Itoa(provinceId), province, districtmap)
	WriteToDataNow(dataCh, cancel, province[strconv.Itoa(provinceId)])

	select {
	case <-cont.Done():
		{
			return
		}
	}
}

//
//func doRequest(district string) ([]*dto.MerchantDetailTotal, error) {
//	web := "https://www.now.in.th"
//	url := web + "/Offer/LoadMore?category=0&pageIndex=0&pageSize=400&district=" + district
//	resp, err := http.Get(url)
//	if err != nil {
//		return nil, err
//	}
//	body, err := ioutil.ReadAll(resp.Body)
//	//fmt.Println(string(body))
//	if err != nil {
//		return nil, err
//	}
//	merchantResp := &dto.TotalPageResp{}
//	fmt.Println(string(body))
//	json.Unmarshal(body, merchantResp)
//	//fmt.Println(merchantResp)
//	merchantDetailList, err := fetchMerchantItem(merchantResp, web)
//	if err != nil {
//		return nil, err
//	}
//	fmt.Println(merchantDetailList)
//	return merchantDetailList, nil
//}

func convertResp(merchant *dto.Merchant) *dto.MerchantDetail {
	merchantDetail := &dto.MerchantDetail{
		MerchantName:     merchant.Restaurant.ResName,
		MerchantPhone:    merchant.Restaurant.Phone,
		MerchantAddress:  merchant.Restaurant.FullAddress,
		District:         merchant.Restaurant.DistrictId,
		City:             merchant.Restaurant.RestaurantLocation,
		MerchantCategory: merchant.AttributeString,
	}
	return merchantDetail

}

func convertRespTotal(merchant *dto.MerchantTotal) *dto.MerchantDetailTotal {
	merchantDetailTotal := &dto.MerchantDetailTotal{
		MerchantName:     merchant.ResName,
		MerchantAddress:  merchant.FullAddress,
		District:         merchant.DistrictId,
		MerchantCategory: merchant.Category,
		PriceRange:       merchant.PriceRange,
		Id:               merchant.Id,
		OperatingHours:   merchant.TimeRange,
		Url:              NOW + merchant.DetailUrl,
		Lat:              merchant.ResLat,
		Lng:              merchant.ResLng,
	}
	return merchantDetailTotal

}

func fetchMerchantItem(merchantResp *dto.TotalPageResp, ch chan *dto.MerchantTotal, district *dto.DistrictResult) {
	merchantCh := ch
	for _, merchant := range merchantResp.Result.MerchantTotalList {
		merchant.DistrictId = district.Id
		merchantCh <- merchant
	}
}

func QueryMore(merchantCh chan *dto.MerchantTotal, detailChan chan *dto.MerchantDetailTotal, web string) {
	for {
		select {
		case merchant := <-merchantCh:
			{
				defer func() {
					if p := recover(); p != nil {
						fmt.Printf("panic recover! p: %v", p)
					}
				}()
				if merchant == nil {
					detailChan <- nil
				} else {
					merchantDetail := convertRespTotal(merchant)
					merchandiseUrl := web + merchant.DetailUrl
					doc, err := goquery.NewDocument(merchandiseUrl)
					menuCategories := make([]*dto.MenuCategory, 0)
					var rating = 0.0
					starsDoc := doc.Find(".rating")
					if starsDoc != nil {
						stars := starsDoc.Find(".stars")
						stars.Find(".full").Each(func(i int, selection *goquery.Selection) {
							rating = rating + 1
						})
						stars.Find(".half").Each(func(i int, selection *goquery.Selection) {
							rating = rating + 0.5
						})
						merchantDetail.Rating = rating
					}

					doc.Find(".scrollspy").Each(func(i int, s *goquery.Selection) {
						menuCategory := &dto.MenuCategory{}
						item := make([]*dto.MerchantItem, 0)
						s.Find(".title-kind-food").Each(func(i int, s *goquery.Selection) {
							menuTitle := strings.TrimSpace(s.Text())
							menuCategory.MenuTitle = menuTitle
							//fmt.Println(menuTitle)
						})
						s.Find(".box-menu-detail.clearfix").Each(func(i int, s *goquery.Selection) {
							//name
							foodNameUnFormat := strings.TrimSpace(s.Find(".title-name-food").Text())

							foodName := convertFoodName(foodNameUnFormat)

							foodOrderTimes := s.Find(".font11.light-grey").Find("span").Find(".bold").Text()

							foodDescribe := s.Find(".desc").Text()
							//price
							var price1 *html.Node
							price := ""
							priceNode := s.Find("p").Find(".txt-blue.font16.bold")
							if priceNode != nil && len(priceNode.Nodes) != 0 {
								price1 = priceNode.Get(0)
								if price1 != nil {
									price = price1.FirstChild.Data

								}
							}
							//fmt.Println(price)
							merchantItem := &dto.MerchantItem{
								Name:           foodName,
								Price:          price,
								FoodOrderTimes: foodOrderTimes,
								Describe:       foodDescribe,
							}
							item = append(item, merchantItem)
						})
						menuCategory.Item = item
						menuCategories = append(menuCategories, menuCategory)
					})
					merchantDetail.MenuCategories = menuCategories
					merchantDetail.OperatingHours, err = queryTime(err, merchandiseUrl)
					data, err := json.Marshal(merchantDetail)
					fmt.Println(string(data))
					detailChan <- merchantDetail
				}
			}

		}

	}

}
func convertFoodName(from string) string {
	to := make([]byte, 0)
	foodNameByte := []byte(from)
	j := 0
	for i, e := range foodNameByte {
		if e == 10 || (e == 32 && (foodNameByte[i+1] == 10 || foodNameByte[i+1] == 32)) {
			break
		}
		j++
		to = append(to, e)
	}
	result := make([]byte, j)
	for i := 0; i < j; i++ {
		result[i] = to[i]
	}
	//fmt.Println(string(result))
	return string(result)
}
func queryTime(err error, merchandiseUrl string) (string, error) {
	resp, err := http.Get(merchandiseUrl)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	context := string(body)
	a := "5px;\\\">\\d{2}:.* -.*:.*</span>"
	reg := regexp.MustCompile(a)
	context1 := string(reg.Find([]byte(context)))
	b := "\\d{2}:.* -.*:.{2}"
	reg2 := regexp.MustCompile(b)
	times := string(reg2.Find([]byte(context1)))
	return times, nil
}

func QueryAll2Ch(merchantDetailCh chan *dto.MerchantDetailTotal, dataCh chan []string, provinceId string, province map[string]string, districtMap map[int]string) {
	i := 0
	for {
		select {
		case merchantDetail := <-merchantDetailCh:
			if merchantDetail == nil {
				dataCh <- nil
				break
			}
			println("write" + strconv.Itoa(i))
			i++
			for _, menuCategory := range merchantDetail.MenuCategories {
				for _, item := range menuCategory.Item {
					row := []string{merchantDetail.MerchantName, merchantDetail.MerchantPhone, merchantDetail.Url, strconv.FormatFloat(merchantDetail.Lat, 'E', -1, 64), strconv.FormatFloat(merchantDetail.Lng, 'E', -1, 64), merchantDetail.MerchantAddress, districtMap[merchantDetail.District], province[provinceId], merchantDetail.MerchantCategory, merchantDetail.PriceRange, merchantDetail.OperatingHours, strconv.FormatFloat(merchantDetail.Rating, 'E', -1, 64), menuCategory.MenuTitle, item.Name, item.Price, item.FoodOrderTimes, item.Describe}
					dataCh <- row
				}
			}
		}
	}
}

//func ProduceMerchant(districtValue string, merchantCh chan *dto.MerchantDetailTotal) {
//	i := 0
//	for {
//		merchantDetails, _, err := doPostRequest(districtValue, i)
//		if err != nil {
//			fmt.Println(err)
//			i++
//			continue
//		}
//		for _, merchant := range merchantDetails {
//			merchantCh <- merchant
//		}
//		if len(merchantDetails) < 30 {
//			return
//		}
//	}
//}

func MakeCVSNow(fileKey string) (*os.File, [][]string) {
	fileName := "merchant_" + fileKey + ".csv"
	f, err := os.Create(fileName) //创建文件
	if err != nil {
		panic(err.Error())
	}
	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	data := make([][]string, 0)
	title := []string{"MerchantName", "MerchantPhone", "URL", "Lat", "Lng", "MerchantAddress", "District", "City", "MerchantCategory", "PriceRange", "OperatingHours", "Rating", "MenuTitle", "ItemName", "Price", "OrderTimes", "Desc"}
	data = append(data, title)
	return f, data
}

func WriteToDataNow(rowCh chan []string, cancel context.CancelFunc, provinceName string) {
	i := 0
	var f *os.File
	var data [][]string
	for {
		if i == 0 {
			f, data = MakeCVSNow(provinceName + "_" + strconv.Itoa(i/4000))
		}
		if i != 0 && i%4000 == 0 {
			w := csv.NewWriter(f) //创建一个新的写入文件流
			w.WriteAll(data)      //写入数据
			w.Flush()
			f.Close()
			f, data = MakeCVSNow(provinceName + "_" + strconv.Itoa(i/4000))
		}
		select {
		case row := <-rowCh:
			{
				if row != nil {
					data = append(data, row)
				} else {
					w := csv.NewWriter(f) //创建一个新的写入文件流
					w.WriteAll(data)      //写入数据
					w.Flush()
					f.Close()
					cancel()
					println("writedata exit")
					return
				}

			}

		}
		i++

	}

}

func WriteToDataNew(row []string, provinceName string, itemSize int) {
	i := 0
	var f *os.File
	var data [][]string
	for {
		if i == 0 {
			f, data = MakeCVSNow(provinceName + "_" + strconv.Itoa(i/itemSize))
		}
		if i != 0 && i%itemSize == 0 {
			w := csv.NewWriter(f) //创建一个新的写入文件流
			w.WriteAll(data)      //写入数据
			w.Flush()
			f.Close()
			f, data = MakeCVSNow(provinceName + "_" + strconv.Itoa(i/itemSize))
		}
		if row != nil {
			data = append(data, row)
		} else {
			w := csv.NewWriter(f) //创建一个新的写入文件流
			w.WriteAll(data)      //写入数据
			w.Flush()
			f.Close()
			println("writedata exit")
			return
		}

		i++

	}

}

//
//func Vnscrap() {
//	district := make(map[string]string, 19)
//	//district["1"] = "bangkok"
//	district["716"] = "nonthaburi"
//	//district["5"] = "District 3"
//	//district["6"] = "District 4"
//	//district["7"] = "District 5"
//	//district["8"] = "District 6"
//	//district["9"] = "District 7"
//	//district["10"] = "District 8"
//	//district["11"] = "District 9"
//	//district["12"] = "District 10"
//	//district["13"] = "District 11"
//	//district["15"] = "Bình Thạnh District"
//	//district["16"] = "Tân Bình District"
//	//district["17"] = "Phú Nhuận District"
//	//district["19"] = "Tân Phú District"
//	//district["2"] = "Gò Vấp District"
//	//district["18"] = "Binh Tan District"
//	//district["693"] = "Thu Duc District"
//	//district["696"] = "Binh Chanh"
//	var f io.Writer
//	var data [][]string
//	dataCh := make(chan []string, 30)
//	for districtValue := range district {
//		dataCh := make(chan []string, 30)
//		go QueryAll2Ch(districtValue, dataCh)
//
//	}
//
//}
