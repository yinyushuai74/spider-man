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
	"io/ioutil"
	"net/http"
	"fmt"
	"grab-crawler/dto"
	"encoding/json"
	"strconv"
	"os"
	"github.com/PuerkitoBio/goquery"
	"context"
	"encoding/csv"
	"regexp"
)
const web = "https://www.wongnai.com/"

func MakeRequest(name string, action string, parameters ... string) ([]byte, error) {
	web := name + action
	for n, param := range parameters {
		if n == 0 {
			web = web + "?" + param
		} else {
			if 0 == n%2 {
				web = web + "&" + param
			} else {
				web = web + "=" + param
			}
		}
	}

	resp, err := http.Get(web)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := string(body)
	fmt.Println(result)
	return body, nil
}
func UnmarshallRegion(body []byte, response *dto.RegionResp) {
	json.Unmarshal(body, response)
}

func UnmarshallMerchantResponse(body []byte, response *dto.ThMerchantResponse) {
	json.Unmarshal(body, response)

}
func QueryCity(city *dto.City, pageNum int) *dto.ThMerchantResponse {
	body, err := MakeRequest(web, "_api/businesses.json", "domain", "1", "page.number", strconv.Itoa(pageNum), "regions", strconv.Itoa(city.Id))
	if err != nil {
		panic(err.Error())

	}
	resp := &dto.ThMerchantResponse{}
	UnmarshallMerchantResponse(body, resp)
	monitor, _ := json.Marshal(resp)
	fmt.Print(string(monitor[:]))
	return resp
}
func QueryCityAll(city *dto.City) []*dto.THEntity {
	i := 1
	entities := make([]*dto.THEntity, 0)
	for {
		resp := QueryCity(city, i)
		if len(resp.THPage.Entities) == 0 {
			break
		}
		//for _, merchant := range resp.THPage.Entities {
		//	entities = append(entities, merchant)
		//}
		entities = append(entities, resp.THPage.Entities...)
		if resp.THPage.Last == resp.THPage.TotalNumberOfEntities || i == resp.THPage.TotalNumberOfPages {
			break
		}
		i++
	}
	return entities
}

func QueryCityAllCh(cities []*dto.City, ch chan *dto.THEntity) {
	for _, city := range cities {
		i := 1
		for {
			resp := QueryCity(city, i)
			if len(resp.THPage.Entities) == 0 {
				break
			}
			for _, merchant := range resp.THPage.Entities {
				ch <- merchant
			}
			if resp.THPage.Last == resp.THPage.TotalNumberOfEntities || i == resp.THPage.TotalNumberOfPages {
				break
			}
			i++
		}
		fmt.Println("city" + strconv.Itoa(city.Id))
	}
	ch <- nil
}

func CombineInfo(data [][]string, merchant *dto.THEntity) [][]string {
	workHours := FindWorkTime(web, merchant.RUrl)
	menu := FindMenu(web, merchant.RUrl)
	row := []string{strconv.Itoa(merchant.Id), merchant.DisplayName, merchant.NameOnly.English, merchant.Branch.Thai, merchant.Branch.English, merchant.Contact.CallablePhoneNo, merchant.Contact.Address.District.ConvertNameId(), merchant.Contact.Address.City.ConvertNameId(),
		merchant.Contact.Address.SubDistrict.ConvertNameId(), merchant.Contact.Address.Hint, merchant.Contact.Address.Street, convertNameIds(merchant.Categories), workHours, strconv.FormatFloat(merchant.Rating, 'E', -1, 64), menu.ConvertMenu()}
	data = append(data, row)
	return data
}

func CombineInfoCh(merchant *dto.THEntity, ch chan []string) {
	workHours := FindWorkTime(web, merchant.RUrl)
	menuResp, _ := MatchItem(web, merchant.RUrl)
	menu := ConvertResp2Categories(menuResp)
	row := []string{menu.ConvertMenu(), strconv.Itoa(merchant.Id), merchant.DisplayName, merchant.NameOnly.English, web+merchant.RUrl,strconv.FormatFloat(merchant.Lat, 'E', -1, 64), strconv.FormatFloat(merchant.Lng, 'E', -1, 64), merchant.Branch.Thai, merchant.Branch.English, merchant.Contact.CallablePhoneNo, merchant.Contact.Address.District.ConvertNameId(), merchant.Contact.Address.City.ConvertNameId(),
		merchant.Contact.Address.SubDistrict.ConvertNameId(), merchant.Contact.Address.Hint, merchant.Contact.Address.Street, convertNameIds(merchant.Categories), workHours, strconv.FormatFloat(merchant.Rating, 'E', -1, 64),}
	ch <- row
}

func ConsumeMerchant(ch chan *dto.THEntity, ctx context.Context, rowCh chan []string, cancel context.CancelFunc) {
	for {
		select {
		case entity := <-ch:
			{
				if entity != nil {
					CombineInfoCh(entity, rowCh)
				} else {
					fmt.Println("finished")
					rowCh <- nil
					cancel()
					return
				}
			}
		case <-ctx.Done():
			{
				fmt.Println("region finished")
				return
			}
		}
	}

}
func WriteToData(rowCh chan []string, cancel context.CancelFunc) {
	i := 0
	var f *os.File
	var data [][]string
	for {
		if i == 0 {
			f, data = MakeCVS("merchant" + "_" + strconv.Itoa(i/4000))
		}
		if i != 0 && i%4000 == 0 {
			w := csv.NewWriter(f) //创建一个新的写入文件流
			w.WriteAll(data)      //写入数据
			w.Flush()
			f.Close()
			f, data = MakeCVS("merchant" + "_" + strconv.Itoa(i/4000))
		}
		select {
		case row := <-rowCh:
			{
				print(i)
				if row != nil {
					category := row[0]
					categoryEntity := &dto.ThaiMenu{}
					json.Unmarshal([]byte(category), categoryEntity)
					data = WriteData(row[1:], data)
					for _, category := range categoryEntity.Categories {
						data = WriteData([]string{"", "", "","","", "", "", "","","", "", "", "", "", "", "", "", category.Name}, data)
						for _, item := range category.Items {
							data = WriteData([]string{"", "", "","","", "", "", "","", "", "", "", "", "", "", "", "", category.Name, item.Name, item.Price}, data)
						}
					}
				} else {
					w := csv.NewWriter(f) //创建一个新的写入文件流
					w.WriteAll(data)      //写入数据
					w.Flush()
					f.Close()
					cancel()
					return
				}

			}

		}
		i++

	}

}

func WriteData(row []string, data [][]string) [][]string {
	data = append(data, row)
	return data
}

func convertNameIds(kvs []*dto.NameId) string {
	result := ""
	for _, kv := range kvs {
		result = result + "{" + kv.ConvertNameId() + "}"
	}
	return result
}

func MakeCVS(fileKey string) (*os.File, [][]string) {
	fileName := "merchant_" + fileKey + ".csv"
	f, err := os.Create(fileName) //创建文件
	if err != nil {
		panic(err.Error())
	}
	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	data := make([][]string, 0)
	title := []string{"Id", "MerchantName", "EnglishName","Page","lat","lng", "Branch", "BranchEnglish", "MerchantPhone", "District", "City", "SubDistrict", "Hint", "Street", "MerchantCategory", "WorkingHoursMessage", "Stars", "MenuTitle", "ItemName", "ItemPrice"}
	data = append(data, title)
	return f, data
}
func FindWorkTime(name string, url string) string {
	webUrl := name + url
	doc, err := goquery.NewDocument(webUrl)
	if err != nil {
		panic(err.Error())
	}
	workHour := ""
	sel := doc.Find(".BusinessDetailBlock__TwoColTable-eelvf4-1.jfqPXy")
	sel.Find("td").Each(func(i int, selection *goquery.Selection) {
		workHour = workHour + selection.Text() + "  "
	})
	return workHour
}

func MatchItem(name string, url string) (*dto.MatchResponse, error) {
	webUrl := name + url + "/menu"
	resp, err := http.Get(webUrl)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	doc := string(body)
	a := "(?:\"businessMenu\".*?value\".)(.*_q\".*?}})"
	reg := regexp.MustCompile(a)
	doc2 := string(reg.Find([]byte(doc)))
	b := "{\"menuGroups\".*"
	reg = regexp.MustCompile(b)
	doc3 := string(reg.Find([]byte(doc2)))
	fmt.Println(doc3)
	response := &dto.MatchResponse{}
	json.Unmarshal([]byte(doc3), response)
	return response, nil
}

func ConvertResp2Categories(resp *dto.MatchResponse) *dto.ThaiMenu {
	target := make([]*dto.Category, 0)
	for _, category := range resp.MenuGroups {
		toItemes := ConvertMatchItem2Item(category.Items)
		to := &dto.Category{
			Name:  category.Name,
			Items: toItemes,
		}
		target = append(target, to)
	}
	result := &dto.ThaiMenu{
		Categories: target,
	}
	return result
}
func ConvertMatchItem2Item(source []*dto.MatchItem) []*dto.ThaiItem {
	target := make([]*dto.ThaiItem, 0)
	for _, from := range source {
		to := &dto.ThaiItem{
			Name:        from.Name,
			Price:       from.Price.Text,
			DisplayName: from.DisplayName,
			Description: from.Description,
		}
		target = append(target, to)

	}
	return target

}

func FindMenu(name string, url string) *dto.ThaiMenu {
	webUrl := name + url + "/menu"
	doc, err := goquery.NewDocument(webUrl)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(url)
	//menuDiv := doc.Find(".BusinessMenuPage__MenuBlockContainer-s1ftsjdm-1.bPcnyY")
	categories := make([]*dto.Category, 0)
	labels := doc.Find(".label")
	labels.Each(func(i int, selection *goquery.Selection) {
		categoryName := selection.Find(".BusinessMenuList__DropdownLabelText-i31xd8-1.AkJkT").Text()
		category := &dto.Category{
			Name: categoryName,
		}
		categories = append(categories, category)
	})
	categoryItems := make([][]*dto.ThaiItem, 0)
	bodes := doc.Find(".body")
	bodes.Each(func(i int, selection *goquery.Selection) {
		itemList := make([]*dto.ThaiItem, 0)
		itemDivs := selection.Find(".BusinessMenuItem__MenuListContainer-s5dc8ij-0.idmSdm")
		itemDivs.Each(func(i int, itemDiv *goquery.Selection) {
			name := itemDiv.Find(".BusinessMenuItem__MenuListDetailContainer-s5dc8ij-2.gMhxgR").Find(".BusinessMenuItem__MenuTitleText-s5dc8ij-3.BpBEe").Text()
			price := itemDiv.Find(".BusinessMenuItem__MenuListPriceContainer-s5dc8ij-6.hTmyma").Find(".BusinessMenuItem__MenuTitleText-s5dc8ij-3.BpBEe").Text()
			item := &dto.ThaiItem{
				Name:  name,
				Price: price,
			}
			itemList = append(itemList, item)
		})
		categoryItems = append(categoryItems, itemList)
	})
	for i, category := range categories {
		category.Items = categoryItems[i]
	}
	thaiMenu := &dto.ThaiMenu{
		Categories: categories,
	}
	return thaiMenu
}
