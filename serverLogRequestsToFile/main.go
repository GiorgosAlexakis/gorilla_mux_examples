package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type RequestLog struct {
	Method        string
	Uri           string
	Host          string
	Header        string
	RemoteAddr    string
	ContentLength string
	Time          time.Time
}

func (l *RequestLog) createLogFromRequest(r *http.Request) *RequestLog {
	logData := RequestLog{
		Method:        r.Method,
		Uri:           r.RequestURI,
		Host:          r.Host,
		Header:        fmt.Sprint(r.Header),
		RemoteAddr:    r.RemoteAddr,
		ContentLength: strconv.FormatInt(r.ContentLength, 10),
		Time:          time.Now().UTC(),
	}
	return &logData
}

func writeLogToFile(logData *RequestLog) error {
	id := uuid.New()
	fileName := fmt.Sprintf("/tmp/%s", id.String())
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	out, err := json.Marshal(logData)
	if err != nil {
		return err
	}
	_, err = f.WriteString(string(out))
	if err != nil {
		return err
	}
	return nil
}

func loggingToFileMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestLog RequestLog
		logData := requestLog.createLogFromRequest(r)
		err := writeLogToFile(logData)
		if err != nil {
			log.Println(err)
			fmt.Fprint(w, fmt.Sprint(http.StatusInternalServerError))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

type MiddlewareFunc func(http.Handler) http.Handler

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Home Page")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.Use(loggingToFileMiddleware)
	log.Fatal(http.ListenAndServe(":8000", r))
}
