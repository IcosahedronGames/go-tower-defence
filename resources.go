package main

import (
	"image/color"
	"strconv"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
)

const (
	backgroundColor = "131a22"

	textIdleColor     = "dff4ff"
	textDisabledColor = "5a7a91"

	labelIdleColor     = textIdleColor
	labelDisabledColor = textDisabledColor

	buttonIdleColor     = textIdleColor
	buttonDisabledColor = labelDisabledColor

	listSelectedBackground         = "4b687a"
	listDisabledSelectedBackground = "2a3944"

	listFocusedBackground = "2a3944"

	headerColor = textIdleColor

	textInputCaretColor         = "e7c34b"
	textInputDisabledCaretColor = "766326"

	toolTipColor = backgroundColor

	separatorColor = listDisabledSelectedBackground
)

type uiResources struct {
	background *image.NineSlice
	text       *textResources
	button     *buttonResources
	checkbox   *checkboxResources
	label      *labelResources
	panel      *panelResources
	textInput  *textInputResources
}

type textResources struct {
	idleColor     color.Color
	disabledColor color.Color
}

type buttonResources struct {
	image   *widget.ButtonImage
	text    *widget.ButtonTextColor
	padding widget.Insets
}

type checkboxResources struct {
	image   *widget.ButtonImage
	graphic *widget.CheckboxGraphicImage
	spacing int
}

type labelResources struct {
	text *widget.LabelColor
}

type panelResources struct {
	padding widget.Insets
}

type textInputResources struct {
	padding widget.Insets
	color   *widget.TextInputColor
}

func newUIResources() (*uiResources, error) {
	background := image.NewNineSliceColor(hexToColor(backgroundColor))

	button, err := newButtonResources()
	if err != nil {
		return nil, err
	}

	checkbox, err := newCheckboxResources()
	if err != nil {
		return nil, err
	}

	panel, err := newPanelResources()
	if err != nil {
		return nil, err
	}

	textInput, err := newTextInputResources()
	if err != nil {
		return nil, err
	}

	return &uiResources{
		background: background,
		text: &textResources{
			idleColor:     hexToColor(textIdleColor),
			disabledColor: hexToColor(textDisabledColor),
		},
		button:    button,
		label:     newLabelResources(),
		checkbox:  checkbox,
		panel:     panel,
		textInput: textInput,
	}, nil
}

func newButtonResources() (*buttonResources, error) {
	idle, err := loadImageNineSlice("assets/graphics/button-idle.png", 12, 0)
	if err != nil {
		return nil, err
	}

	hover, err := loadImageNineSlice("assets/graphics/button-hover.png", 12, 0)
	if err != nil {
		return nil, err
	}
	pressed_hover, err := loadImageNineSlice("assets/graphics/button-selected-hover.png", 12, 0)
	if err != nil {
		return nil, err
	}
	pressed, err := loadImageNineSlice("assets/graphics/button-pressed.png", 12, 0)
	if err != nil {
		return nil, err
	}

	disabled, err := loadImageNineSlice("assets/graphics/button-disabled.png", 12, 0)
	if err != nil {
		return nil, err
	}

	i := &widget.ButtonImage{
		Idle:         idle,
		Hover:        hover,
		Pressed:      pressed,
		PressedHover: pressed_hover,
		Disabled:     disabled,
	}

	return &buttonResources{
		image: i,

		text: &widget.ButtonTextColor{
			Idle:     hexToColor(buttonIdleColor),
			Disabled: hexToColor(buttonDisabledColor),
		},

		padding: widget.Insets{
			Left:  30,
			Right: 30,
		},
	}, nil
}

func newLabelResources() *labelResources {
	return &labelResources{
		text: &widget.LabelColor{
			Idle:     hexToColor(labelIdleColor),
			Disabled: hexToColor(labelDisabledColor),
		},
	}
}

func newCheckboxResources() (*checkboxResources, error) {
	idle, err := loadImageNineSlice("assets/graphics/checkbox-idle.png", 20, 0)
	if err != nil {
		return nil, err
	}

	hover, err := loadImageNineSlice("assets/graphics/checkbox-hover.png", 20, 0)
	if err != nil {
		return nil, err
	}

	disabled, err := loadImageNineSlice("assets/graphics/checkbox-disabled.png", 20, 0)
	if err != nil {
		return nil, err
	}

	checked, err := loadGraphicImages("assets/graphics/checkbox-checked-idle.png", "assets/graphics/checkbox-checked-disabled.png")
	if err != nil {
		return nil, err
	}

	unchecked, err := loadGraphicImages("assets/graphics/checkbox-unchecked-idle.png", "assets/graphics/checkbox-unchecked-disabled.png")
	if err != nil {
		return nil, err
	}

	greyed, err := loadGraphicImages("assets/graphics/checkbox-greyed-idle.png", "assets/graphics/checkbox-greyed-disabled.png")
	if err != nil {
		return nil, err
	}

	return &checkboxResources{
		image: &widget.ButtonImage{
			Idle:     idle,
			Hover:    hover,
			Pressed:  hover,
			Disabled: disabled,
		},

		graphic: &widget.CheckboxGraphicImage{
			Checked:   checked,
			Unchecked: unchecked,
			Greyed:    greyed,
		},

		spacing: 10,
	}, nil
}

func newPanelResources() (*panelResources, error) {
	return &panelResources{
		padding: widget.Insets{
			Left:   30,
			Right:  30,
			Top:    20,
			Bottom: 20,
		},
	}, nil
}
func newTextInputResources() (*textInputResources, error) {

	return &textInputResources{

		padding: widget.Insets{
			Left:   8,
			Right:  8,
			Top:    4,
			Bottom: 4,
		},

		color: &widget.TextInputColor{
			Idle:          hexToColor(textIdleColor),
			Disabled:      hexToColor(textDisabledColor),
			Caret:         hexToColor(textInputCaretColor),
			DisabledCaret: hexToColor(textInputDisabledCaretColor),
		},
	}, nil
}

func hexToColor(h string) color.Color {
	u, err := strconv.ParseUint(h, 16, 0)
	if err != nil {
		panic(err)
	}

	return color.NRGBA{
		R: uint8(u & 0xff0000 >> 16),
		G: uint8(u & 0xff00 >> 8),
		B: uint8(u & 0xff),
		A: 255,
	}
}
