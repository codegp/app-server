package main

import (
	"bytes"
	"image/png"
	"net/http"

	"github.com/nfnt/resize"
)

func PostIcon(w http.ResponseWriter, r *http.Request) *requestError {
	typeID, rerr := readIDFromRequest(r, "typeID")
	if rerr != nil {
		return rerr
	}

	f, _, err := r.FormFile("iconData")
	if err != nil {
		return requestErrorf(err, "Error getting icon from form body")
	}

	ogImg, err := png.Decode(f)
	if err != nil {
		return requestErrorf(err, "Error decoding image")
	}

	resizedImage := resize.Resize(100, 100, ogImg, resize.Lanczos3)
	buf := new(bytes.Buffer)
	err = png.Encode(buf, resizedImage)
	if err != nil {
		return requestErrorf(err, "Error encoding resized image")
	}

	err = cp.WriteIcon(typeID, buf.Bytes())
	if err != nil {
		return requestErrorf(err, "Error writing icon")
	}

	return nil
}
