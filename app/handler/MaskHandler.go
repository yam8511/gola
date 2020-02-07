package handler

import (
	"bytes"
	"encoding/csv"
	errorCode "gola/app/common/error_code"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// SearchMask 搜尋口罩
func SearchMask(c *gin.Context) {
	req := c.DefaultQuery("q", "")
	q := strings.TrimSpace(req)
	if q == "" {
		c.JSON(http.StatusOK, []struct{}{})
		return
	}

	words := strings.Split(q, " ")
	tmpWords := map[string]int{}
	for _, word := range words {
		word = strings.TrimSpace(word)
		tmpWords[word] = 0
	}
	words = []string{}
	for word := range tmpWords {
		word = strings.TrimSpace(word)
		if word == "" {
			continue
		}
		words = append(words, word)
	}

	if len(words) == 0 {
		c.JSON(http.StatusOK, []struct{}{})
		return
	}

	buffer, err := fetchMaskAPI()
	if err != nil {
		c.JSON(http.StatusOK, errorCode.GetAPIError("mask_api_error", err))
		return
	}

	r := csv.NewReader(buffer)
	_, err = r.Read()
	if err != nil {
		c.JSON(http.StatusOK, errorCode.GetAPIError("read_mask_csv_error", err))
		return
	}

	type maskInfo struct {
		Name       string `json:"name"`
		NameURL    string `json:"name_url"`
		Address    string `json:"address"`
		AddressURL string `json:"address_url"`
		Number     string `json:"number"`
		Child      string `json:"child"`
		Phone      string `json:"phone"`
		UpdatedAt  string `json:"updated_at"`
	}

	data := []maskInfo{}
	for {
		record, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			c.JSON(http.StatusOK, errorCode.GetAPIError("read_mask_csv_error", err))
			return
		}

		address := strings.TrimSpace(record[2])
		address = strings.ReplaceAll(address, "臺", "台")
		address = strings.ReplaceAll(address, "－", "-")
		address = strings.ReplaceAll(address, " ", "")
		for i := 0; i <= 9; i++ {
			switch i {
			case 0:
				address = strings.ReplaceAll(address, "０", "0")
			case 1:
				address = strings.ReplaceAll(address, "１", "1")
			case 2:
				address = strings.ReplaceAll(address, "２", "2")
			case 3:
				address = strings.ReplaceAll(address, "３", "3")
			case 4:
				address = strings.ReplaceAll(address, "４", "4")
			case 5:
				address = strings.ReplaceAll(address, "５", "5")
			case 6:
				address = strings.ReplaceAll(address, "６", "6")
			case 7:
				address = strings.ReplaceAll(address, "７", "7")
			case 8:
				address = strings.ReplaceAll(address, "８", "8")
			case 9:
				address = strings.ReplaceAll(address, "９", "9")
			}
		}

		num := strings.TrimSpace(record[4])
		child := strings.TrimSpace(record[5])
		name := strings.TrimSpace(record[1])
		phone := strings.TrimSpace(record[3])
		updatedAt := strings.TrimSpace(record[6])

		stage1 := words[0]
		if (!strings.Contains(address, stage1) && !strings.Contains(name, stage1)) ||
			(num == "0" && child == "0") {
			continue
		}

		if len(words) > 1 {
			stage2 := words[1:]
			for _, word := range stage2 {
				if (strings.Contains(address, word) || strings.Contains(name, word)) &&
					(num != "0" || child != "0") {
					data = append(data, maskInfo{
						Name:       name,
						NameURL:    "https://www.google.com.tw/maps/search/" + name,
						Address:    address,
						AddressURL: "https://www.google.com.tw/maps/search/" + address,
						Child:      child,
						Number:     num,
						Phone:      phone,
						UpdatedAt:  updatedAt,
					})
					break
				}
			}
		} else {
			data = append(data, maskInfo{
				Name:       name,
				NameURL:    "https://www.google.com.tw/maps/search/" + name,
				Address:    address,
				AddressURL: "https://www.google.com.tw/maps/search/" + address,
				Child:      child,
				Number:     num,
				Phone:      phone,
				UpdatedAt:  updatedAt,
			})
		}
	}

	c.JSON(http.StatusOK, data)
}

func fetchMaskAPI() (*bytes.Buffer, error) {
	url := "https://data.nhi.gov.tw/resource/mask/maskdata.csv"
	method := "GET"
	client := &http.Client{Timeout: time.Second * 5}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(body)
	return buf, nil
}
