package ui

// Layout holds calculated dimensions for the UI split.
type Layout struct {
	TotalWidth  int
	TotalHeight int

	TabBarHeight   int
	StatusBarHeight int

	ContentWidth  int
	ContentHeight int

	TableWidth  int
	PanelWidth  int
	PanelOpen   bool
}

const (
	tabBarHeight    = 2
	statusBarHeight = 1
	panelRatio      = 0.38
	minPanelWidth   = 30
)

// Calculate computes layout dimensions.
func Calculate(width, height int, panelOpen bool) Layout {
	l := Layout{
		TotalWidth:      width,
		TotalHeight:     height,
		TabBarHeight:    tabBarHeight,
		StatusBarHeight: statusBarHeight,
		PanelOpen:       panelOpen,
	}

	l.ContentWidth = width
	l.ContentHeight = height - tabBarHeight - statusBarHeight

	if l.ContentHeight < 3 {
		l.ContentHeight = 3
	}

	if panelOpen && width > 80 {
		l.PanelWidth = int(float64(width) * panelRatio)
		if l.PanelWidth < minPanelWidth {
			l.PanelWidth = minPanelWidth
		}
		l.TableWidth = width - l.PanelWidth - 1 // 1 for border
	} else {
		l.TableWidth = width
		l.PanelWidth = 0
		l.PanelOpen = false
	}

	return l
}
