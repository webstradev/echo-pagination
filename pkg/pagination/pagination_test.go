package pagination_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/matryer/is"

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
			is := is.New(t)

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

				gotPage, err := pagination.GetPage(c, pagination.WithPageText(pageText))
				is.NoErr(err)

				gotSize, err := pagination.GetPageSize(c, pagination.WithSizeText(sizeText))
				is.NoErr(err)

				is.Equal(gotPage, tt.expectedPage) // Didn't get the expected page.

				is.Equal(gotSize, tt.expectedSize) // Didn't get the expected size.

				return nil
			})

			// Execute the middleware and handler
			if err := handler(c); err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestGetPage(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		name          string
		setupContext  func(c echo.Context)
		customOptions []pagination.CustomOption
		expectedPage  int
		expectError   bool
	}{
		{
			name: "Default page text",
			setupContext: func(c echo.Context) {
				c.Set("page", 5)
			},
			expectedPage: 5,
			expectError:  false,
		},
		{
			name: "Custom page text",
			setupContext: func(c echo.Context) {
				c.Set("custom_page", 10)
			},
			customOptions: []pagination.CustomOption{pagination.WithPageText("custom_page")},
			expectedPage:  10,
			expectError:   false,
		},
		{
			name:         "Page not found in context",
			setupContext: func(c echo.Context) {},
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			tt.setupContext(c)

			page, err := pagination.GetPage(c, tt.customOptions...)

			if tt.expectError {
				is.True(err != nil)
			} else {
				is.NoErr(err)
				is.Equal(page, tt.expectedPage)
			}
		})
	}
}

func TestGetPageSize(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		name          string
		setupContext  func(c echo.Context)
		customOptions []pagination.CustomOption
		expectedSize  int
		expectError   bool
	}{
		{
			name: "Default size text",
			setupContext: func(c echo.Context) {
				c.Set("size", 25)
			},
			expectedSize: 25,
			expectError:  false,
		},
		{
			name: "Custom size text",
			setupContext: func(c echo.Context) {
				c.Set("custom_size", 50)
			},
			customOptions: []pagination.CustomOption{pagination.WithSizeText("custom_size")},
			expectedSize:  50,
			expectError:   false,
		},
		{
			name:         "Size not found in context",
			setupContext: func(c echo.Context) {},
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			tt.setupContext(c)

			size, err := pagination.GetPageSize(c, tt.customOptions...)

			if tt.expectError {
				is.True(err != nil)
			} else {
				is.NoErr(err)
				is.Equal(size, tt.expectedSize)
			}
		})
	}
}

func TestPaginationHeaders(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		name            string
		middleware      echo.MiddlewareFunc
		queryParams     url.Values
		expectedHeaders map[string]string
	}{
		{
			name:       "Default headers are set correctly",
			middleware: pagination.New(),
			queryParams: url.Values{
				"page": {"2"},
				"size": {"20"},
			},
			expectedHeaders: map[string]string{
				"X-Page": "2",
				"X-Size": "20",
			},
		},
		{
			name:       "Default headers without prefix",
			middleware: pagination.New(pagination.WithHeaderPrefix("")),
			queryParams: url.Values{
				"page": {"2"},
				"size": {"20"},
			},
			expectedHeaders: map[string]string{
				"page": "2",
				"size": "20",
			},
		},
		{
			name: "Custom text headers are set correctly",
			middleware: pagination.New(
				pagination.WithPageText("offset"),
				pagination.WithSizeText("limit"),
			),
			queryParams: url.Values{
				"offset": {"3"},
				"limit":  {"15"},
			},
			expectedHeaders: map[string]string{
				"X-Offset": "3",
				"X-limit":  "15",
			},
		},
		{
			name:            "No headers on invalid input",
			middleware:      pagination.New(),
			queryParams:     url.Values{"page": {"invalid"}},
			expectedHeaders: map[string]string{},
		},
		{
			name:        "Default values are set in headers when no query params",
			middleware:  pagination.New(pagination.WithHeaderPrefix("X-")),
			queryParams: url.Values{},
			expectedHeaders: map[string]string{
				"X-Page": "1",
				"X-Size": "10",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/?"+tt.queryParams.Encode(), nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Call middleware
			handler := tt.middleware(func(c echo.Context) error {
				return nil
			})

			// Execute the middleware and handler
			err := handler(c)
			is.NoErr(err) // Middleware execution should not error

			// Check for expected headers
			for headerKey, expectedValue := range tt.expectedHeaders {
				gotValue := rec.Header().Get(headerKey)
				is.Equal(gotValue, expectedValue) // Header value should match expected value
			}
		})
	}
}
