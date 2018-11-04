package fiftheditionsrd

import (
	"net/url"
)

type NamedAPIResource struct {
	Name string  `json:"name,omitempty"`
	Url  url.URL `json:"url,omitempty"`
}

type ClassAPIResource struct {
	Class string
	Url   url.URL
}

type Choice struct {
	Choose int
	Type   string
	from   []NamedAPIResource
}

type Cost struct {
	Quantity int
	Unit     string
}

type NamedAPIResourceList struct {
	Count   int
	Results []NamedAPIResource
}

type ClassAPIResourceList struct {
	Count   int
	Results []ClassAPIResource
}

type AbilityScore struct {
	ID          string             `json:"_id,omitempty"`
	Index       int                `json:"index,omitempty"`
	Name        string             `json:"name,omitempty"`
	FullName    string             `json:"full_name,omitempty"`
	Description string             `json:"desc,omitempty"`
	Skills      []NamedAPIResource `json:"skills,omitempty"`
	URL         url.URL            `json:"url,omitempty"`
}
