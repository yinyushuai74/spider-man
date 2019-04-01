package dto

import (
	"strconv"
	"encoding/json"
)

type RegionResp struct {
	Cities []*City `json:"cities"`
}

type City struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
type ThMerchantResponse struct {
	THPage THPage `json:"page"`
}

type THPage struct {
	First                 int         `json:"first"`
	Last                  int         `json:"last"`
	TotalNumberOfPages    int         `json:"totalNumberOfPages"`
	TotalNumberOfEntities int         `json:"totalNumberOfEntities"`
	Entities              []*THEntity `json:"entities"`
}
type THEntity struct {
	Id                 int                `json:"id"`
	DisplayName        string             `json:"displayName"`
	NameOnly           Explain            `json:"NameOnly"`
	WorkingHoursStatus WorkingHoursStatus `json:"WorkingHoursStatus"`
	Branch             Explain            `json:"branch"`
	Contact            Contact            `json:"contact"`
	Categories         []*NameId          `json:"categories"`
	Rating             float64            `json:"rating"`
	RUrl               string             `json:"rUrl"`
	Lat                float64            `json:"lat"`
	Lng                float64            `json:"lng"`

}
type Explain struct {
	Primary string `json:"primary"`
	Thai    string `json:"thai"`
	English string `json:"english"`
}

type MerchantRow struct {
	Categories Category `json:"categories"`
	Basic      []string `json:"basic"`
}
type WorkingHoursStatus struct {
	Open    bool   `json:"open"`
	Message string `json:"message"`
}
type Contact struct {
	Address         Address `json:"address"`
	HomePage        string  `json:"homePage"`
	CallablePhoneNo string  `json:"callablePhoneno"`
}
type Address struct {
	Street      string `json:"street"`
	Hint        string `json:"hint"`
	SubDistrict NameId `json:"subDistrict"`
	District    NameId `json:"district"`
	City        NameId `json:"city"`
}
type NameId struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ThaiMenu struct {
	Categories []*Category `json:"categories"`
}
type Category struct {
	Name  string      `json:"name"`
	Items []*ThaiItem `json:"items"`
}
type ThaiItem struct {
	Name        string `json:"name"`
	Price       string `json:"price"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}

type MenuGroup struct {
	Name  string       `json:"name"`
	Items []*MatchItem `json:"items"`
}

type MatchItem struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Price       Price  `json:"price"`
}
type MatchResponse struct {
	MenuGroups []*MenuGroup `json:"menuGroups"`
}

type Price struct {
	Exact int    `json:"exact"`
	Text  string `json:"text"`
}

func (thaiMenu *ThaiMenu) ConvertMenu() string {
	data, _ := json.Marshal(thaiMenu)
	return string(data[:])
}

func (nameId *NameId) ConvertNameId() string {
	return nameId.Name + ":" + strconv.Itoa(nameId.Id)

}
