// Package pagination provides a middleware for the echo web framework to handle
// pagination. It allows for the usage of url parameters like `?page=1&size=25`
// to paginate data on your API. The values will be propagated throughout the
// request context.
package pagination

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// New returns a new pagination middleware with custom values.
func New(customOptions ...CustomOption) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {
			p := &paginator{
				opts: applyCustomOptionsToDefault(customOptions...),
				c:    c,
			}

			page, err := p.getPageFromRequest()
			if err != nil {
				return p.abortWithBadRequest(err)
			}

			if err := p.validatePage(page); err != nil {
				return p.abortWithBadRequest(err)
			}

			pageSize, err := p.getPageSizeFromRequest()
			if err != nil {
				return p.abortWithBadRequest(err)
			}

			if err := p.validatePageSize(pageSize); err != nil {
				return p.abortWithBadRequest(err)
			}

			p.setPageAndPageSize(page, pageSize)

			return next(c)
		}
	}
}

type paginator struct {
	opts options
	c    echo.Context
}

func (p *paginator) abortWithBadRequest(err error) error {
	return p.c.String(http.StatusBadRequest, err.Error())
}

func (p *paginator) getPageFromRequest() (int, error) {
	return p.getIntFromContextWithDefault(p.opts.PageText, p.opts.DefaultPage)
}

func (p *paginator) getPageSizeFromRequest() (int, error) {
	return p.getIntFromContextWithDefault(p.opts.SizeText, p.opts.DefaultPageSize)
}

func (p *paginator) getIntFromContextWithDefault(key string, defaultValue int) (int, error) {
	valueStr := p.c.QueryParam(key)
	if valueStr == "" {
		return defaultValue, nil
	}

	valueInt, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("%s parameter must be an integer", key)
	}

	return valueInt, nil
}

func (p *paginator) validatePage(page int) error {
	if page < 0 {
		return errors.New("page number must be positive")
	}

	return nil
}

func (p *paginator) validatePageSize(size int) error {
	if size < p.opts.MinPageSize || size > p.opts.MaxPageSize {
		return errors.New("page size must be between %d and %d")
	}

	return nil
}

func (p *paginator) setPageAndPageSize(page int, size int) {
	p.c.Set(p.opts.PageText, page)
	p.c.Set(p.opts.SizeText, size)
}
