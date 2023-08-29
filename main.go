package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type WeatherResponse struct {
	City    string `json:"city"`
	Weather string `json:"weather"`
}

var WeatherData []WeatherResponse

func main() {
	r := gin.Default()

	r.GET("/city/:name", func(c *gin.Context) {
		cityName := c.Param("name")

		weather, err := getWeather(cityName)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to fetch weather")
			return
		}

		response := fmt.Sprintf("%s %s", cityName, weather)

		saveOrUpdateJSON(cityName, response)

		c.String(http.StatusOK, response)

	})

	port := ":8080"
	fmt.Printf("Server listening on port %s\n", port)
	log.Fatal(r.Run(port))
}

// func getPincode(cityName string) (string, error) {
// }

func getWeather(cityName string) (string, error) {
	url := "https://weatherapi-com.p.rapidapi.com/current.json?q=53.1%2C-0.13"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", "1d04ba4a72mshb73180c7156b027p146ea5jsnefd4d25a8fd3")
	req.Header.Add("X-RapidAPI-Host", "weatherapi-com.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	return fmt.Sprint(string(body)), nil

}

func saveOrUpdateJSON(cityName, weather string) {
	response := WeatherResponse{
		City:    cityName,
		Weather: weather,
	}

	fileName := "city.json"
	filePath := filepath.Join(".", "responses", fileName)

	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Println("Failed to create directory:", err)
			return
		}
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Failed to open file:", err)
		return
	}
	defer file.Close()

	// Encode the response struct into JSON format
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(response); err != nil {
		fmt.Println("Failed to encode JSON:", err)
	}
}
