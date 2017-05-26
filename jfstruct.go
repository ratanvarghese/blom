package main

type Author struct {
	Name string `json:"name"`
}

type Attachment struct {
	Url       string `json:"url"`
	Mime_type string `json:"mime_type"`
}

type ArticleItem struct {
	Id             string       `json:"id"`
	Url            string       `json:"url"`
	Title          string       `json:"title"`
	Content_html   string       `json:"content_html"`
	Date_published string       `json:"date_published"`
	Date_modified  string       `json:"date_modified"`
	Tags           []string     `json:"tags"`
	Attachments    []Attachment `json:"attachments"`
}

type JSONFeed struct {
	Version       string        `json:"version"`
	Title         string        `json:"title"`
	Home_page_url string        `json:"home_page_url"`
	Feed_url      string        `json:"feed_url"`
	Next_url      string        `json:"next_url"`
	Icon          string        `json:"icon"`
	Favicon       string        `json:"favicon"`
	Author        Author        `json:"author"`
	Items         []ArticleItem `json:"items"`
}
