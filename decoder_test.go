package goodreads

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpDecoder_Decode(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/foo/bar?p1=v1&p2=v2", r.URL.String())
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><response><id>SampleID</id></response>`))
	}))
	defer s.Close()

	v := url.Values{}
	v.Set("p1", "v1")
	v.Set("p2", "v2")
	var res struct {
		ID string `xml:"id"`
	}
	h := HttpDecoder{Client: http.DefaultClient, ApiRoot: s.URL, Verbose: true}
	err := h.Decode("foo/bar", v, &res)
	assert.Nil(t, err)
	assert.Equal(t, "SampleID", res.ID)
}
