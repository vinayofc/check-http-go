package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HttpCheck(website string) string {
	url := fmt.Sprintf("https://check-host.net/check-http?host=%s&max_nodes=10", website)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return err.Error()
	}
	req.Header.Add("Accept", "application/json") // Fixed Accept header

	res, err := client.Do(req)
	if err != nil {
		return err.Error()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err.Error()
	}
	var data map[string]interface{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return err.Error()
	}
	permanentLink, ok := data["request_id"].(string)
	if !ok {
		return "CHK_ERROR"
	}
	return permanentLink
}

func ResultHTTP(req_id string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://check-host.net/check-result/%s", req_id)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", " application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	// Remove null values
	for key, value := range result {
		if value == nil {
			delete(result, key)
		}
	}

	return result, nil
}

func main() {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"author": "Vinay Chaudhary",
			"result": "GO Lang HTTP Request API",
		})
	})

	r.GET("/check", func(c *gin.Context) {
		website := c.Query("url")
		if website == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "URL parameter is required"})
			return
		}

		res := HttpCheck(website)
		var msg string
		var finalRes map[string]interface{}
		if res == "CHK_ERROR" {
			msg = "Error"
		} else {
			var err error
			finalRes, err = ResultHTTP(res)
			if err != nil {
				msg = "Error fetching result"
				finalRes = nil
			} else {
				msg = "Success"
			}
		}
		c.JSON(200, gin.H{
			"author":   "Vinay Chaudhary",
			"response": msg,
			"result":   finalRes,
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (default)
}
