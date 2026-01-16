package calendar

import (
	"errors"
	"testing"
	"time"
)

func mustDay(t *testing.T, s string) time.Time {
	t.Helper()
	d, err := ParseDay(s)
	if err != nil {
		t.Fatalf("ParseDay(%q) error: %v", s, err)
	}
	return d
}

func TestParseDay_OK(t *testing.T) {
	d, err := ParseDay("2026-01-17")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if d.Location() != time.UTC {
		t.Fatalf("expected UTC, got %v", d.Location())
	}
	if d.Hour() != 0 || d.Minute() != 0 || d.Second() != 0 || d.Nanosecond() != 0 {
		t.Fatalf("expected normalized time at 00:00:00, got %v", d)
	}
}

func TestParseDay_Invalid(t *testing.T) {
	_, err := ParseDay("17-01-2026")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestCreate_OK_ReturnsIDAndPersists(t *testing.T) {
	svc := NewService()

	day := mustDay(t, "2026-01-17")
	id, err := svc.Create(1, day, "Buy milk")
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if id <= 0 {
		t.Fatalf("expected id > 0, got %d", id)
	}

	// Провери что событие действительно хранится и выбирается
	events, err := svc.EventsForDay(1, day)
	if err != nil {
		t.Fatalf("EventsForDay error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	ev := events[0]
	if ev.ID != id {
		t.Fatalf("expected ID=%d, got %d", id, ev.ID)
	}
	if ev.UserID != 1 {
		t.Fatalf("expected UserID=1, got %d", ev.UserID)
	}
	if !ev.Date.Equal(day) {
		t.Fatalf("expected Date=%v, got %v", day, ev.Date)
	}
	if ev.Text != "Buy milk" {
		t.Fatalf("expected Text=%q, got %q", "Buy milk", ev.Text)
	}
}

func TestCreate_InvalidInput(t *testing.T) {
	svc := NewService()
	day := mustDay(t, "2026-01-17")

	t.Run("invalid user_id", func(t *testing.T) {
		_, err := svc.Create(0, day, "x")
		if err == nil || !errors.Is(err, ErrInvalid) {
			t.Fatalf("expected ErrInvalid, got %v", err)
		}
	})

	t.Run("empty text", func(t *testing.T) {
		_, err := svc.Create(1, day, "")
		if err == nil || !errors.Is(err, ErrInvalid) {
			t.Fatalf("expected ErrInvalid, got %v", err)
		}
	})
}

func TestUpdate_OK(t *testing.T) {
	svc := NewService()

	day1 := mustDay(t, "2026-01-17")
	day2 := mustDay(t, "2026-01-18")

	id, err := svc.Create(1, day1, "A")
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}

	if err := svc.Update(id, 1, day2, "B"); err != nil {
		t.Fatalf("Update error: %v", err)
	}

	events1, err := svc.EventsForDay(1, day1)
	if err != nil {
		t.Fatalf("EventsForDay error: %v", err)
	}
	if len(events1) != 0 {
		t.Fatalf("expected 0 events on day1, got %d", len(events1))
	}

	events2, err := svc.EventsForDay(1, day2)
	if err != nil {
		t.Fatalf("EventsForDay error: %v", err)
	}
	if len(events2) != 1 {
		t.Fatalf("expected 1 event on day2, got %d", len(events2))
	}
	if events2[0].ID != id || events2[0].Text != "B" {
		t.Fatalf("unexpected event after update: %+v", events2[0])
	}
}

func TestUpdate_NotFound(t *testing.T) {
	svc := NewService()
	day := mustDay(t, "2026-01-17")

	err := svc.Update(999, 1, day, "X")
	if err == nil || !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDelete_OK(t *testing.T) {
	svc := NewService()

	day := mustDay(t, "2026-01-17")
	id, err := svc.Create(1, day, "X")
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}

	if err := svc.Delete(id, 1); err != nil {
		t.Fatalf("Delete error: %v", err)
	}

	events, err := svc.EventsForDay(1, day)
	if err != nil {
		t.Fatalf("EventsForDay error: %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("expected 0 events after delete, got %d", len(events))
	}
}

func TestDelete_NotFound(t *testing.T) {
	svc := NewService()

	err := svc.Delete(999, 1)
	if err == nil || !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestEventsForWeek_Includes7DaysFromStart(t *testing.T) {
	svc := NewService()
	start := mustDay(t, "2026-01-17")

	// Внутри диапазона [start, start+7days)
	in1 := start
	in2 := start.AddDate(0, 0, 6)

	// Вне диапазона
	out := start.AddDate(0, 0, 7)

	id1, _ := svc.Create(1, in1, "in1")
	id2, _ := svc.Create(1, in2, "in2")
	_, _ = svc.Create(1, out, "out")

	events, err := svc.EventsForWeek(1, start)
	if err != nil {
		t.Fatalf("EventsForWeek error: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}

	// Проверим, что вернулись нужные id (порядок не гарантирован)
	seen := map[int64]bool{}
	for _, e := range events {
		seen[e.ID] = true
	}
	if !seen[id1] || !seen[id2] {
		t.Fatalf("expected ids %d and %d, got %+v", id1, id2, seen)
	}
}

func TestEventsForMonth_CalendarMonth(t *testing.T) {
	svc := NewService()

	anchor := mustDay(t, "2026-01-17")

	in1 := mustDay(t, "2026-01-01")
	in2 := mustDay(t, "2026-01-31")
	out := mustDay(t, "2026-02-01")

	id1, _ := svc.Create(1, in1, "jan1")
	id2, _ := svc.Create(1, in2, "jan31")
	_, _ = svc.Create(1, out, "feb1")

	events, err := svc.EventsForMonth(1, anchor)
	if err != nil {
		t.Fatalf("EventsForMonth error: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events in January, got %d", len(events))
	}

	seen := map[int64]bool{}
	for _, e := range events {
		seen[e.ID] = true
	}
	if !seen[id1] || !seen[id2] {
		t.Fatalf("expected ids %d and %d, got %+v", id1, id2, seen)
	}
}

func TestUserIsolation(t *testing.T) {
	svc := NewService()
	day := mustDay(t, "2026-01-17")

	_, _ = svc.Create(1, day, "u1")
	_, _ = svc.Create(2, day, "u2")

	ev1, _ := svc.EventsForDay(1, day)
	ev2, _ := svc.EventsForDay(2, day)

	if len(ev1) != 1 || ev1[0].UserID != 1 {
		t.Fatalf("expected 1 event for user 1, got %+v", ev1)
	}
	if len(ev2) != 1 || ev2[0].UserID != 2 {
		t.Fatalf("expected 1 event for user 2, got %+v", ev2)
	}
}
