package internal

import (
	"github.com/zhanglp92/plugins/imports/internal/comments"
	"github.com/zhanglp92/plugins/imports/internal/imports"
)

// Process ...
func Process(data []byte, updateComment bool) ([]byte, error) {
	var handlers []func([]byte) ([]byte, error)

	if updateComment {
		handlers = append(handlers, comments.Process)
	}
	handlers = append(handlers, imports.Process)

	var err error
	for _, h := range handlers {
		if data, err = h(data); err != nil {
			return nil, err
		}
	}
	return data, nil
}
