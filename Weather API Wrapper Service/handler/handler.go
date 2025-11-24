package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RohithBN/redis"
	"github.com/RohithBN/types"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var apiKey string
var apiEndpoint = "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/"

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Get API key
	apiKey = os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY environment variable is required")
	}
}

func FetchWeatherData(c *gin.Context) {

	var weatherRequest types.WeatherAPIRequest

	if err := c.ShouldBindJSON(&weatherRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if weatherRequest.Location == "" {
		c.JSON(400, gin.H{"error": "Location is required"})
		return
	}

	cachedLocation, err := redis.RDB.Get(redis.Ctx,weatherRequest.Location).Result()
	if err == nil {
		var cachedResponse types.WeatherAPIResponse
		if err:= json.Unmarshal([]byte(cachedLocation), &cachedResponse); err == nil {
			c.JSON(200, cachedResponse)
			return
		}
		fmt.Println("Cache Hit for location: ", weatherRequest.Location)
	}

	fmt.Println("Cache miss for location: ", weatherRequest.Location)

	weatherData, err := getWeatherDataFromAPI(weatherRequest.Location, weatherRequest.Date1, weatherRequest.Date2)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, weatherData)

}

func getWeatherDataFromAPI(location, startDate, endDate string) (*types.WeatherAPIResponse, error) {
	var weatherResponse types.WeatherAPIResponse

	// Build request URL
	requestURL := apiEndpoint + location
	if startDate != "" && endDate != "" {
		requestURL += "/" + startDate + "/" + endDate
	}
	requestURL += "?unitGroup=metric&key=" + apiKey + "&contentType=json"

	// Make HTTP GET request
	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}
	// Decode JSON response
	if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		return nil, err
	}
	// Cache the response in Redis
	responseBytes,err:= json.Marshal(weatherResponse)
	if err==nil{
		err= redis.RDB.Set(redis.Ctx,location,string(responseBytes),0).Err()
		if err!=nil{
			log.Printf("Failed to cache response for location %s: %v", location, err)
		}
	}else{
		log.Printf("Failed to marshal response for caching for location %s: %v", location, err)
	}
	return &weatherResponse, nil
}
