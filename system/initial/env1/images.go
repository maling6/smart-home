package env1

import (
	"os"
	"strings"
	"path"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	. "github.com/e154/smart-home/system/initial/assertions"
	m "github.com/e154/smart-home/models"
)

func images(adaptors *adaptors.Adaptors) (imageList map[string]*m.Image) {

	imageList = map[string]*m.Image{
		"button_v1_off": {
			Image:    "30d2f4116a09fd14b49c266985db8109.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     2518,
			Name:     "button_v1_off.svg",
		},
		"button_v1_refresh": {
			Image:    "86486ca5d086aafd5724d61251b94bba.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     3212,
			Name:     "button_v1_refresh.svg",
		},
		"lamp_v1_r": {
			Image:    "2d4a761241e24a77725287180656b466.svg",
			MimeType: "text/xml; charset=utf-8",
			Size:     2261,
			Name:     "lamp_v1_r.svg",
		},
		"socket_v1_b": {
			Image:    "bef910d70c56f38b22cea0c00d92d8cc.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     7326,
			Name:     "socket_v1_b.svg",
		},
		"button_v1_on": {
			Image:    "7c145f62dcaf8da2a9eb43f2b23ea2b1.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     2398,
			Name:     "button_v1_on.svg",
		},
		"socket_v1_def": {
			Image:    "4c28edf0700531731df43ed055ebf56d.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     7326,
			Name:     "socket_v1_def.svg",
		},
		"socket_v1_r": {
			Image:    "e91e461f7c9a800eed5a074101d3e5a5.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     7326,
			Name:     "socket_v1_r.svg",
		},
		"lamp_v1_def": {
			Image:    "91e93ee7e7734654083dee0a5cbe55e9.svg",
			MimeType: "text/xml; charset=utf-8",
			Size:     2266,
			Name:     "lamp_v1_def.svg",
		},
		"socket_v1_g": {
			Image:    "4819b36056dfa786f5856fa45e9a3151.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     7326,
			Name:     "socket_v1_g.svg",
		},
		"lamp_v1_y": {
			Image:    "c1c5ec4e75bb6ec33f5f8cfd87b0090e.svg",
			MimeType: "text/xml; charset=utf-8",
			Size:     2261,
			Name:     "lamp_v1_y.svg",
		},
		"socket_v2_b": {
			Image:    "c813ac54bb4dd6b99499d097eda67310.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     3060,
			Name:     "socket_v2_b.svg",
		},
		"socket_v2_def": {
			Image:    "f0ea38f2b388dc2bb2566f6efc7731b0.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     3060,
			Name:     "socket_v2_def.svg",
		},
		"socket_v2_g": {
			Image:    "fa6b42c81056069d03857cfbb2cf95eb.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     3060,
			Name:     "socket_v2_g.svg",
		},
		"socket_v2_r": {
			Image:    "e565f191030491cfdc39ad728559c18f.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     3060,
			Name:     "socket_v2_r.svg",
		},
		"socket_v3_b": {
			Image:    "297d56426098a53091fb8f91aabe3cd7.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     2718,
			Name:     "socket_v3_b.svg",
		},
		"socket_v3_def": {
			Image:    "becf0f8f635061c143acb4329f744615.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     2718,
			Name:     "socket_v3_def.svg",
		},
		"socket_v3_g": {
			Image:    "850bf4da00cb9de85e1442695230a127.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     2718,
			Name:     "socket_v3_g.svg",
		},
		"socket_v3_r": {
			Image:    "434514389e95cab6d684b978378055d5.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     2718,
			Name:     "socket_v3_r.svg",
		},
	}

	var err error
	var subDir string
	for _, image := range imageList {
		image.Id, err = adaptors.Image.Add(image)
		So(err, ShouldBeNil)

		fullPath := common.GetFullPath(image.Image)
		to := path.Join(fullPath, image.Image)
		if exist := common.FileExist(to); !exist {

			os.MkdirAll(fullPath, os.ModePerm)

			switch {
			case strings.Contains(image.Name, "button"):
				subDir = "buttons"
			case strings.Contains(image.Name, "lamp"):
				subDir = "lamp"
			case strings.Contains(image.Name, "socket"):
				subDir = "socket"
			}

			from := path.Join("data", "icons", subDir, image.Name)
			common.CopyFile(from, to)
		}
	}

	return
}
