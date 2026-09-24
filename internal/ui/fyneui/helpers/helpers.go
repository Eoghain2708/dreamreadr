package helpers

import "time"

func FormatCreatedAt(t time.Time) string {
	return t.Format("Jan 2 2006")
}
