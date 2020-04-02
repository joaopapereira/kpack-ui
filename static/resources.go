package static

import "fyne.io/fyne/theme"

var (
	iconSadEmoji *theme.ThemedResource
	appIcon      *theme.ThemedResource
	appIcon16    *theme.ThemedResource
	appIcon32    *theme.ThemedResource
)

func init() {
	//nolint:ST1003
	iconSadEmoji = theme.NewThemedResource(iconEmojisadPng, nil)
	//nolint:ST1003
	appIcon = theme.NewThemedResource(iconKpackUiIco, nil)
	//nolint:ST1003
	appIcon16 = theme.NewThemedResource(iconKpackUi16x16Png, nil)
	//nolint:ST1003
	appIcon32 = theme.NewThemedResource(iconKpackUi32x32Png, nil)
}

func SadEmojiIcon() *theme.ThemedResource {
	return iconSadEmoji
}
