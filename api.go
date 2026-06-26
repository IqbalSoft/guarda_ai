package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/getlantern/systray"
)

func fetchZaiUsage(m5H, m5HR, mWk, mWkR, mMo, mMoR *systray.MenuItem, cfg AIConfig) {
	if !cfg.Enabled {
		return
	}
	if cfg.APIKey == "" {
		m5H.SetTitle("[Missing API Key]")
		return
	}

	req, _ := http.NewRequest("GET", "https://api.z.ai/api/monitor/usage/quota/limit", nil)
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	type ZaiResp struct {
		Data struct {
			Limits []struct {
				Type          string  `json:"type"`
				Unit          int     `json:"unit"`
				Percentage    float64 `json:"percentage"`
				NextResetTime int64   `json:"nextResetTime"`
			} `json:"limits"`
		} `json:"data"`
	}

	var apiData ZaiResp
	json.Unmarshal(body, &apiData)

	pct5Hour, pctWeekly, pctMonthly := 0, 0, 0
	reset5Hour, resetWeekly, resetMonthly := "N/A", "N/A", "N/A"

	jakartaLocation, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		jakartaLocation = time.UTC
	}

	for _, limit := range apiData.Data.Limits {
		resetTimeStr := "N/A"
		if limit.NextResetTime > 0 {
			resetTime := time.UnixMilli(limit.NextResetTime).In(jakartaLocation)
			resetTimeStr = resetTime.Format("2006-01-02 15:04") + " WIB"
		}

		switch {
		case limit.Type == "TOKENS_LIMIT" && limit.Unit == 3:
			pct5Hour = int(limit.Percentage)
			reset5Hour = resetTimeStr
		case limit.Type == "TOKENS_LIMIT" && limit.Unit == 6:
			pctWeekly = int(limit.Percentage)
			resetWeekly = resetTimeStr
		case limit.Type == "TIME_LIMIT" && limit.Unit == 5:
			pctMonthly = int(limit.Percentage)
			resetMonthly = resetTimeStr
		}
	}

	m5H.SetTitle(fmt.Sprintf("[%s] %d%% Used", generateBar(pct5Hour), pct5Hour))
	m5HR.SetTitle(fmt.Sprintf("Reset Time: %s", reset5Hour))
	mWk.SetTitle(fmt.Sprintf("[%s] %d%% Used", generateBar(pctWeekly), pctWeekly))
	mWkR.SetTitle(fmt.Sprintf("Reset Time: %s", resetWeekly))
	mMo.SetTitle(fmt.Sprintf("[%s] %d%% Used", generateBar(pctMonthly), pctMonthly))
	mMoR.SetTitle(fmt.Sprintf("Reset Time: %s", resetMonthly))
}

func fetchGenericAI(name string, mBar *systray.MenuItem, cfg AIConfig) {
	if !cfg.Enabled {
		return
	}
	if cfg.APIKey == "" {
		mBar.SetTitle("✨ API key is waiting — set it in Settings")
		return
	}

	// Placeholder logic until you add real API endpoints for these tools
	mockPercentage := 25
	mBar.SetTitle(fmt.Sprintf("[%s] %d%% Used", generateBar(mockPercentage), mockPercentage))
}
