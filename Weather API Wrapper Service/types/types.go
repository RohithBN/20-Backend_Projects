package types



type WeatherAPIRequest struct{
	Location string `json:"location"`
	Date1 string `json:"date1"`
	Date2 string `json:"date2"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}


type WeatherAPIResponse struct{
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	ResolvedAddress string `json:"resolvedAddress"`
	Timezone string `json:"timezone"`
	Address string `json:"address"`
	Days []any `json:"days"`
}