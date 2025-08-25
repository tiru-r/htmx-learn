package db

import (
	"testing"
)

func TestNewPaginationParams(t *testing.T) {
	tests := []struct {
		name           string
		page           int
		pageSize       int
		expectedPage   int
		expectedSize   int
		expectedOffset int
	}{
		{
			name:           "valid parameters",
			page:           2,
			pageSize:       20,
			expectedPage:   2,
			expectedSize:   20,
			expectedOffset: 20,
		},
		{
			name:           "page less than 1",
			page:           0,
			pageSize:       10,
			expectedPage:   1,
			expectedSize:   10,
			expectedOffset: 0,
		},
		{
			name:           "page size too small",
			page:           1,
			pageSize:       3,
			expectedPage:   1,
			expectedSize:   DefaultPageSize,
			expectedOffset: 0,
		},
		{
			name:           "page size too large",
			page:           1,
			pageSize:       150,
			expectedPage:   1,
			expectedSize:   MaxPageSize,
			expectedOffset: 0,
		},
		{
			name:           "negative page",
			page:           -5,
			pageSize:       10,
			expectedPage:   1,
			expectedSize:   10,
			expectedOffset: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := NewPaginationParams(tt.page, tt.pageSize)

			if params.Page != tt.expectedPage {
				t.Errorf("Page = %d, expected %d", params.Page, tt.expectedPage)
			}
			if params.PageSize != tt.expectedSize {
				t.Errorf("PageSize = %d, expected %d", params.PageSize, tt.expectedSize)
			}
			if params.Offset != tt.expectedOffset {
				t.Errorf("Offset = %d, expected %d", params.Offset, tt.expectedOffset)
			}
		})
	}
}

func TestNewPaginatedResult(t *testing.T) {
	data := []string{"item1", "item2", "item3"}
	params := PaginationParams{Page: 2, PageSize: 10, Offset: 10}
	total := 25

	result := NewPaginatedResult(data, params, total)

	if result.Page != 2 {
		t.Errorf("Page = %d, expected 2", result.Page)
	}
	if result.PageSize != 10 {
		t.Errorf("PageSize = %d, expected 10", result.PageSize)
	}
	if result.Total != 25 {
		t.Errorf("Total = %d, expected 25", result.Total)
	}
	if result.TotalPages != 3 {
		t.Errorf("TotalPages = %d, expected 3", result.TotalPages)
	}
	if !result.HasNext {
		t.Error("HasNext = false, expected true")
	}
	if !result.HasPrev {
		t.Error("HasPrev = false, expected true")
	}
	if len(result.Data) != 3 {
		t.Errorf("Data length = %d, expected 3", len(result.Data))
	}
}

func TestPaginatedResultEdgeCases(t *testing.T) {
	t.Run("first page", func(t *testing.T) {
		data := []string{"item1"}
		params := PaginationParams{Page: 1, PageSize: 10, Offset: 0}
		total := 25

		result := NewPaginatedResult(data, params, total)

		if result.HasPrev {
			t.Error("HasPrev = true, expected false for first page")
		}
		if !result.HasNext {
			t.Error("HasNext = false, expected true for first page with more data")
		}
	})

	t.Run("last page", func(t *testing.T) {
		data := []string{"item1"}
		params := PaginationParams{Page: 3, PageSize: 10, Offset: 20}
		total := 25

		result := NewPaginatedResult(data, params, total)

		if !result.HasPrev {
			t.Error("HasPrev = false, expected true for last page")
		}
		if result.HasNext {
			t.Error("HasNext = true, expected false for last page")
		}
	})

	t.Run("single page", func(t *testing.T) {
		data := []string{"item1", "item2"}
		params := PaginationParams{Page: 1, PageSize: 10, Offset: 0}
		total := 2

		result := NewPaginatedResult(data, params, total)

		if result.HasPrev {
			t.Error("HasPrev = true, expected false for single page")
		}
		if result.HasNext {
			t.Error("HasNext = true, expected false for single page")
		}
		if result.TotalPages != 1 {
			t.Errorf("TotalPages = %d, expected 1 for single page", result.TotalPages)
		}
	})

	t.Run("empty result", func(t *testing.T) {
		data := []string{}
		params := PaginationParams{Page: 1, PageSize: 10, Offset: 0}
		total := 0

		result := NewPaginatedResult(data, params, total)

		if result.HasPrev {
			t.Error("HasPrev = true, expected false for empty result")
		}
		if result.HasNext {
			t.Error("HasNext = true, expected false for empty result")
		}
		if result.TotalPages != 1 {
			t.Errorf("TotalPages = %d, expected 1 for empty result", result.TotalPages)
		}
	})
}