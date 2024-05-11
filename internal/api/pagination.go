package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

var ErrIncorrectLimit = errors.New("limit must be greater than 0")
var ErrIncorrectPageNo = errors.New("page number must be greater than 0")

var (
	defaultLimit = 30
	defaultPage  = 1
)

type PageInfo struct {
	Limit  int
	Offset int
	PageNo int
}

func (s *Server) getPageInfo(r *http.Request) (*PageInfo, error) {
	var limit int
	limitStr := r.FormValue("limit")
	switch limitStr {
	case "":
		limit = defaultLimit
	default:
		l, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, fmt.Errorf("failed to get limit: %w", err)
		}
		limit = l
	}
	if limit <= 0 {
		return nil, ErrIncorrectLimit
	}

	var pageNo int
	pageNoStr := r.FormValue("page")
	switch pageNoStr {
	case "":
		pageNo = defaultPage
	default:
		p, err := strconv.Atoi(pageNoStr)
		if err != nil {
			return nil, fmt.Errorf("failed to get page number: %w", err)
		}
		pageNo = p
	}
	if pageNo <= 0 {
		return nil, ErrIncorrectPageNo
	}

	offset := (pageNo - 1) * limit

	page := PageInfo{
		Limit:  limit,
		Offset: offset,
		PageNo: pageNo,
	}

	return &page, nil
}
