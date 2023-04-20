package models

import (
	"fmt"
	"strconv"

	"userstyles.world/modules/config"
)

type Pagination struct {
	Min   int
	Prev3 int
	Prev2 int
	Prev1 int
	Now   int
	Next1 int
	Next2 int
	Next3 int
	Max   int
	Path  string
	Sort  string
}

// NewPagination is a convenience function that initializes pagination struct.
func NewPagination(page int, sort, path string) Pagination {
	return Pagination{
		Now:  page,
		Path: path,
		Sort: sort,
	}
}

// URL generates a dynamic path from available items.
func (p Pagination) URL(page int) string {
	s := fmt.Sprintf("%s?page=%d", p.Path, page)
	if p.Sort != "" {
		s += fmt.Sprintf("&sort=%s", p.Sort)
	}
	return s
}

// IsValidPage checks whether a passed parameter is a valid number.
func IsValidPage(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return i, err
}

func (p *Pagination) CalcItems(total int) {
	if total == 0 {
		p.Max = 1
		p.Min = 1
		return
	}

	p.Min = 1

	// Calculate max page and remainder.
	p.Max = total / config.AppPageMaxItems
	if total%config.AppPageMaxItems > 0 {
		p.Max++
	}

	// Set prev/next.
	p.Prev1 = p.Now - 1
	p.Prev2 = p.Now - 2
	p.Prev3 = p.Now - 3
	p.Next1 = p.Now + 1
	p.Next2 = p.Now + 2
	p.Next3 = p.Now + 3
}

func (p *Pagination) OutOfBounds() bool {
	// Display last page if requested page is greater than max page.
	if p.Now > p.Max {
		p.Now = p.Max
		return true
	}

	// Display first page if requested page is less than 1.
	if p.Now < 1 {
		p.Now = 1
		return true
	}

	return false
}

func (p *Pagination) SortStyles() string {
	switch p.Sort {
	case "newest":
		return "styles.created_at DESC"
	case "oldest":
		return "styles.created_at ASC"
	case "recentlyupdated":
		return "styles.updated_at DESC"
	case "leastupdated":
		return "styles.updated_at ASC"
	case "mostinstalls":
		return "installs DESC"
	case "leastinstalls":
		return "installs ASC"
	case "mostviews":
		return "views DESC"
	case "leastviews":
		return "views ASC"
	default:
		return "styles.id ASC"
	}
}
