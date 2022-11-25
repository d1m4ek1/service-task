package seterror

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

// MessSomethingWentWrong Text errors message
const (
	MessSomethingWentWrong = "Что-то пошло не так!"
)

type SetTmplTime interface {
	Date() string
	Time(isTimeAmPm bool) string
}

type setTmplTime struct {
	NowTime time.Time
}

func (d *setTmplTime) Date() string {
	var date string = fmt.Sprintf("%d %s %d", d.NowTime.Day(), d.NowTime.Month().String(), d.NowTime.Year())
	return date
}

func (d *setTmplTime) Time(isTimeAmPm bool) string {
	var hour string = fmt.Sprint(d.NowTime.Hour())
	var minute string = fmt.Sprint(d.NowTime.Minute())
	var second string = fmt.Sprint(d.NowTime.Second())

	if d.NowTime.Minute() != 0 && d.NowTime.Minute() < 10 {
		minute = fmt.Sprintf("0%d", d.NowTime.Minute())
	}

	if d.NowTime.Second() != 0 && d.NowTime.Second() < 10 {
		second = fmt.Sprintf("0%d", d.NowTime.Second())
	}

	var hourMinute string = fmt.Sprintf("%s:%s", hour, minute)

	if isTimeAmPm {
		var h = d.NowTime.Hour()

		if h == 0 {
			hour = fmt.Sprint(12)
			hourMinute = fmt.Sprintf("%s:%s AM", hour, minute)
		} else if h == 12 {
			hourMinute = fmt.Sprintf("%d:%s PM", h, minute)
		} else if h > 12 {
			hour = fmt.Sprint(h - 12)
			hourMinute = fmt.Sprintf("%s:%s PM", hour, minute)
		} else {
			hourMinute = fmt.Sprintf("%d:%s AM", h, minute)
		}
	}

	var time string = fmt.Sprintf("%s Second: %s", hourMinute, second)
	return time
}

type AppErrorLogger interface {
	LogError() error
	Wrap() string
	NewError()
}

type AppError struct {
	CustomMessage, File, DirToLogError string
	NowTime                            time.Time
	ErrorMessage                       error
	IsLine                             int
	IsTimeAmPm                         bool
}

func (w *AppError) LogError() error {
	if w.DirToLogError == "" {
		return nil
	}

	if _, err := os.Stat(w.DirToLogError); err != nil {
		file, err := os.Create(w.DirToLogError)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	file, err := os.OpenFile(w.DirToLogError, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString(w.Wrap())

	return nil
}

func (w *AppError) Wrap() string {
	var setTime SetTmplTime = &setTmplTime{
		NowTime: w.NowTime,
	}

	return fmt.Sprintf("\n\nERROR:\n\tError: %s;\n\tMessage: %s;\n\tLine: %d;\n\tFile: %s;\n\tDate: %s;\n\tTime: %s;",
		w.ErrorMessage, w.CustomMessage, w.IsLine, w.File, setTime.Date(), setTime.Time(w.IsTimeAmPm))
}

func SetAppError(customMessage string, errorMessage error) {
	_, file, isLine, _ := runtime.Caller(1)

	var wrap = &AppError{
		DirToLogError: "./logs/errors.log",
		NowTime:       time.Now(),
		IsTimeAmPm:    true,
		CustomMessage: customMessage,
		IsLine:        isLine,
		File:          file,
		ErrorMessage:  errorMessage,
	}

	fmt.Println(wrap.Wrap())
	if err := wrap.LogError(); err != nil {
		log.Fatal(err)
		return
	}
}
