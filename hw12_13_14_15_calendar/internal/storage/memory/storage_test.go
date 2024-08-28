package memorystorage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/storage"
)

func TestStorage(t *testing.T) {
	// Create test events.
	userID := 1
	date1, _ := time.Parse(time.RFC3339, "2024-09-02T10:00:00+03:00")
	date2, _ := time.Parse(time.RFC3339, "2024-09-03T10:00:00+03:00")
	date3, _ := time.Parse(time.RFC3339, "2024-09-10T19:00:00+03:00")
	date4, _ := time.Parse(time.RFC3339, "2024-09-15T12:00:00+03:00")
	events := []*storage.Event{
		{
			Title:       "Event 1",
			StartDate:   date1,
			EndDate:     date1.Add(time.Minute * 30),
			Description: "Test event 1",
			UserID:      userID,
		},
		{
			Title:       "Event 2",
			StartDate:   date2,
			EndDate:     date2.Add(time.Minute * 15),
			Description: "Test event 2",
			UserID:      userID,
		},
		{
			Title:       "Event 3",
			StartDate:   date3,
			EndDate:     date3.Add(time.Hour * 2),
			Description: "Test event 3",
			UserID:      userID,
		},
		{
			Title:       "Event 4",
			StartDate:   date4,
			EndDate:     date4.Add(time.Hour),
			Description: "Test event 4",
			UserID:      userID,
		},
	}

	date5, _ := time.Parse(time.RFC3339, "2024-10-23T17:50:00+03:00")
	update := storage.Event{
		ID:        4,
		Title:     "Updated event 4",
		StartDate: date5,
		EndDate:   date5.Add(time.Minute * 45),
		UserID:    userID,
	}

	ErrNotFound := storage.ErrEventNotFound

	// Create new storage.
	storage := New()

	// Test storage.
	t.Run("test storage", func(t *testing.T) {
		// Create events.
		for _, event := range events {
			err := storage.Create(*event)
			require.Nil(t, err)
		}

		// List events for day.
		dayEvents, err := storage.ListDay(date1, userID)
		require.Nil(t, err)
		require.Equal(t, events[0].Title, dayEvents[0].Title)

		// List events for day.
		weekEvents, err := storage.ListWeek(date1, userID)
		require.Nil(t, err)
		for i, event := range weekEvents {
			require.Equal(t, events[i].Title, event.Title)
		}

		// List events for month.
		monthEvents, err := storage.ListMonth(date1, userID)
		require.Nil(t, err)
		for i, event := range monthEvents {
			require.Equal(t, events[i].Title, event.Title)
		}

		// Update event.
		id := 4
		err = storage.Update(id, update)
		require.Nil(t, err)
		require.Equal(t, update, storage.events[4])

		// Delete event.
		id = 1
		err = storage.Delete(id)
		require.Nil(t, err)
		err = storage.Delete(id)
		require.Equal(t, ErrNotFound, err)
	})
}
