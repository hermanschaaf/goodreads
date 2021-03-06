// Package goodreads provides a REST client for the public goodreads.com API.
//
// https://www.goodreads.com/api
package goodreads

import (
	"fmt"
	"net/url"
	"strconv"
)

// Client wraps the public Goodreads API.
type Client struct {
	ApiKey  string
	Decoder Decoder
}

// NewClient initializes a Client with default parameters.
func NewClient(key string) *Client {
	return &Client{
		ApiKey:  key,
		Decoder: DefaultDecoder,
	}
}

// AuthorBooks returns a list of books by a particular author.
// https://www.goodreads.com/api/index#author.books
func (c *Client) AuthorBooks(authorID string, page int) (*Author, error) {
	v := c.defaultValues()
	if page > 0 {
		v.Set("page", strconv.Itoa(page))
	}

	var r struct {
		Author Author `xml:"author"`
	}
	err := c.Decoder.Decode(fmt.Sprintf("author/list/%s", authorID), v, &r)
	if err != nil {
		return nil, err
	}
	return &r.Author, nil
}

// AuthorShow returns the full details of an author.
// https://www.goodreads.com/api/index#author.show
func (c *Client) AuthorShow(authorID string) (*Author, error) {
	var r struct {
		Author Author `xml:"author"`
	}
	err := c.Decoder.Decode(fmt.Sprintf("author/show/%s", authorID), c.defaultValues(), &r)
	if err != nil {
		return nil, err
	}
	return &r.Author, nil
}

// ReviewList returns the books on a members shelf.
// https://www.goodreads.com/api/index#reviews.list
func (c *Client) ReviewList(userID, shelf, sort, search, order string, page, perPage int) ([]Review, error) {
	v := c.defaultValues()
	v.Set("id", userID)
	v.Set("v", "2")
	if shelf != "" {
		v.Set("shelf", shelf)
	}
	if sort != "" {
		v.Set("sort", sort)
	}
	if search != "" {
		v.Set("search", search)
	}
	if order != "" {
		v.Set("order", order)
	}
	if page > 0 {
		v.Set("page", strconv.Itoa(page))
	}
	if perPage > 0 {
		v.Set("per_page", strconv.Itoa(perPage))
	}

	var r struct {
		Reviews []Review `xml:"reviews>review"`
	}
	err := c.Decoder.Decode("review/list", v, &r)
	if err != nil {
		return nil, err
	}
	return r.Reviews, nil
}

// ShelvesList returns the list of shelves belonging to a user.
// https://www.goodreads.com/api/index#shelves.list
func (c *Client) ShelvesList(userID string) ([]UserShelf, error) {
	v := c.defaultValues()
	v.Set("user_id", userID)
	var r struct {
		Shelves []UserShelf `xml:"shelves>user_shelf"`
	}
	err := c.Decoder.Decode("shelves/list", v, &r)
	if err != nil {
		return nil, err
	}
	return r.Shelves, nil
}

// UserShow returns the public information about a given Goodreads user.
// https://www.goodreads.com/api/index#user.show
func (c *Client) UserShow(id string) (*User, error) {
	var r struct {
		User User `xml:"user"`
	}
	err := c.Decoder.Decode(fmt.Sprintf("user/show/%s.xml", id), c.defaultValues(), &r)
	if err != nil {
		return nil, err
	}
	return &r.User, nil
}

func (c *Client) defaultValues() url.Values {
	v := url.Values{}
	v.Set("key", c.ApiKey)
	return v
}
