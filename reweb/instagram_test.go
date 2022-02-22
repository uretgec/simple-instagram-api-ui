package main

import (
	"strings"
	"testing"
)

func Test_parseAdditionalData(t *testing.T) {
	type args struct {
		body []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Null Additional Data",
			args: args{
				body: []byte(strings.ReplaceAll("<script type=\"text/javascript\">window.__initialDataLoaded(window._sharedData);</script><script type=\"text/javascript\">window.__additionalDataLoaded('extra',null);</script><script type=\"text/javascript\">", "\\", "")),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseAdditionalData(tt.args.body); got != tt.want {
				t.Errorf("parseAdditionalData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMagicItem_validateUrl(t *testing.T) {
	type fields struct {
		Url               string
		Type              string
		Caption           string
		DisplayUrl        string
		VideoUrl          string
		ProfilePicUrl     string
		Username          string
		Shortcode         string
		MultipleMediaList []MultipleMediaItem
	}
	type args struct {
		url string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Correct Url",
			fields: fields{
				Url:               "https://www.instagram.com/p/CYC09h9lLxh/",
				Type:              "",
				Caption:           "",
				DisplayUrl:        "",
				VideoUrl:          "",
				ProfilePicUrl:     "",
				Username:          "",
				Shortcode:         "CYC09h9lLxh",
				MultipleMediaList: []MultipleMediaItem{},
			},
			args: args{
				url: "https://www.instagram.com/p/CYC09h9lLxh/",
			},
			want: true,
		},
		{
			name: "Wrong Url",
			fields: fields{
				Url:               "https://www.instagram.com/p/CYC09h9lLxh/",
				Type:              "",
				Caption:           "",
				DisplayUrl:        "",
				VideoUrl:          "",
				ProfilePicUrl:     "",
				Username:          "",
				Shortcode:         "CYC09h9lLxh",
				MultipleMediaList: []MultipleMediaItem{},
			},
			args: args{
				url: "https://www.instagram.com/magic/CYC09h9lLxh/?utm_source=ig_web_copy_link",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mi := &MagicItem{
				Url:               tt.fields.Url,
				Type:              tt.fields.Type,
				Caption:           tt.fields.Caption,
				DisplayUrl:        tt.fields.DisplayUrl,
				VideoUrl:          tt.fields.VideoUrl,
				ProfilePicUrl:     tt.fields.ProfilePicUrl,
				Username:          tt.fields.Username,
				Shortcode:         tt.fields.Shortcode,
				MultipleMediaList: tt.fields.MultipleMediaList,
			}
			if got := mi.validateUrl(tt.args.url); got != tt.want {
				t.Errorf("MagicItem.validateUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
