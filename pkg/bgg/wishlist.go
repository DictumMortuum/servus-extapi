package bgg

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Status struct {
	Own              string `xml:"own,attr" json:"own"`
	PrevOwned        string `xml:"prevowned,attr" json:"prevowned"`
	ForTrade         string `xml:"fortrade,attr" json:"fortrade"`
	Want             string `xml:"want,attr" json:"want"`
	WantToBuy        string `xml:"wanttobuy,attr" json:"wanttobuy"`
	Wishlist         string `xml:"wishlist,attr" json:"wishlist"`
	WishlistPriority string `xml:"wishlistpriority,attr" json:"wishlistpriority"`
	Preordered       string `xml:"preordered,attr" json:"preordered"`
}

type WishlistItem struct {
	// XMLName  xml.Name `xml:"item" json:"-"`
	ObjectId int64  `xml:"objectid,attr" json:"id"`
	Name     string `xml:"name" json:"name"`
	Status   Status `xml:"status"`
}

type WishlistRs struct {
	// XMLName xml.Name       `xml:"items" json:"-"`
	Items []WishlistItem `xml:"item" json:"items"`
}

func Wishlist(name string) (*WishlistRs, error) {
	link := fmt.Sprintf("https://api.geekdo.com/xmlapi2/collection?username=%s", url.QueryEscape(name))
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, err
	}

	conn := &http.Client{}

	resp, err := conn.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	tmp := WishlistRs{}

	err = xml.Unmarshal(body, &tmp)
	if err != nil {
		return nil, err
	}

	return &tmp, nil
}
