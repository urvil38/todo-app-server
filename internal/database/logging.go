package database

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"
	"unicode"
)

// QueryLoggingDisabled stops logging of queries when true.
// For use in tests only: not concurrency-safe.
var QueryLoggingDisabled bool

var queryCounter int64 // atomic: per-process counter for unique query IDs

type queryEndLogEntry struct {
	ID              string
	Query           string
	Args            string
	DurationSeconds float64
	Error           string `json:",omitempty"`
}

func logQuery(ctx context.Context, query string, args []interface{}, instanceID string, retryable bool) func(*error) {
	if QueryLoggingDisabled {
		return func(*error) {}
	}
	const maxlen = 300 // maximum length of displayed query

	// To make the query more compact and readable, replace newlines with spaces
	// and collapse adjacent whitespace.
	var r []rune
	for _, c := range query {
		if c == '\n' {
			c = ' '
		}
		if len(r) == 0 || !unicode.IsSpace(r[len(r)-1]) || !unicode.IsSpace(c) {
			r = append(r, c)
		}
	}
	query = string(r)
	if len(query) > maxlen {
		query = query[:maxlen] + "..."
	}

	uid := generateLoggingID(instanceID)

	// Construct a short string of the args.
	const (
		maxArgs   = 20
		maxArgLen = 50
	)
	var argStrings []string
	for i := 0; i < len(args) && i < maxArgs; i++ {
		s := fmt.Sprint(args[i])
		if len(s) > maxArgLen {
			s = s[:maxArgLen] + "..."
		}
		argStrings = append(argStrings, s)
	}
	if len(args) > maxArgs {
		argStrings = append(argStrings, "...")
	}
	argString := strings.Join(argStrings, ", ")

	log.Debugf("%s %s args=%s", uid, query, argString)
	start := time.Now()
	return func(errp *error) {
		dur := time.Since(start)
		if errp == nil { // happens with queryRow
			log.Debugf("%s done", uid)
		} else {
			entry := queryEndLogEntry{
				ID:              uid,
				Query:           query,
				Args:            argString,
				DurationSeconds: dur.Seconds(),
			}
			if *errp == nil {
				log.Debugf("%+v", entry)
			} else {
				entry.Error = (*errp).Error()
				logf := log.Errorf
				if errors.Is(ctx.Err(), context.Canceled) ||
					strings.Contains(entry.Error, "pq: canceling statement due to user request") {
					logf = log.Debugf
				}
				if retryable && isSerializationFailure(*errp) {
					logf = log.Debugf
				}
				logf("%+v", entry)
			}
		}
	}
}

func (db *DB) logTransaction(ctx context.Context) func(*error) {
	if QueryLoggingDisabled {
		return func(*error) {}
	}
	uid := generateLoggingID(db.instanceID)
	log.Debugf("%s transaction (isolation %s) started", uid, db.opts.Isolation)
	start := time.Now()
	return func(errp *error) {
		log.Debugf("%s transaction (isolation %s) finished in %s with error %v",
			uid, db.opts.Isolation, time.Since(start), *errp)
	}
}

func generateLoggingID(instanceID string) string {
	if instanceID == "" {
		instanceID = "local"
	}
	n := atomic.AddInt64(&queryCounter, 1)
	return fmt.Sprintf("%s-%d", instanceID, n)
}
