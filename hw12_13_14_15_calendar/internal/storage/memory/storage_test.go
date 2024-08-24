package memorystorage

import (
	"context"
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
			Title:  "Event 1",
			Date:   date1,
			UserID: userID,
		},
		{
			Title:  "Event 2",
			Date:   date2,
			UserID: userID,
		},
		{
			Title:  "Event 3",
			Date:   date3,
			UserID: userID,
		},
		{
			Title:  "Event 4",
			Date:   date4,
			UserID: userID,
		},
	}

	date5, _ := time.Parse(time.RFC3339, "2024-10-23T17:50:00+03:00")
	update := storage.Event{
		ID:     4,
		Title:  "Updated event 4",
		Date:   date5,
		UserID: userID,
	}

	notFound := storage.ErrEventNotFound

	// Create new storage.
	storage := New()

	ctx := context.Background()

	// Test storage.
	t.Run("test storage", func(t *testing.T) {
		// Create events.
		for _, event := range events {
			_, err := storage.Create(ctx, *event)
			require.Nil(t, err)
		}

		// List events for day.
		dayEvents, err := storage.ListDay(ctx, date1)
		require.Nil(t, err)
		require.Equal(t, events[0].Title, dayEvents[0].Title)

		// List events for day.
		weekEvents, err := storage.ListWeek(ctx, date1)
		require.Nil(t, err)
		for i, event := range weekEvents {
			require.Equal(t, events[i].Title, event.Title)
		}

		// List events for month.
		monthEvents, err := storage.ListMonth(ctx, date1)
		require.Nil(t, err)
		for i, event := range monthEvents {
			require.Equal(t, events[i].Title, event.Title)
		}

		// Update event.
		id := 4
		err = storage.Update(ctx, id, update)
		require.Nil(t, err)
		require.Equal(t, update, storage.events[4])

		// Delete event.
		id = 1
		err = storage.Delete(ctx, id)
		require.Nil(t, err)
		err = storage.Delete(ctx, id)
		require.Equal(t, notFound, err)
	})

}
