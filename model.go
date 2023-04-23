package main

type Book struct {
	Id          string     `json:"id"`
	Title       string     `json:"title"`
	ISBN        string     `json:"isbn"`
	Language    string     `json:"language"`
	ImageUrl    string     `json:"image_url"`
	Description string     `json:"description"`
	Link        string     `json:"link"`
	PageCount   int        `json:"page_count"`
	Categories  [][]string `json:"categories"`
	Authors     []Author   `json:"authors"`
}

type Author struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Biography string `json:"biography"`
	ImageUrl  string `json:"image_url"`
	Link      string `json:"link"`
}
