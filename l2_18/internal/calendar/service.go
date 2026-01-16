package calendar

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrNotFound = errors.New("event not found") // бизнес-ошибка HTTP 503
	ErrConflict = errors.New("event conflict")  // опционально
	ErrInvalid  = errors.New("invalid input")   // если валидируете в сервисе
)

type Service struct {
	mu     sync.RWMutex
	nextID int64
	events map[int64]Event // ключ - ID
}

func NewService() *Service {
	return &Service{
		events: make(map[int64]Event),
	}
}

// ParseDay - парсер даты (HTTP 400)
func ParseDay(dateStr string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", dateStr) // UTC
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date %q (expected YYYY-MM-DD): %w", dateStr, err)
	}
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC), nil
}

func normalizeDay(t time.Time) time.Time {
	t = t.In(time.UTC)
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func (s *Service) Create(userID int64, day time.Time, text string) (int64, error) {
	if userID <= 0 {
		return 0, fmt.Errorf("user_id must be > 0: %w", ErrInvalid)
	}
	if text == "" {
		return 0, fmt.Errorf("event text is empty: %w", ErrInvalid)
	}

	day = normalizeDay(day)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextID++
	id := s.nextID

	s.events[id] = Event{
		ID:     id,
		UserID: userID,
		Date:   day,
		Text:   text,
	}
	return id, nil
}

func (s *Service) Update(id int64, userID int64, day time.Time, text string) error {
	if id <= 0 {
		return fmt.Errorf("id must be > 0: %w", ErrInvalid)
	}
	if userID <= 0 {
		return fmt.Errorf("user_id must be > 0: %w", ErrInvalid)
	}
	if text == "" {
		return fmt.Errorf("event text is empty: %w", ErrInvalid)
	}

	day = normalizeDay(day)

	s.mu.Lock()
	defer s.mu.Unlock()

	ev, ok := s.events[id]
	if !ok {
		return fmt.Errorf("update id=%d: %w", id, ErrNotFound)
	}

	//проверка принадлежности
	if ev.UserID != userID {
		return fmt.Errorf("update id=%d: %w", id, ErrNotFound)
	}

	ev.Date = day
	ev.Text = text
	s.events[id] = ev
	return nil
}

func (s *Service) Delete(id int64, userID int64) error {
	if id <= 0 {
		return fmt.Errorf("id must be > 0: %w", ErrInvalid)
	}
	if userID <= 0 {
		return fmt.Errorf("user_id must be > 0: %w", ErrInvalid)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	ev, ok := s.events[id]
	if !ok {
		return fmt.Errorf("delete id=%d: %w", id, ErrNotFound)
	}
	if ev.UserID != userID {
		return fmt.Errorf("delete id=%d: %w", id, ErrNotFound)
	}

	delete(s.events, id)
	return nil
}

func (s *Service) EventsForDay(userID int64, day time.Time) ([]Event, error) {
	start := normalizeDay(day)
	end := start.AddDate(0, 0, 1)
	return s.listInRange(userID, start, end), nil
}

func (s *Service) EventsForWeek(userID int64, day time.Time) ([]Event, error) {
	start := normalizeDay(day)
	end := start.AddDate(0, 0, 7)
	return s.listInRange(userID, start, end), nil
}

func (s *Service) EventsForMonth(userID int64, day time.Time) ([]Event, error) {
	day = normalizeDay(day)
	y, m, _ := day.Date()
	start := time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)
	return s.listInRange(userID, start, end), nil
}

func (s *Service) listInRange(userID int64, start, end time.Time) []Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Event, 0)
	for _, ev := range s.events {
		if ev.UserID != userID {
			continue
		}
		if !ev.Date.Before(start) && ev.Date.Before(end) {
			out = append(out, ev) // value copy
		}
	}
	return out
}
