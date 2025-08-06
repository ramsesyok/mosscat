package repository

import (
	"database/sql"
	"testing"
	"time"
)

func TestWhereClause(t *testing.T) {
	if whereClause(nil) != "" {
		t.Fatalf("expected empty string")
	}
	wh := []string{"a = 1", "b = 2"}
	got := whereClause(wh)
	if got != "WHERE a = 1 AND b = 2" {
		t.Fatalf("unexpected where: %s", got)
	}
}

func TestStrPtrAndTimePtr(t *testing.T) {
	ns := sql.NullString{String: "x", Valid: true}
	if *strPtr(ns) != "x" {
		t.Fatal("strPtr returned wrong value")
	}
	if strPtr(sql.NullString{}) != nil {
		t.Fatal("strPtr empty should return nil")
	}

	now := sql.NullTime{Time: time.Now(), Valid: true}
	if timePtr(now) == nil {
		t.Fatal("timePtr should not be nil")
	}
	if timePtr(sql.NullTime{}) != nil {
		t.Fatal("timePtr empty should be nil")
	}
}
