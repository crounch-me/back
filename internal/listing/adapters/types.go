package adapters

import "time"

type List struct {
	ID              string
	Name            string
	CreationDate    time.Time
	ArchivationDate *time.Time
}
