// Copyright 2018 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// This file was modified from the example found: https://github.com/hajimehoshi/ebiten/blob/main/examples/tiles/main.go with permission from the author
package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"

	"github.com/ebitenui/ebitenui"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	screenWidth  = 920
	screenHeight = 920
	tileSize     = 16
	tileMapWidth = 15
	title        = "Icosahedron Games: Tower Defense"
)

func main() {
	g := &Game{
		layers:     getLayers(),
		tilesImage: getTileImage(),
		settings: &Settings{
			showFPS: false,
			vSynch:  ebiten.IsVsyncEnabled(),
		},
		player: NewPlayer(),
	}
	g.ui = g.getEbitenUI()
	v := mgl32.Vec2{}
	fmt.Printf("%f\n", v[0])

	ebiten.SetTPS(60)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle(title)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

type Player struct {
	position mgl32.Vec2
}

func NewPlayer() Player {
	return Player{
		mgl32.Vec2{0, 0},
	}
}

func (player *Player) UpdatePlayer(deltaTime float32) {
	// Game units / second
	var movementSpeed float32 = 100

	var movementDir = mgl32.Vec2{0, 0}

	if inpututil.KeyPressDuration(ebiten.KeyW) > 0 {
		movementDir = movementDir.Add(mgl32.Vec2{0, -1})
	}
	if inpututil.KeyPressDuration(ebiten.KeyS) > 0 {
		movementDir = movementDir.Add(mgl32.Vec2{0, 1})
	}
	if inpututil.KeyPressDuration(ebiten.KeyA) > 0 {
		movementDir = movementDir.Add(mgl32.Vec2{-1, 0})
	}
	if inpututil.KeyPressDuration(ebiten.KeyD) > 0 {
		movementDir = movementDir.Add(mgl32.Vec2{1, 0})
	}

	if movementDir[0] != float32(0) || movementDir[1] != float32(0) {
		movementDir = movementDir.Normalize()
		var frameDisplacement = movementDir.Mul(movementSpeed * deltaTime)
		player.position = player.position.Add(frameDisplacement)
		log.Println("Position {", player.position[0], ",", player.position[1], "}")
	}
}

// Enum of windows that can be open
type Window string

const (
	MainMenu Window = "mainMenu"
	None     Window = "none"
)

type PerFrame struct {
	deltaTime32 float32
	deltaTime64 float64
}

// God class
type Game struct {
	tilesImage *ebiten.Image
	layers     [][]int

	ui        *ebitenui.UI
	headerLbl *widget.Text
	settings  *Settings
	window    Window
	player    Player
	perFrame  PerFrame
}

type Settings struct {
	showFPS bool
	vSynch  bool
}

func (g *Game) Update() error {
	// Ensure that the UI is updated to receive events
	g.ui.Update()
	g.perFrame.deltaTime64 = 1.0 / ebiten.ActualTPS()
	g.perFrame.deltaTime64 = max(0.001, g.perFrame.deltaTime64)
	g.perFrame.deltaTime32 = float32(g.perFrame.deltaTime64)
	g.player.UpdatePlayer(g.perFrame.deltaTime32)

	// Update the Label text to indicate if the ui is currently being hovered over or not
	g.headerLbl.Label = fmt.Sprintf("Game Demo!\nUI is hovered: %t", input.UIHovered)

	// Log out if we have clicked on the gamefield and NOT the ui
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && !input.UIHovered {
		log.Println("Mouse clicked on gamefield")
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		log.Println("Escape is pressed")
		if g.window != MainMenu {
			g.window = MainMenu
			openMainMenu(g)
		}
	}

	return nil
}

// Open the main menu, triggered by pressing escape key
func openMainMenu(g *Game) {
	res, _ := newUIResources()
	var rw widget.RemoveWindowFunc
	var window *widget.Window

	titleFace, _ := loadFont(24)
	face, _ := loadFont(20)

	titleBar := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(res.background),
		widget.ContainerOpts.Layout(widget.NewGridLayout(widget.GridLayoutOpts.Columns(3), widget.GridLayoutOpts.Stretch([]bool{true, false, false}, []bool{true}), widget.GridLayoutOpts.Padding(widget.Insets{
			Left:   30,
			Right:  5,
			Top:    6,
			Bottom: 5,
		}))))

	titleBar.AddChild(widget.NewText(
		widget.TextOpts.Text("Main Menu", titleFace, res.textInput.color.Idle),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
	))

	titleBar.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(res.button.image),
		widget.ButtonOpts.TextPadding(res.button.padding),
		widget.ButtonOpts.Text("X", face, res.button.text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			g.window = None
			rw()
		}),
		widget.ButtonOpts.TabOrder(99),
	))

	c := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(res.background),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
				widget.GridLayoutOpts.Padding(res.panel.padding),
				widget.GridLayoutOpts.Spacing(0, 15),
			),
		),
	)

	// Show FPS setting
	cb1 := widget.NewLabeledCheckbox(
		widget.LabeledCheckboxOpts.Spacing(res.checkbox.spacing),
		widget.LabeledCheckboxOpts.CheckboxOpts(
			widget.CheckboxOpts.InitialState(boolToCheck(g.settings.showFPS)),
			widget.CheckboxOpts.ButtonOpts(widget.ButtonOpts.Image(res.checkbox.image)),
			widget.CheckboxOpts.Image(res.checkbox.graphic),
			widget.CheckboxOpts.StateChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
				if g.settings.showFPS {
					g.settings.showFPS = false
				} else {
					g.settings.showFPS = true
				}
			})),
		widget.LabeledCheckboxOpts.LabelOpts(widget.LabelOpts.Text("Show FPS", face, res.label.text)))

	c.AddChild(cb1)

	// VSync
	cb2 := widget.NewLabeledCheckbox(
		widget.LabeledCheckboxOpts.Spacing(res.checkbox.spacing),
		widget.LabeledCheckboxOpts.CheckboxOpts(
			widget.CheckboxOpts.InitialState(boolToCheck(g.settings.vSynch)),
			widget.CheckboxOpts.ButtonOpts(widget.ButtonOpts.Image(res.checkbox.image)),
			widget.CheckboxOpts.Image(res.checkbox.graphic),
			widget.CheckboxOpts.StateChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
				if g.settings.vSynch {
					g.settings.vSynch = false
				} else {
					g.settings.vSynch = true
				}
			})),
		widget.LabeledCheckboxOpts.LabelOpts(widget.LabelOpts.Text("VSynch", face, res.label.text)))

	c.AddChild(cb2)

	bc := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(15),
		)),
	)
	c.AddChild(bc)

	window = widget.NewWindow(
		widget.WindowOpts.Modal(),
		widget.WindowOpts.Contents(c),
		widget.WindowOpts.TitleBar(titleBar, 30),
		widget.WindowOpts.Draggable(),
		widget.WindowOpts.Resizeable(),
		widget.WindowOpts.MinSize(500, 200),
		widget.WindowOpts.MaxSize(700, 400),
		widget.WindowOpts.ResizeHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Resize: ", args.Rect)
		}),
		widget.WindowOpts.MoveHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Move: ", args.Rect)
		}),
	)
	windowSize := input.GetWindowSize()
	r := image.Rect(0, 0, 550, 250)
	r = r.Add(image.Point{windowSize.X / 4 / 2, windowSize.Y * 2 / 3 / 2})
	window.SetLocation(r)

	rw = g.ui.AddWindow(window)
}

func boolToCheck(test bool) widget.WidgetState {
	if test {
		return widget.WidgetChecked
	}
	return widget.WidgetUnchecked
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the tilemap
	g.drawGameWorld(screen)
	// Ensure ui.Draw is called after the gameworld is drawn
	g.ui.Draw(screen)
	// Print FPS on screen
	if g.settings.showFPS {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f", ebiten.ActualFPS()))
	}
	ebiten.SetVsyncEnabled(g.settings.vSynch)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *Game) getEbitenUI() *ebitenui.UI {
	// load label text font
	face, _ := loadFont(18)

	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(5)))),
	)

	// Because this container has a backgroundImage set we track that the ui is hovered over.
	headerContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(color.Black)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  true,
				StretchVertical:    false,
			}),
			// Uncomment this to not track that you are hovering over this header
			// widget.WidgetOpts.TrackHover(false),
		),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	g.headerLbl = widget.NewText(
		widget.TextOpts.Text("Game Demo!", face, color.White),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
			// Uncomment to force tracking hover of this element
			// widget.WidgetOpts.TrackHover(true),
		),
	)
	headerContainer.AddChild(g.headerLbl)
	rootContainer.AddChild(headerContainer)

	hProgressbar := widget.NewProgressBar(
		widget.ProgressBarOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  true,
				StretchVertical:    false,
			}),
			// Set the minimum size for the progress bar.
			// This is necessary if you wish to have the progress bar be larger than
			// the provided track image. In this exampe since we are using NineSliceColor
			// which is 1px x 1px we must set a minimum size.
			widget.WidgetOpts.MinSize(200, 20),
			// Set this parameter to indicate we want do not want to track that this ui element is being hovered over.
			// widget.WidgetOpts.TrackHover(false),
		),
		widget.ProgressBarOpts.Images(
			// Set the track images (Idle, Disabled).
			&widget.ProgressBarImage{
				Idle: eimage.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			},
			// Set the progress images (Idle, Disabled).
			&widget.ProgressBarImage{
				Idle: eimage.NewNineSliceColor(color.NRGBA{255, 255, 100, 255}),
			},
		),
		// Set the min, max, and current values.
		widget.ProgressBarOpts.Values(0, 10, 7),
		// Set how much of the track is displayed when the bar is overlayed.
		widget.ProgressBarOpts.TrackPadding(widget.Insets{
			Top:    2,
			Bottom: 2,
		}),
	)

	rootContainer.AddChild(hProgressbar)

	// Create a label to show the percentage on top of the progress bar
	label2 := widget.NewText(
		widget.TextOpts.Text("70%", face, color.Black),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
	)
	rootContainer.AddChild(label2)

	return &ebitenui.UI{
		Container: rootContainer,
	}
}

func (g *Game) drawGameWorld(screen *ebiten.Image) {
	w := g.tilesImage.Bounds().Dx()
	tileXCount := w / tileSize

	for _, l := range g.layers {
		for i, t := range l {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64((i%tileMapWidth)*tileSize), float64((i/tileMapWidth)*tileSize))
			op.GeoM.Translate(float64(-g.player.position[0]), float64(-g.player.position[1]))
			op.GeoM.Scale(4, 4)

			sx := (t % tileXCount) * tileSize
			sy := (t / tileXCount) * tileSize
			screen.DrawImage(g.tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
		}
	}
}

func getTileImage() *ebiten.Image {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(images.Tiles_png))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}

func getLayers() [][]int {
	return [][]int{
		{
			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 218, 243, 244, 243,
			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,

			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			243, 243, 244, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			243, 243, 243, 243, 243, 243, 243, 243, 243, 219, 243, 243, 243, 219, 243,
			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,

			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 244, 243, 243, 243,
			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
		},
		{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 26, 27, 28, 29, 30, 31, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 51, 52, 53, 54, 55, 56, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 76, 77, 78, 79, 80, 81, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 101, 102, 103, 104, 105, 106, 0, 0, 0, 0,

			0, 0, 0, 0, 0, 126, 127, 128, 129, 130, 131, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 303, 303, 245, 242, 303, 303, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,

			0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 245, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
		},
	}
}

func loadFont(size float64) (font.Face, error) {
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}
