package bundles

type PaginationBundle struct {
	Page    int64
	Limit   int64
	MaxFlag string `validate:"objectid"`
	MinFlag string `validate:"objectid"`
}

type PaginationResponseBundle struct {
	StartFlag string `validate:"objectid"`
	EndFlag   string `validate:"objectid"`
}
