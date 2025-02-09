package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type House struct {
	Title  string `json:"title"`
	Rooms  int    `json:"rooms"`
	Suburb string `json:"suburb"`
	Price  string `json:"price"`
	Link   string `json:"link"`
}

func fetchHouses(suburb string, rooms int, price int) ([]House, error) {
	url := fmt.Sprintf("http://localhost:5000/api/houses?suburb=%s&rooms=%d&price=%d", suburb, rooms, price)
	log.Printf("Fetching houses from URL: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching houses: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, err
	}

	log.Printf("Response body: %s", string(body))

	var houses []House
	err = json.Unmarshal(body, &houses)
	if err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		return nil, err
	}

	return houses, nil
}

func main() {
	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")
	r.StaticFile("/", "./index.html")

	r.GET("/houses", func(c *gin.Context) {
		suburb := c.Query("suburb")
		if suburb == "" {
			log.Printf("Suburb parameter is required")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Suburb parameter is required"})
			return
		}

		roomsStr := c.Query("rooms")
		priceStr := c.Query("price")

		var rooms int
		var price int
		var err error

		if roomsStr != "" {
			rooms, err = strconv.Atoi(roomsStr)
			if err != nil {
				log.Printf("Error parsing rooms: %v", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rooms parameter"})
				return
			}
		}

		if priceStr != "" {
			price, err = strconv.Atoi(priceStr)
			if err != nil {
				log.Printf("Error parsing price: %v", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price parameter"})
				return
			}
		}

		log.Printf("Received request: suburb=%s, rooms=%d, price=%d", suburb, rooms, price)

		houses, err := fetchHouses(suburb, rooms, price)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Printf("Found %d houses matching criteria", len(houses))
		c.JSON(http.StatusOK, houses)
	})

	r.Run(":8080")
}
