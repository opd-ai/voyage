//go:build headless

package rendering

import (
	"image/color"
	"testing"
)

func TestNewPortraitGenerator(t *testing.T) {
	pg := NewPortraitGenerator(12345, 32)
	if pg == nil {
		t.Fatal("NewPortraitGenerator returned nil")
	}
	if pg.portraitSize != 32 {
		t.Errorf("expected portraitSize 32, got %d", pg.portraitSize)
	}
}

func TestGenerateAnimatedPortrait(t *testing.T) {
	pg := NewPortraitGenerator(12345, 32)
	primary := color.RGBA{50, 50, 150, 255}
	secondary := color.RGBA{100, 80, 60, 255}
	skin := color.RGBA{200, 160, 140, 255}

	portrait := pg.GenerateAnimatedPortrait(primary, secondary, skin)
	if portrait == nil {
		t.Fatal("GenerateAnimatedPortrait returned nil")
	}

	if portrait.IdleFrameCount != 4 {
		t.Errorf("expected 4 idle frames, got %d", portrait.IdleFrameCount)
	}
	if portrait.HurtFrameCount != 4 {
		t.Errorf("expected 4 hurt frames, got %d", portrait.HurtFrameCount)
	}
	if portrait.DeathFrameCount != 8 {
		t.Errorf("expected 8 death frames, got %d", portrait.DeathFrameCount)
	}
}

func TestPortraitAnimationStates(t *testing.T) {
	pg := NewPortraitGenerator(12345, 32)
	primary := color.RGBA{50, 50, 150, 255}
	secondary := color.RGBA{100, 80, 60, 255}
	skin := color.RGBA{200, 160, 140, 255}

	portrait := pg.GenerateAnimatedPortrait(primary, secondary, skin)

	if portrait.State() != PortraitIdle {
		t.Errorf("initial state should be PortraitIdle, got %d", portrait.State())
	}

	portrait.SetState(PortraitHurt)
	if portrait.State() != PortraitHurt {
		t.Errorf("state should be PortraitHurt, got %d", portrait.State())
	}

	portrait.SetState(PortraitDeath)
	if portrait.State() != PortraitDeath {
		t.Errorf("state should be PortraitDeath, got %d", portrait.State())
	}
}

func TestPortraitAnimationUpdate(t *testing.T) {
	pg := NewPortraitGenerator(12345, 32)
	primary := color.RGBA{50, 50, 150, 255}
	secondary := color.RGBA{100, 80, 60, 255}
	skin := color.RGBA{200, 160, 140, 255}

	portrait := pg.GenerateAnimatedPortrait(primary, secondary, skin)

	if portrait.currentFrame != 0 {
		t.Errorf("initial frame should be 0, got %d", portrait.currentFrame)
	}

	// FrameTime is 0.25, update by less than that
	portrait.Update(0.1)
	if portrait.currentFrame != 0 {
		t.Errorf("frame should still be 0 after 0.1s, got %d", portrait.currentFrame)
	}

	// Advance enough to change frame
	portrait.Update(0.2)
	if portrait.currentFrame != 1 {
		t.Errorf("frame should be 1 after 0.3s total, got %d", portrait.currentFrame)
	}
}

func TestPortraitIdleLooping(t *testing.T) {
	pg := NewPortraitGenerator(12345, 32)
	primary := color.RGBA{50, 50, 150, 255}
	secondary := color.RGBA{100, 80, 60, 255}
	skin := color.RGBA{200, 160, 140, 255}

	portrait := pg.GenerateAnimatedPortrait(primary, secondary, skin)

	// Advance through all 4 frames and then some
	for i := 0; i < 8; i++ {
		portrait.Update(0.26) // Slightly more than frame time
	}

	// Should have looped back
	if portrait.currentFrame >= 4 {
		t.Errorf("idle animation should loop, frame %d >= 4", portrait.currentFrame)
	}
}

func TestPortraitDeathNoLoop(t *testing.T) {
	pg := NewPortraitGenerator(12345, 32)
	primary := color.RGBA{50, 50, 150, 255}
	secondary := color.RGBA{100, 80, 60, 255}
	skin := color.RGBA{200, 160, 140, 255}

	portrait := pg.GenerateAnimatedPortrait(primary, secondary, skin)
	portrait.SetState(PortraitDeath)

	// Advance through all 8 death frames and beyond
	for i := 0; i < 15; i++ {
		portrait.Update(0.26)
	}

	// Should stay at last frame
	if portrait.currentFrame != 7 {
		t.Errorf("death animation should stop at last frame (7), got %d", portrait.currentFrame)
	}
}

func TestPortraitReset(t *testing.T) {
	pg := NewPortraitGenerator(12345, 32)
	primary := color.RGBA{50, 50, 150, 255}
	secondary := color.RGBA{100, 80, 60, 255}
	skin := color.RGBA{200, 160, 140, 255}

	portrait := pg.GenerateAnimatedPortrait(primary, secondary, skin)

	portrait.Update(0.5)
	if portrait.currentFrame == 0 {
		t.Error("frame should have advanced after update")
	}

	portrait.Reset()
	if portrait.currentFrame != 0 {
		t.Errorf("frame should be 0 after reset, got %d", portrait.currentFrame)
	}
	if portrait.elapsed != 0 {
		t.Errorf("elapsed should be 0 after reset, got %f", portrait.elapsed)
	}
}

func TestPortraitCurrentFrameIndex(t *testing.T) {
	pg := NewPortraitGenerator(12345, 32)
	primary := color.RGBA{50, 50, 150, 255}
	secondary := color.RGBA{100, 80, 60, 255}
	skin := color.RGBA{200, 160, 140, 255}

	portrait := pg.GenerateAnimatedPortrait(primary, secondary, skin)

	idx := portrait.CurrentFrameIndex()
	if idx != 0 {
		t.Errorf("initial frame index should be 0, got %d", idx)
	}

	// Advance to next frame
	portrait.Update(0.3)
	idx = portrait.CurrentFrameIndex()
	if idx != 1 {
		t.Errorf("frame index should be 1 after advancing, got %d", idx)
	}
}

func TestPortraitStateChangeResetsFrame(t *testing.T) {
	pg := NewPortraitGenerator(12345, 32)
	primary := color.RGBA{50, 50, 150, 255}
	secondary := color.RGBA{100, 80, 60, 255}
	skin := color.RGBA{200, 160, 140, 255}

	portrait := pg.GenerateAnimatedPortrait(primary, secondary, skin)

	// Advance a few frames
	portrait.Update(0.5)
	if portrait.currentFrame == 0 {
		t.Error("frame should have advanced")
	}

	// Change state
	portrait.SetState(PortraitHurt)
	if portrait.currentFrame != 0 {
		t.Errorf("frame should reset to 0 on state change, got %d", portrait.currentFrame)
	}
	if portrait.elapsed != 0 {
		t.Errorf("elapsed should reset to 0 on state change, got %f", portrait.elapsed)
	}
}

func TestPortraitDeterminism(t *testing.T) {
	pg1 := NewPortraitGenerator(12345, 32)
	pg2 := NewPortraitGenerator(12345, 32)

	primary := color.RGBA{50, 50, 150, 255}
	secondary := color.RGBA{100, 80, 60, 255}
	skin := color.RGBA{200, 160, 140, 255}

	p1 := pg1.GenerateAnimatedPortrait(primary, secondary, skin)
	p2 := pg2.GenerateAnimatedPortrait(primary, secondary, skin)

	if p1.IdleFrameCount != p2.IdleFrameCount {
		t.Errorf("idle frame counts should match: %d vs %d", p1.IdleFrameCount, p2.IdleFrameCount)
	}
	if p1.FrameTime != p2.FrameTime {
		t.Errorf("frame times should match: %f vs %f", p1.FrameTime, p2.FrameTime)
	}
}
