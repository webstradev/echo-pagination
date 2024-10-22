package pagination_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/webstradev/echo-pagination/pkg/pagination"
)

func TestPaginationMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		middleware     echo.MiddlewareFunc
		queryParams    url.Values
		expectedPage   int
		expectedSize   int
		customPageText string
		customSizeText string
	}{
		{
			"Non int Page Param - Bad Request",
			pagination.New(),
			url.Values{
				"page": {"notanumber"},
			},
			0,
			0,
			"",
			"",
		},
		{
			"Non int Size Param - Bad Request",
			pagination.New(),
			url.Values{
				"page": {"1"},
				"size": {"notanumber"},
			},
			0,
			0,
			"",
			"",
		},
		{
			"Negative Page Param - Bad Request",
			pagination.New(),
			url.Values{
				"page": {"-1"},
			},
			0,
			0,
			"",
			"",
		},
		{
			"Size below min - Bad Request",
			pagination.New(),
			url.Values{
				"page": {"1"},
				"size": {"0"},
			},
			0,
			0,
			"",
			"",
		},
		{
			"Size above max - Bad Request",
			pagination.New(),
			url.Values{
				"page": {"1"},
				"size": {"101"},
			},
			0,
			0,
			"",
			"",
		},
		{
			"Default Handling",
			pagination.New(),
			url.Values{},
			1,
			10,
			"",
			"",
		},
		{
			"The first 100 results",
			pagination.New(),
			url.Values{
				"page": {"1"},
				"size": {"100"},
			},
			1,
			100,
			"",
			"",
		},
		{
			"The second 20 results",
			pagination.New(),
			url.Values{
				"page": {"2"},
				"size": {"20"},
			},
			2,
			20,
			"",
			"",
		},
		{
			"Custom Handling",
			pagination.New(
				pagination.WithPageText("pages"),
				pagination.WithSizeText("items"),
				pagination.WithDefaultPage(0),
				pagination.WithDefaultPageSize(5),
				pagination.WithMinPageSize(1),
				pagination.WithMaxPageSize(25),
			),
			url.Values{},
			0,
			5,
			"pages",
			"items",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/?"+tt.queryParams.Encode(), nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Call middleware
			handler := tt.middleware(func(c echo.Context) error {
				// Handle custom page and size text
				pageText := "page"
				sizeText := "size"
				if tt.customPageText != "" {
					pageText = tt.customPageText
				}
				if tt.customSizeText != "" {
					sizeText = tt.customSizeText
				}

				gotPage := c.Get(pageText)
				gotSize := c.Get(sizeText)

				// Check if the page and pageSize are set correctly
				if gotPage != tt.expectedPage {
					t.Errorf("Expected page %d, got %v", tt.expectedPage, gotPage)
				}

				if gotSize != tt.expectedSize {
					t.Errorf("Expected size %d, got %v", tt.expectedSize, gotSize)
				}

				return nil
			})

			// Execute the middleware and handler
			if err := handler(c); err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
