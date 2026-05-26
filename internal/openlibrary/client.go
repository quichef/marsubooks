package openlibrary

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type SearchResult struct {
	Title     string `json:"title"`
	Author    string `json:"author"`
	ISBN      string `json:"isbn"`
	CoverURL  string `json:"cover_url"`
}

type olSearchResponse struct {
	Docs []struct {
		Title       string   `json:"title"`
		AuthorNames []string `json:"author_name"`
		ISBN        []string `json:"isbn"`
		CoverI      int      `json:"cover_i"`
	} `json:"docs"`
}

func Search(query string) ([]SearchResult, error) {
	u := "https://openlibrary.org/search.json?limit=8&q=" + url.QueryEscape(query)
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openlibrary: status %d", resp.StatusCode)
	}

	var data olSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var results []SearchResult
	for _, d := range data.Docs {
		r := SearchResult{Title: d.Title}
		if len(d.AuthorNames) > 0 {
			r.Author = d.AuthorNames[0]
		}
		if len(d.ISBN) > 0 {
			r.ISBN = d.ISBN[0]
		}
		if d.CoverI > 0 {
			r.CoverURL = fmt.Sprintf("https://covers.openlibrary.org/b/id/%d-M.jpg", d.CoverI)
		}
		results = append(results, r)
	}
	return results, nil
}
