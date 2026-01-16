package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"cmd/internal/calendar"
)

type Handler struct {
	Svc *calendar.Service
}

type resultResponse struct {
	Result string `json:"result"`
	ID     int64  `json:"id,omitempty"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid method"})
		return
	}

	if err := r.ParseForm(); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "failed to parse form"})
		return
	}

	userID, err := strconv.ParseInt(r.FormValue("user_id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid user_id"})
		return
	}

	day, err := calendar.ParseDay(r.FormValue("date"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	text := r.FormValue("event")

	id, err := h.Svc.Create(userID, day, text)
	if err != nil {
		// Ошибка ввода из сервиса 400
		if errors.Is(err, calendar.ErrInvalid) {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
			return
		}
		// Бизнес-логика 503
		if errors.Is(err, calendar.ErrConflict) {
			writeJSON(w, http.StatusServiceUnavailable, errorResponse{Error: err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal error"})
		return
	}

	writeJSON(w, http.StatusOK, resultResponse{Result: "created", ID: id})
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid metod"})
		return
	}
	if err := r.ParseForm(); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "failed to parse form"})
		return
	}

	id, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return
	}
	userID, err := strconv.ParseInt(r.FormValue("user_id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid user_id"})
		return
	}
	day, err := calendar.ParseDay(r.FormValue("date"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}
	text := r.FormValue("event")

	err = h.Svc.Update(id, userID, day, text)
	if err != nil {
		if errors.Is(err, calendar.ErrInvalid) {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, calendar.ErrNotFound) {
			writeJSON(w, http.StatusServiceUnavailable, errorResponse{Error: err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal error"})
		return
	}

	writeJSON(w, http.StatusOK, resultResponse{Result: "updated"})
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid method"})
		return
	}
	if err := r.ParseForm(); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "failed to parse form"})
		return
	}

	id, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return
	}
	userID, err := strconv.ParseInt(r.FormValue("user_id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid user_id"})
		return
	}

	err = h.Svc.Delete(id, userID)
	if err != nil {
		if errors.Is(err, calendar.ErrInvalid) {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, calendar.ErrNotFound) {
			writeJSON(w, http.StatusServiceUnavailable, errorResponse{Error: err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal error"})
		return
	}

	writeJSON(w, http.StatusOK, resultResponse{Result: "deleted"})
}

func (h *Handler) EventsForDay(w http.ResponseWriter, r *http.Request) {
	h.eventsForRange(w, r, "day")
}
func (h *Handler) EventsForWeek(w http.ResponseWriter, r *http.Request) {
	h.eventsForRange(w, r, "week")
}
func (h *Handler) EventsForMonth(w http.ResponseWriter, r *http.Request) {
	h.eventsForRange(w, r, "month")
}

func (h *Handler) eventsForRange(w http.ResponseWriter, r *http.Request, mode string) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid method"})
		return
	}

	userID, err := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid user_id"})
		return
	}
	day, err := calendar.ParseDay(r.URL.Query().Get("date"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	var events []calendar.Event
	switch mode {
	case "day":
		events, err = h.Svc.EventsForDay(userID, day)
	case "week":
		events, err = h.Svc.EventsForWeek(userID, day)
	case "month":
		events, err = h.Svc.EventsForMonth(userID, day)
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"result": events})
}
