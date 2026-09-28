package nav

import (
	"testing"
	"time"

	"github.com/Falasefemi2/worldcupdashboard/data"
)

func TestNav(t *testing.T) {
	t.Run("format live match date", func(t *testing.T) {
		m := data.Match{
			Status: data.StatusLive,
			Minute: "90",
			Date:   time.Unix(1668771420, 0),
		}

		want := "LIVE 90"
		got := renderDatetime(m)
		if want != got {
			t.Fatalf("want %s, got %s", want, got)
		}
	})
}
