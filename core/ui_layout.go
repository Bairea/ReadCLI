package core

func applyLayoutFromApp() {
	applyLayoutFromAppWithReflow(true)
}

func applyLayoutFromAppWithoutReflow() {
	applyLayoutFromAppWithReflow(false)
}

func applyLayoutFromAppWithReflow(reflow bool) {
	w, h := lastTermWidth, lastTermHeight
	if w <= 0 {
		w = fixedWidth
	}
	if h <= 0 {
		h = 24
	}
	applyLayoutWithReflow(w, h, reflow)
}

func applyLayout(termWidth, termHeight int) {
	applyLayoutWithReflow(termWidth, termHeight, true)
}

func applyLayoutWithReflow(termWidth, termHeight int, reflow bool) {
	if termHeight <= 0 {
		termHeight = 3
	}
	width := termWidth
	if width <= 0 {
		width = fixedWidth
	}

	focusLevel := effectiveFocusLevel()
	leftWidth := 24
	rightWidth := 28
	headerHeight := 4
	footerHeight := 4
	if focusLevel >= 1 {
		leftWidth = 0
		rightWidth = 0
	}
	if focusLevel >= 2 {
		headerHeight = 0
		footerHeight = 1
	}
	if focusLevel == 0 && width < 120 {
		leftWidth = 20
		rightWidth = 22
	}
	if focusLevel == 0 && width < 90 {
		leftWidth = 18
		rightWidth = 0
	}
	mainWidth := width - leftWidth - rightWidth
	if mainWidth < 30 {
		mainWidth = width
		leftWidth = 0
		rightWidth = 0
	}

	if focusLevel < 2 && termHeight < 16 {
		headerHeight = 4
		footerHeight = 3
	}

	contentHeight := termHeight - headerHeight - footerHeight
	if contentHeight < 1 {
		contentHeight = 1
	}

	if focusLevel >= 2 {
		mainContentWidth = mainWidth
		mainContentHeight = contentHeight
	} else {
		borderPad := 0
		if app != nil && app.uiState.showBorder {
			borderPad = 2
		}
		mainContentWidth = mainWidth - borderPad
		mainContentHeight = contentHeight - borderPad
	}

	nextContentWidth := readingContentWidth(mainWidth)
	if reflow {
		app.readingState.contentWidth = nextContentWidth
	}
	if reflow && app.reader != nil {
		app.reader.Reflow(app.readingState.contentWidth)
	}

	root.ResizeItem(header, headerHeight, 0)
	root.ResizeItem(midRow, contentHeight, 1)
	root.ResizeItem(footer, footerHeight, 0)

	if leftWidth > 0 {
		midRow.ResizeItem(left, leftWidth, 0)
	} else {
		midRow.ResizeItem(left, 0, 0)
	}
	midRow.ResizeItem(main, 0, 1)
	if rightWidth > 0 {
		midRow.ResizeItem(right, rightWidth, 0)
	} else {
		midRow.ResizeItem(right, 0, 0)
	}

	refreshChrome()
}

// effectiveFocusLevel resolves the focus level that applies to the current scene.
// Reader surfaces keep the configured level, overlays opened from reading drop
// back to the half level so the header and key hints stay visible, and every
// other scene uses the full layout.
func effectiveFocusLevel() int {
	if app == nil || app.bossKey {
		return 0
	}
	level := app.readingState.focusLevel
	if level < 1 {
		return 0
	}
	switch app.mode {
	case modeReading:
		return level
	case modeTOC, modeBookmarks, modeReadingSettings, modeSearchInput, modeReadingJumpInput, modeReadingColorInput:
		return 1
	default:
		return 0
	}
}

func compactReadingUI() bool {
	return effectiveFocusLevel() >= 2
}
