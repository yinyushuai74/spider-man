package dto

type MerchantResp struct {
	Success bool        `json:"success"`
	Data    []*Merchant `json:"data"`
}

type Merchant struct {
	Id              int        `json:"Id"`
	Title           string     `json:"title"`
	AttributeString string     `json:"AttributeString"`
	DetailUrl       string     `json:"DetailUrl"`
	Restaurant      Restaurant `json:"Restaurant"`
}

type Restaurant struct {
	ResName            string `json:"ResName"`
	FullAddress        string `json:"FullAddress"`
	RestaurantLocation string `json:"RestaurantLocation"`
	Phone              string `json:"Phone"`
	DistrictId         int    `json:"DistrictId"`
}

type MerchantDetail struct {
	MerchantName     string          `json:"merchant_name"`
	MerchantPhone    string          `json:"merchant_phone"`
	MerchantAddress  string          `json:"merchant_address"`
	District         int             `json:"district"`
	City             string          `json:"city"`
	MerchantCategory string          `json:"merchant_category"`
	OperatingHours   string          `json:"operating_hours"`
	Rating           float64         `json:"rating"`
	MenuCategories   []*MenuCategory `json:"menu_categories"`
}
type MenuCategory struct {
	MenuTitle string          `json:"menu_title"`
	Item      []*MerchantItem `json:"item"`
}
type MerchantItem struct {
	Name           string `json:"name"`
	Price          string `json:"price"`
	FoodOrderTimes string `json:"food_order_times"`
	Describe       string `json:"describe"`
}
type TotalPageResp struct {
	Result Result `json:"result"`
}
type Result struct {
	MerchantTotalList []*MerchantTotal `json:"ListResult"`
	TotalCount        int              `json:"totalCount"`
}
type MerchantTotal struct {
	Id          int     `json:"RestaurantId"`
	DetailUrl   string  `json:"DetailUrl"`
	ResName     string  `json:"RestaurantName"`
	FullAddress string  `json:"RestaurantAddress"`
	TimeRange   string  `json:"TimeRange"`
	Category    string  `json:"Category"`
	PriceRange  string  `json:"PriceRange"`
	DistrictId  int     `json:"DistricId"`
	ResLat      float64 `json:"ResLat"`
	ResLng      float64 `json:"ResLng"`
}
type MerchantDetailTotal struct {
	Id               int             `json:"RestaurantId"`
	MerchantName     string          `json:"merchant_name"`
	MerchantPhone    string          `json:"merchant_phone"`
	MerchantAddress  string          `json:"merchant_address"`
	District         int             `json:"district"`
	PriceRange       string          `json:"PriceRange"`
	City             string          `json:"city"`
	MerchantCategory string          `json:"merchant_category"`
	OperatingHours   string          `json:"operating_hours"`
	Rating           float64         `json:"rating"`
	MenuCategories   []*MenuCategory `json:"menu_categories"`
	Lat              float64         `json:"lat"`
	Lng              float64         `json:"lng"`
	Url              string          `json:"url"`
}
type NOWDistrict struct {
	Result []*DistrictResult `json:"result"`
}
type DistrictResult struct {
	AsciiName string `json:"AsciiName"`
	Id        int    `json:"Id"`
}
