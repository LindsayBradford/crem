package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

//var (
//	Grey       = color.Gray16{Y:16}
//)
//

type Crem struct{}

var _ fyne.Theme = (*Crem)(nil)

func (c Crem) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	//if name == theme.ColorNameBackground {
	//	if variant == theme.VariantLight {
	//		return color.White
	//	}
	//	return color.Black
	//}
	//
	return theme.DefaultTheme().Color(name, variant)
}

func (c Crem) Icon(name fyne.ThemeIconName) fyne.Resource {
	//if name == theme.IconNameHome {
	//	fyne.NewStaticResource("myHome", homeBytes)
	//}
	//
	return theme.DefaultTheme().Icon(name)
}

func (c Crem) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (c Crem) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
