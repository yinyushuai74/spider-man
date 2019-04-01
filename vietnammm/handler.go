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

package vietnammm

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
	"strconv"
	"sync"
)

const (
	host     = "https://www.vietnammm.com"
	takeaway = "/dat-mon-giao-tan-noi--"
	info     = "/them-thong-tin-ve-"
)

func ScrapeVN(ctx context.Context, zone string) {
	rowCh := make(chan []string)
	restaurantCh := make(chan *RestaurantName)
	listctx,listcancel := context.WithCancel(ctx)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go GetRestList(zone, restaurantCh,listcancel)
	go convertRestaurantToRow(restaurantCh,rowCh)
	go WriteToData(rowCh,zone,4000,listctx,wg)
	wg.Wait()
	fmt.Println("exit")
}


func GetRestList(area string,RestaurantCh chan *RestaurantName,listcancel context.CancelFunc) {
	Url := host + takeaway + area
	doc, err := goquery.NewDocument(Url)
	if err != nil {
		fmt.Println(err)
	}
	restaurants := doc.Find(".restaurant.grid")
	restaurants.Each(func(i int, s *goquery.Selection) {
		restaurant := &RestaurantName{}
		name := s.Find(".detailswrapper.grid-13").Find("a.restaurantname")
		href, IsExist := name.Attr("href")
		if IsExist == true {
			restaurant.Href = href
		}
		restaurant.Name = name.Text()
		style, IsExist := s.Find(".review-stars-range").Attr("style")
		if IsExist {
			restaurant.Rate = style
		}
		restaurant.Tag = s.Find("div.kitchens").Find("span").Text()
		restaurant.InfoHref = info + restaurant.Href[1:]
		if restaurant.Name=="{{RestaurantName}}" {
			return
		}
		QueryDetail(restaurant)
		RestaurantCh <- restaurant
	})
	listcancel()

}

func convertRestaurantToRow(dataCh chan *RestaurantName ,rowCh chan []string) {
	for   {
		select {
		case merchant := <-dataCh:
			for _, category := range merchant.Category {
				for _, item := range category.Item {
					row := []string{merchant.Name, merchant.Rate, merchant.Tag, merchant.Href, merchant.InfoHref, merchant.Street, merchant.Locality, merchant.AddressPic, merchant.OperateHour, category.Name, category.Desc, item.Name, item.Price, item.Desc}
					rowCh <- row
				}
			}
		}
	}
}

type RestaurantName struct {
	Name        string      `json:"name"`
	Rate        string      `json:"rate"`
	Tag         string      `json:"tag"`
	Href        string      `json:"href"`
	InfoHref    string      `json:"info_href"`
	Street      string      `json:"street"`
	Locality    string      `json:"locality"`
	AddressPic  string      `json:"address_pic"`
	OperateHour string      `json:"operate_hour"`
	Category    []*Category `json:"category"`
}

type Category struct {
	Name string  `json:"name"`
	Desc string  `json:"desc"`
	Item []*Item `json:"item"`
}

type Item struct {
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Price  string `json:"price"`
	Choose string `json:"choose"`
}

func QueryDetail(name *RestaurantName) {
	dishPage := host + name.Href
	detailPage := host + name.InfoHref
	doc, err := goquery.NewDocument(detailPage)
	if err != nil {
		fmt.Println("detailPage",err)
	}
	doc.Find(".restaurantopentimes").Find("td").Each(func(i int, selection *goquery.Selection) {
		if i == 1 {
			name.OperateHour = selection.Text()
		}
	})
	addressDoc := doc.Find(".moreinfo_address.grid-12")
	addressDoc.Find("span").Each(func(i int, selection *goquery.Selection) {
		itemprop, IsExist := selection.Attr("itemprop")
		if IsExist {
			if itemprop == "streetAddress" {
				name.Street = selection.Text()
			}
			if itemprop == "addressLocality" {
				name.Locality = selection.Text()
			}
		}
	})
	img, ifExist := addressDoc.Find("img").Attr("src")
	if ifExist {
		name.AddressPic = img
	}

	dishDoc, err := goquery.NewDocument(dishPage)
	if err != nil {
		fmt.Println("dishPage",err)
	}
	categories := make([]*Category, 0)
	dishDoc.Find(".menucard").Find(".menu-meals-group").Each(func(i int, selection *goquery.Selection) {
		id, idexist := selection.Attr("id")
		if idexist && id != "0" {
			category := &Category{}
			group := selection.Find("div.menu-meals-group-category")
			groupName := group.Find(".menu-category-head").Find("span").Text()
			category.Name = groupName
			grouDesc := group.Find(".menu-category-description").Text()
			category.Desc = grouDesc
			ItemList := make([]*Item, 0)
			selection.Find(".category-menu-meals").Find(".meal").Each(func(i int, meal *goquery.Selection) {
				item := &Item{}
				mealName := meal.Find("span.meal-name").Text()
				item.Name = mealName
				choose := meal.Find(".meal-description-choose-from").Text()
				item.Desc = choose
				mealDesc := meal.Find(".meal-description-additional-info").Text()
				item.Choose = mealDesc
				price := meal.Find("button.menu-meal-add.add-btn-icon").Find("span").Text()
				item.Price = price
				ItemList = append(ItemList, item)
			})
			category.Item = ItemList
			categories = append(categories, category)
		}
	})
	name.Category = categories
}

func MakeCVS(fileKey string) (*os.File, [][]string) {
	fileName := "merchant_" + fileKey + ".csv"
	f, err := os.Create(fileName) //创建文件
	if err != nil {
		panic(err.Error())
	}
	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	data := make([][]string, 0)
	title := []string{"Name", "Rate", "Tag", "Href", "InfoHref", "Street", "Locality", "AddressPic", "OperateHour", "CategoryName","categoryDesc","itemName","itemprice","itemDesc","itemChoose"}
	data = append(data, title)
	return f, data
}


func WriteToData(rowCh chan []string, provinceName string, rowSize int,ctx context.Context,wg *sync.WaitGroup) {
	 i:= 0
	var f *os.File
	var data [][]string
	var w *csv.Writer
	 defer wg.Done()
	 defer func() {
		 w = csv.NewWriter(f) //创建一个新的写入文件流
		 err := w.WriteAll(data)      //写入数据
		 if err != nil {
			 fmt.Println(err)
		 }
		 w.Flush()
 		 f.Close()
	 }()
	for {
		if i == 0 {
			f, data = MakeCVS(provinceName + "_" + strconv.Itoa(i/rowSize))
		}
		if i != 0 && i%rowSize == 0 {
			w = csv.NewWriter(f) //创建一个新的写入文件流
			w.WriteAll(data)      //写入数据
			w.Flush()
			f.Close()
			f, data = MakeCVS(provinceName + "_" + strconv.Itoa(i/rowSize))
		}
		select {
		case row := <-rowCh:
			{
				if row != nil {
					fmt.Println(row,strconv.Itoa(i))
					data = append(data, row)
				} else {
					return
				}

			}
		case <-ctx.Done():

			return
		}
		i++

	}

}

//func ()  {
//
//}
