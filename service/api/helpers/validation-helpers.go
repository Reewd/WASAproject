package helpers

import (
	"errors"
	"strings"

	"github.com/Reewd/WASAproject/service/api/dto"
	"github.com/rivo/uniseg"
	"github.com/ucarion/emoji"
)

func IsSingleEmoji(s string) error {
	s = strings.TrimSpace(s)

	// 1) must be exactly one grapheme-cluster
	if uniseg.GraphemeClusterCount(s) != 1 {
		return errors.New("reaction must contain exactly one emoji")
	}

	// 2) that cluster must be recognised as an emoji
	if _, ok := emoji.Lookup(s); !ok {
		return errors.New("reaction must be an emoji")
	}

	return nil
}

func ExtractPhoto(reqPhoto *dto.Photo) (*string, *dto.Photo) {
	if reqPhoto != nil {
		return &reqPhoto.PhotoId, reqPhoto
	}
	return nil, nil
}
