package data

import "time"

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime     `json:"runtime,omitempty"`  // Runtime类型是自定义类型，实现了json.Unmarshaler接口
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}
