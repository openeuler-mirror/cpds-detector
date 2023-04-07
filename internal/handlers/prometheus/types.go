package prometheus

type queryParams struct {
	expression string
	time       int64
}

type rangeQueryParams struct {
	expression string
	startTime  int64
	endTime    int64
	step       int
}
