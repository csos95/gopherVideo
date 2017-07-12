package gopherVideo

import (
	"fmt"

	"github.com/gopherjs/gopherjs/js"

	"honnef.co/go/js/dom"
)

// NewPlayer returns a new gopher video player and the contained video
func (p *Player) setupHTML() {

	// div container for the video and controls
	container := document.CreateElement("div").(*dom.HTMLDivElement)
	container.SetClass("gopherVideo")
	container.SetID(fmt.Sprintf("%s", p.ID))

	// the video
	video := document.CreateElement("video").(*dom.HTMLVideoElement)
	video.SetClass("gopherVideo-video")

	// the source for the video
	source := document.CreateElement("source").(*dom.HTMLSourceElement)
	source.SetAttribute("src", p.URL)
	video.AppendChild(source)
	container.AppendChild(video)

	// div for the controls
	controls := document.CreateElement("div").(*dom.HTMLDivElement)
	controls.SetClass("gopherVideo-controls")

	bottomControls := document.CreateElement("div").(*dom.HTMLDivElement)
	bottomControls.SetClass("gopherVideo-bottom-controls")

	// a button to play/pause the video
	object := js.Global.Get("document").Call("createElementNS", "http://www.w3.org/2000/svg", "svg")
	playpause := objectToBasicHTMLElement(object)
	playpause.SetAttribute("xmlns", "http://www.w3.org/2000/svg")
	playpause.SetAttribute("class", "gopherVideo-playpause")
	playpause.SetAttribute("width", "20px")
	playpause.SetAttribute("height", "20px")
	playpause.SetAttribute("viewBox", "0 0 8 8")

	object = js.Global.Get("document").Call("createElementNS", "http://www.w3.org/2000/svg", "path")
	playpausePath := objectToBasicHTMLElement(object)
	playpausePath.SetAttribute("d", "M0 0v6l6-3-6-3z")
	playpausePath.SetAttribute("transform", "translate(1 1)")
	playpause.AppendChild(playpausePath)
	bottomControls.AppendChild(playpause)

	// the current playtime text
	timeText := document.CreateElement("pre").(*dom.HTMLPreElement)
	timeText.SetClass("gopherVideo-time")
	timeText.SetTextContent("0:00")
	bottomControls.AppendChild(timeText)

	// the progress bar for the video
	progressBar := document.CreateElement("input").(*dom.HTMLInputElement)
	progressBar.SetClass("gopherVideo-progressbar")
	progressBar.SetAttribute("type", "range")
	progressBar.SetAttribute("min", "0")
	progressBar.SetAttribute("max", "1")
	progressBar.Value = "0"
	bottomControls.AppendChild(progressBar)

	// the video duration text
	durationText := document.CreateElement("pre").(*dom.HTMLPreElement)
	durationText.SetClass("gopherVideo-duration")
	durationText.SetTextContent("0:00")
	bottomControls.AppendChild(durationText)

	// a button to enter fullscreen
	fullscreen := document.CreateElement("button").(*dom.HTMLButtonElement)
	fullscreen.SetClass("gopherVideo-fullscreen")
	fullscreen.SetTextContent("fullscreen")
	bottomControls.AppendChild(fullscreen)

	controls.AppendChild(bottomControls)
	container.AppendChild(controls)

	p.Container = container
	p.Video = video
	p.Controls = controls
	p.PlayPause = playpause
	p.ProgressBar = progressBar
	p.TimeText = timeText
	p.DurationText = durationText
	p.FullscreenButton = fullscreen
}

func objectToBasicHTMLElement(object *js.Object) dom.Element {
	return &dom.BasicHTMLElement{&dom.BasicElement{&dom.BasicNode{object}}}
}
