package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type InstagramItem struct {
	ShortcodeMedia struct {
		ID                 string `json:"id"`
		Shortcode          string `json:"shortcode"`
		Typename           string `json:"__typename"`
		DisplayURL         string `json:"display_url"`
		EdgeMediaToCaption struct {
			Edges []struct {
				Node struct {
					Text string `json:"text"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"edge_media_to_caption"`

		EdgeSidecarToChildren struct {
			Edges []struct {
				Node struct {
					Typename   string `json:"__typename"`
					ID         string `json:"id"`
					Shortcode  string `json:"shortcode"`
					DisplayURL string `json:"display_url"`
					VideoUrl   string `json:"video_url"`
					IsVideo    bool   `json:"is_video"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"edge_sidecar_to_children"`

		HasAudio       bool        `json:"has_audio"`
		VideoURL       string      `json:"video_url"`
		VideoViewCount int         `json:"video_view_count"`
		VideoPlayCount interface{} `json:"video_play_count"`

		IsAd              bool        `json:"is_ad"`
		IsAffiliate       bool        `json:"is_affiliate"`
		IsPaidPartnership bool        `json:"is_paid_partnership"`
		IsVideo           bool        `json:"is_video"`
		Location          interface{} `json:"location"`

		Owner struct {
			ID            string `json:"id"`
			IsVerified    bool   `json:"is_verified"`
			ProfilePicURL string `json:"profile_pic_url"`
			Username      string `json:"username"`
			FullName      string `json:"full_name"`
		} `json:"owner"`
	} `json:"shortcode_media"`
}

type MagicItem struct {
	Url               string              `json:"url"`
	Type              string              `json:"type"`
	Caption           string              `json:"caption"`
	DisplayUrl        string              `json:"display_url"`
	VideoUrl          string              `json:"video_url"`
	ProfilePicUrl     string              `json:"profile_pic_url"`
	Username          string              `json:"username"`
	Shortcode         string              `json:"shortcode"`
	MultipleMediaList []MultipleMediaItem `json:"multiple_media_list"`
}

type MultipleMediaItem struct {
	Type       string `json:"type"`
	DisplayUrl string `json:"display_url"`
	VideoUrl   string `json:"video_url"`
}

func NewMagicItem() *MagicItem {
	return &MagicItem{}
}

func (mi *MagicItem) MarshalBinary() ([]byte, error) {
	return json.Marshal(mi)
}

func (mi *MagicItem) makeRequest() ([]byte, error) {

	//fmt.Println("url", generateUrl(mi.Shortcode), mi.Shortcode)
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Get(generateUrl(mi.Shortcode))
	if os.IsTimeout(err) {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

func (mi *MagicItem) convertRequest(body []byte) error {
	parseData := parseAdditionalData(body)
	if parseData != "" {
		var instaItem InstagramItem

		err := json.Unmarshal([]byte(parseData), &instaItem)
		if err != nil {
			return err
		}

		mi.Type = instaItem.ShortcodeMedia.Typename
		if len(instaItem.ShortcodeMedia.EdgeMediaToCaption.Edges) > 0 {
			mi.Caption = instaItem.ShortcodeMedia.EdgeMediaToCaption.Edges[0].Node.Text
		}
		mi.DisplayUrl = instaItem.ShortcodeMedia.DisplayURL
		mi.VideoUrl = instaItem.ShortcodeMedia.VideoURL
		mi.ProfilePicUrl = instaItem.ShortcodeMedia.Owner.ProfilePicURL
		mi.Username = instaItem.ShortcodeMedia.Owner.Username
		mi.Shortcode = instaItem.ShortcodeMedia.Shortcode
		mi.MultipleMediaList = []MultipleMediaItem{}

		if len(instaItem.ShortcodeMedia.EdgeSidecarToChildren.Edges) > 0 {
			for _, item := range instaItem.ShortcodeMedia.EdgeSidecarToChildren.Edges {
				mi.MultipleMediaList = append(mi.MultipleMediaList, MultipleMediaItem{
					Type:       item.Node.Typename,
					DisplayUrl: item.Node.DisplayURL,
					VideoUrl:   item.Node.VideoUrl,
				})
			}
		}
	} else {
		// Look at the html content
		// Type
		re := regexp.MustCompile(`(?s)data-media-type="(.*?)"`)
		mediaTypeData := re.FindAllStringSubmatch(string(body), -1)

		//fmt.Println("MediaType", mediaTypeData[0][1])
		if len(mediaTypeData) > 0 && len(mediaTypeData[0]) > 1 {
			mi.Type = mediaTypeData[0][1]
		}

		// Caption
		re = regexp.MustCompile(`(?s)class="Caption"(.*?)class="CaptionUsername"(.*?)<\/a>(.*?)<div`)
		captionData := re.FindAllStringSubmatch(string(body), -1)

		//fmt.Println("CaptionData", captionData[0])
		if len(captionData) > 0 && len(captionData[0]) > 2 {
			re = regexp.MustCompile(`<[^>]*>`)
			mi.Caption = strings.TrimSpace(re.ReplaceAllString(captionData[0][3], ""))
		}

		// Main Media
		re = regexp.MustCompile(`(?s)class="Content(.*?)src="(.*?)"`)
		mainMediaData := re.FindAllStringSubmatch(string(body), -1)

		//fmt.Println("mainMediaData", strings.ReplaceAll(mainMediaData[0][2], "amp;", ""))
		if len(mainMediaData) > 0 && len(mainMediaData[0]) > 1 {
			// Do not find video url (GraphVideo). Only image url SORRY!!!
			if mi.Type == "GraphImage" {
				mi.Type = "GraphImage"
				// VideoUrl
				mi.VideoUrl = ""
			}

			// DisplayUrl
			mi.DisplayUrl = strings.ReplaceAll(mainMediaData[0][2], "amp;", "")
		}

		// UserData
		re = regexp.MustCompile(`(?s)class="Header(.*?)<a class="Avatar(.*?)" href="(.*?)"(.*?)img src="(.*?)" alt="(.*?)"(.*?)<div`)
		userData := re.FindAllStringSubmatch(string(body), -1)

		//fmt.Println("userData", userData[0][4], userData[0][5], userData[0][6])
		if len(userData) > 0 && len(userData[0]) > 6 {
			// ProfilePicUrl
			mi.ProfilePicUrl = strings.ReplaceAll(userData[0][5], "amp;", "")

			// Username
			mi.Username = userData[0][6]
		}
	}

	return nil
}

func (mi *MagicItem) validateUrl(url string) bool {
	v := regexp.MustCompile(`(https?:\/\/(?:www\.)?instagram\.com\/(p|tv|reel)\/([^\/?#&]+)).*`)
	parseData := v.FindAllStringSubmatch(url, -1)

	fmt.Println("validateUrl", parseData)
	if len(parseData) > 0 {
		//parseData[0]
		//parseData[0][1] // url
		//parseData[0][2] // type
		//parseData[0][3] // shortcode

		if len(parseData[0]) > 2 {
			mi.Url = parseData[0][1]
			mi.Shortcode = parseData[0][3]

			return true
		}
	}

	return false
}

func generateUrl(shortcode string) string {
	return "https://www.instagram.com/p/" + shortcode + "/embed/captioned/"
}

func parseAdditionalData(body []byte) string {
	re := regexp.MustCompile(`(?s)(window\.__additionalDataLoaded\('extra',)(.*?)\)(;<\/script>)`)
	additionalData := re.FindAllStringSubmatch(string(body), -1)

	//fmt.Println("parseAdditionalData", additionalData)
	if len(additionalData) > 0 && len(additionalData[0]) > 2 && additionalData[0][2] != "null" {
		return additionalData[0][2]
	}

	return ""
}
