package main

// album represent data about a record album
type album struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Price  string `json:"price"`
}
