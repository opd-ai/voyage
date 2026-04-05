// Package input provides unified input handling for keyboard, mouse, and touch.
//
// This package abstracts away the differences between desktop (keyboard/mouse)
// and mobile/touch (tap/swipe) input methods, providing a consistent interface
// for game logic.
//
// # Touch Support
//
// Touch input is fully supported for mobile browsers and devices:
//   - Tap: equivalent to mouse click or Enter/Space key
//   - Swipe: equivalent to arrow key movement
//   - Multi-touch: not currently used
//
// # Usage
//
// Create an InputManager and call Update() each frame:
//
//	manager := input.NewManager()
//
//	func (g *Game) Update() error {
//	    manager.Update()
//
//	    if manager.JustConfirmed() {
//	        // Handle confirm action
//	    }
//
//	    if dir := manager.GetDirection(); dir != DirectionNone {
//	        // Handle movement
//	    }
//
//	    return nil
//	}
package input
