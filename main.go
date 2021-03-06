package gophervideo

import (
	"fmt"

	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

var document = dom.GetWindow().Document()
var documentElement = js.Global.Get("document")
var body = document.DocumentElement().GetElementsByTagName("body")[0].(*dom.HTMLBodyElement)
var cssSet = false

// Player represents a gopher video player
type Player struct {
	// player data
	ID               string
	URL              string
	Duration         int
	ProgressBarWidth int
	TimeTextWidth    int
	Fullscreen       bool
	FirstPlay        bool
	Removed          bool
	Seeking          bool
	MouseInContainer bool
	MouseMoved       bool
	SecondsSinceMove int

	// player elements
	Parent           dom.HTMLElement
	Container        *dom.HTMLDivElement
	Video            *dom.HTMLVideoElement
	Controls         *dom.HTMLDivElement
	PlayPause        *dom.BasicHTMLElement
	ProgressBarBack  *dom.HTMLDivElement
	ProgressBarFront *dom.HTMLDivElement
	TimeText         *dom.HTMLSpanElement
	DurationText     *dom.HTMLSpanElement
	VolumeIcon       *dom.BasicHTMLElement
	VolumeBar        *dom.HTMLInputElement
	FullscreenButton *dom.BasicHTMLElement

	// listeners
	playpauseListener           func(*js.Object)
	playpauseMouseListener      func(*js.Object)
	videoTimeUpdateListener     func(*js.Object)
	ProgressBarClickListener    func(*js.Object)
	ProgressBarDragListener     func(*js.Object)
	ProgressBarDownListener     func(*js.Object)
	ProgressBarUpListener       func(*js.Object)
	volumeBarListener           func(*js.Object)
	fullscreenButtonListener    func(*js.Object)
	fullscreenMouseListener     func(*js.Object)
	fullscreenListener          func(*js.Object)
	webkitFullscreenListener    func(*js.Object)
	mozillaFullscreenListener   func(*js.Object)
	microsoftFullscreenListener func(*js.Object)
	keyPressListener            func(*js.Object)
}

// NewPlayer returns a new gopher video player and the contained video
func NewPlayer(parent dom.HTMLElement, url string) *Player {
	id := "1"

	player := &Player{
		ID:         id,
		URL:        url,
		Parent:     parent,
		Fullscreen: false,
		FirstPlay:  true,
		Removed:    false,
		Seeking:    false,
	}

	if !cssSet {
		player.setupCSS()
	}
	player.setupHTML()
	player.setupListeners()

	player.Parent.AppendChild(player.Container)

	return player
}

// Remove the player from the document
func (p *Player) Remove() {
	// remove all listeners
	p.PlayPause.RemoveEventListener("click", true, p.playpauseListener)
	p.PlayPause.RemoveEventListener("click", true, p.playpauseMouseListener)
	p.Video.RemoveEventListener("timeupdate", false, p.videoTimeUpdateListener)
	p.Video.RemoveEventListener("click", false, p.ProgressBarClickListener)
	p.Video.RemoveEventListener("mousemove", false, p.ProgressBarDragListener)
	p.Video.RemoveEventListener("mousedown", false, p.ProgressBarDownListener)
	p.Video.RemoveEventListener("mouseup", false, p.ProgressBarDownListener)
	p.VolumeBar.RemoveEventListener("input", true, p.volumeBarListener)
	p.FullscreenButton.RemoveEventListener("click", true, p.fullscreenButtonListener)
	p.FullscreenButton.RemoveEventListener("click", true, p.fullscreenMouseListener)
	document.RemoveEventListener("fullscreenchange", false, p.fullscreenListener)
	document.RemoveEventListener("webkitfullscreenchange", false, p.webkitFullscreenListener)
	document.RemoveEventListener("mozfullscreenchange", false, p.mozillaFullscreenListener)
	document.RemoveEventListener("MSFullscreenChange", false, p.microsoftFullscreenListener)
	document.RemoveEventListener("keypress", false, p.keyPressListener)

	p.Pause()
	p.Video.SetAttribute("src", "")
	p.Video.Call("load")
	p.Container.RemoveChild(p.Video)
	p.Parent.RemoveChild(p.Container)

	p.Removed = true
}

// formats the time in days:hours:minutes:seconds leaving off empty fields to the left
func (p *Player) timeFormat(seconds int) string {
	if p.Duration < 60 {
		return fmt.Sprintf("%d", seconds)
	} else if p.Duration < 3600 {
		return fmt.Sprintf("%d:%02d", seconds/60, seconds%60)
	} else if p.Duration < 86400 {
		return fmt.Sprintf("%d:%02d:%02d", seconds/3600, seconds/60%60, seconds%60)
	}
	return fmt.Sprintf("%d:%02d:%02d:%02d", seconds/86400, seconds/3600%24, seconds/60%60, seconds%60)
}
