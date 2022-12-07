package log

import (
	"log"
)

var DefaultLogger = NewStdLogger(log.Writer())

type Logger interface {
	Log(level Level, keyvals ...interface{}) error
}

// type logger struct {
// 	logger    Logger
// 	prefix    []interface{}
// }

// func (c *logger) Log(level Level, keyvals ...interface{}) error {
// 	if err := c.logger.Log(level,keyvals); err != nil {
// 		return err
// 	}
// 	return nil
// }
