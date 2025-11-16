package roulette

import "testing"

type fixedRand struct {
	value int
}

func (f fixedRand) Intn(n int) int {
	return f.value % n
}

func TestRouletteNumberWin(t *testing.T) {
	s := newSessionWithRand(fixedRand{value: 17})
	if _, _, err := s.Handle("start", nil); err != nil {
		t.Fatalf("start failed: %v", err)
	}
	if _, _, err := s.Handle("wager", []string{"10", "17"}); err != nil {
		t.Fatalf("wager failed: %v", err)
	}
	messages, finished, err := s.Handle("spin", nil)
	if err != nil {
		t.Fatalf("spin failed: %v", err)
	}
	if finished {
		t.Fatalf("spin should not finish session")
	}
	if s.bankroll != 550 {
		t.Fatalf("expected bankroll 550, got %d", s.bankroll)
	}
	if len(messages) < 2 {
		t.Fatalf("expected spin messages")
	}
}

func TestRouletteColorLoss(t *testing.T) {
	s := newSessionWithRand(fixedRand{value: 2}) // black
	if _, _, err := s.Handle("start", nil); err != nil {
		t.Fatalf("start failed: %v", err)
	}
	if _, _, err := s.Handle("wager", []string{"20", "red"}); err != nil {
		t.Fatalf("wager failed: %v", err)
	}
	_, _, err := s.Handle("spin", nil)
	if err != nil {
		t.Fatalf("spin failed: %v", err)
	}
	if s.bankroll != 180 {
		t.Fatalf("expected bankroll 180 after loss, got %d", s.bankroll)
	}
}
