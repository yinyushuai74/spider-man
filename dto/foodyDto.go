// Copyright (c) 2012-2018 Grabtaxi Holdings PTE LTD (GRAB), All Rights Reserved. NOTICE: All information contained herein
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

package dto

type LocationResponse struct {
	AllDistricts []*FDistrict `json:"AllDistricts"`
	AllLocations []*FProvince `json:"AllLocations"`
}

type FDistrict struct {
	Id     int64  `json:"Id"`
	Name   string `json:"Name"`
	CityId string `json:"CityId"`
}

type FProvince struct {
	CountryId   int64        `json:"CountryId"`
	CountryName string       `json:"CountryName"`
	DisplayName string       `json:"DisplayName"`
	Id          int64        `json:"Id"`
	Url         string       `json:"Url"`
	Districts   []*FDistrict `json:"Districts"`
}

type MeteData struct {
	Reply *MetaReply `json:"reply"`
}

type MetaReply struct {
	Country *MetaCountry `json:"country"`
}

type MetaCountry struct {
	Cities []*MetaCity `json:"cities"`
}

type MetaCity struct {
	Id        int64           `json:"id"`
	Name      string          `json:"name"`
	Districts []*MetaDistrict `json:"districts"`
}
type MetaDistrict struct {
	DistrictID     int64  `json:"district_id"`
	Name           string `json:"name"`
	UrlRewriteName string `json:"url_rewrite_name"`
}
type DeliveryIDsRequest struct {
	CategoryGroup int64   `json:"category_group"`
	CityID        int64   `json:"city_id"`
	DeliveryOnly  bool    `json:"delivery_only"`
	FoodyServices []int64 `json:"foody_services"`
	Keyword       string  `json:"keyword"`
	SortType      int64   `json:"sort_type"`
	DistrictIds   []int64 `json:"district_ids"`
}
type DeliveryIDsResp struct {
	Reply  *DeliveryIdReply `json:"reply"`
	Result string           `json:"result"`
}
type DeliveryIdReply struct {
	DeliveryIds []int64 `json:"delivery_ids"`
}

type DeliveryDetailResp struct {
	Reply *DetailReply `json:"reply"`
}

type DetailReply struct {
	DeliveryDetail *DeliveryDetail `json:"delivery_detail"`
}
type DetailPosition struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
type PriceRange struct {
	MinPrice float64 `json:"min_price"`
	MaxPrice float64 `json:"max_price"`
}

type Rating struct {
	Avg float64 `json:"avg"`
	TotalReview int64 `json:"total_review"`
}

type Delivery struct {
	Time *Time `json:"time"`
}

type Time struct {
	WeekDays []*WeekDays `json:"week_days"`
}

type WeekDays struct {
	StartTime string `json:"start_time"`
	WeekDay   int64  `json:"week_day"`
	EndTime   string `json:"end_time"`
}

type Brand struct {
	BrandId int64  `json:"brand_id"`
	Name    string `json:"name"`
}

type DeliveryDetail struct {
	Address     string          `json:"address"`
	DeliveryId  int64           `json:"delivery_id"`
	CityId      int64           `json:"city_id"`
	Name        string          `json:"name"`
	Category    []string        `json:"category"`
	Cuisines    []string        `json:"cuisines"`
	LocationUrl string          `json:"location_url"`
	Position    *DetailPosition `json:"position"`
	PriceRange  *PriceRange     `json:"price_range"`
	Rating      *Rating         `json:"rating"`
	Brand       *Brand          `json:"brand"`
	Delivery    *Delivery       `json:"delivery"`
	Url         string          `json:"url"`
}

type RestaurantMenu struct {
	DishTypeName string   `json:"dish_type_name"`
	Dishes       []*Dishes `json:"dishes"`
}

type Dishes struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	Price       *NewPrice `json:"price"`
	TotalOrder int64 `json:"total_order"`
}

type NewPrice struct {
	Text string `json:"text"`
}

type DishResp struct {
	Reply *Reply `json:"reply"`
}


type Reply struct {
	MenuInfos []*RestaurantMenu `json:"menu_infos"`
}


type DetailTotal struct {
	RestaurantMenu []*RestaurantMenu
	DeliveryDetail *DeliveryDetail
}
