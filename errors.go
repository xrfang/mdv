package main

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

func assert(e interface{}, ntfy ...interface{}) {
	switch e.(type) {
	case nil:
	case bool:
		if !e.(bool) {
			mesg := "assertion failed"
			if len(ntfy) > 0 {
				mesg = ntfy[0].(string)
				if len(ntfy) > 1 {
					mesg = fmt.Sprintf(mesg, ntfy[1:]...)
				}
			}
			panic(errors.New(mesg))
		}
	case error:
		panic(e)
	default:
		panic(fmt.Errorf("assert: expect error or bool, got %T", e))
	}
}

type exception []string

func (e exception) Error() string {
	return strings.Join(e, "\n")
}

func trace(msg string, args ...interface{}) error {
	ex := exception{fmt.Sprintf(msg, args...)}
	n := 1
	for {
		n++
		pc, file, line, ok := runtime.Caller(n)
		if !ok {
			break
		}
		f := runtime.FuncForPC(pc)
		name := f.Name()
		if strings.HasPrefix(name, "runtime.") {
			continue
		}
		fn := strings.Split(file, "/")
		file = strings.Join(fn[len(fn)-2:], "/")
		ex = append(ex, fmt.Sprintf("\t(%s:%d) %s", file, line, name))
	}
	return ex
}
