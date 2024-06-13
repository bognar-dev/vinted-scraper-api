package vinted_scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Order string

const (
	NEWEST_FIRST      Order = "newest_first"
	RELEVANCE         Order = "relevance"
	PRICE_HIGH_TO_LOW Order = "price_high_to_low"
	PRICE_LOW_TO_HIGH Order = "price_low_to_high"
)

func ToOrder(order string) Order {
	switch order {
	case "newest_first":
		return NEWEST_FIRST
	case "relevance":
		return RELEVANCE
	case "price_high_to_low":
		return PRICE_HIGH_TO_LOW
	case "price_low_to_high":
		return PRICE_LOW_TO_HIGH
	default:
		return NEWEST_FIRST
	}
}

func FetchCookie(domain string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.vinted.%s", domain), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Cookie", os.Getenv("VINTED_COOKIE"))
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	sessionCookie := resp.Header.Get("Set-Cookie")
	cookies := strings.Split(sessionCookie, ";")
	for _, cookie := range cookies {
		fmt.Println(cookie)
		if strings.Contains(cookie, "_vinted_fr_session") {
			return cookie, nil
		}
	}
	return "", fmt.Errorf("cookie not found")
}

func Search(query string, order Order, currency string) ([]byte, error) {

	cookie, err := FetchCookie("co.uk")
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.vinted.co.uk/api/v2/catalog/items?search_text=%s&currency=%s&order=%s", query, currency, order), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Cookie", fmt.Sprintf("%s", cookie))
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response VintedApi_Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return body, nil
}
