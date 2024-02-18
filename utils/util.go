package utils

import (
	"math"
	"os"
	"strconv"
)

// GetEnvDefault function to get the os env config or return the given default value if empty
func GetEnvDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GetLimitOffset sanitize and return default pagination limit and offset
func GetLimitOffset(limit string, defaultLimit int, page string, defaultPage int) (int, int, int) {
	// Check the page if it's a valid integer otherwise set it to default page
	curr_page, err := strconv.Atoi(page)
	if err != nil {
		curr_page = defaultPage
	}

	// Check the limit if it's a valid integer otherwise set it to default limit
	count, err := strconv.Atoi(limit)
	if err != nil {
		count = defaultLimit
	}

	// set minimal page will be 1
	curr_page = int(math.Max(1, float64(curr_page)))

	// set maximum page size is 500
	count = int(math.Min(500, float64(count)))

	// Calculate the offset of the data for the query
	offset := (curr_page - 1) * count

	return offset, count, curr_page
}
