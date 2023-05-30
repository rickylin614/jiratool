package lib

import (
	"image/color"

	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

//go:embed setofont.ttf
var ttfFile []byte

// 自定義字體
type MyTheme struct{}

func (*MyTheme) Font(s fyne.TextStyle) fyne.Resource {
	if len(ttfFile) == 0 {
		// 載入失敗，使用預設字體
		return theme.DefaultTheme().Font(s)
	}
	fontData := fyne.NewStaticResource("setofont.ttf", ttfFile)
	return fontData
}

func (*MyTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}

func (*MyTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (*MyTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}
