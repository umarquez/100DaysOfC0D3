package main

type colorTotal struct {
	RGB   string `json:"rgb"`
	Total int    `json:"total"`
}

type processResults struct {
	Colors    []colorTotal `json:"colors"`
	Height    int          `json:"height"`
	Width     int          `json:"width"`
	TotalTime float64      `json:"total_time"`
}
