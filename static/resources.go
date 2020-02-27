package static

import "fyne.io/fyne/theme"

var (
	iconSadEmoji *theme.ThemedResource
)

func init() {
	iconSadEmoji = theme.NewThemedResource(iconEmojisadPng, nil)
}

func SadEmojiIcon() *theme.ThemedResource {
	return iconSadEmoji
}
