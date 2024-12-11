package functions

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type FinalGeo struct {
	Latitude  string `json:"lat"`
	Longitude string `json:"lon"`
	MapURL    string
}

func MapsApiResp(city, country string) (FinalGeo, error) {
		baseURL := "https://nominatim.openstreetmap.org/search"
	
		//define
		params := url.Values{}
		//add place to params q is the key
		params.Add("q", city+","+ country)
		//add desired format of response
		params.Add("format", "json")
	
		// make the URL required and encode commas etc betweenn params
		fullURL := baseURL+"?"+ params.Encode()
	
		//request for data
		res, err := http.Get(fullURL)
		if err != nil {
			return FinalGeo{}, fmt.Errorf("error sending request: %v", err)
		}
		defer res.Body.Close()
	
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return FinalGeo{}, fmt.Errorf("error reading response: %v", err)
		}
	
		fmt.Println("Raw API Response:", string(body))
	
		var finalgeo []FinalGeo
		err = json.Unmarshal(body, &finalgeo)
		if err != nil {
			return FinalGeo{}, fmt.Errorf("error unmarshaling response: %v", err)
		}

		finalgeo[0].MapURL = "https://www.google.com/maps?q=" +finalgeo[0].Latitude+","+ finalgeo[0].Longitude


	fmt.Println("finalgeo",finalgeo)
		if len(finalgeo) == 0 {
			return FinalGeo{}, fmt.Errorf("no results found for city: %s", city)
		}
		fmt.Println("finalgeo[0]:",finalgeo[0])
		return finalgeo[0], nil
	}
	
	