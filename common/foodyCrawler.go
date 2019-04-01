package common

import (
	"net/http"
	"io/ioutil"
	"grab-crawler/dto"
	"encoding/json"
	"fmt"
)

const host = "https://www.foody.vn"
const locationUrl = "/__get/Common/GetPopupLocation"

func QueryLocationList() []*dto.FProvince {
	url := host + locationUrl
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("X-Requested-With", "XMLHttpRequest")
	response, _ := client.Do(request)
	body, err := ioutil.ReadAll(response.Body)
	fmt.Println(body)
	if err != nil {
		return nil
	}
	locationResponse := &dto.LocationResponse{}
	json.Unmarshal(body,locationResponse)
	return locationResponse.AllLocations
}
