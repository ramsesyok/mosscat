package dbtime

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type DBTime struct{ time.Time }

var timeLayouts = []string{
	time.RFC3339Nano,
	"2006-01-02 15:04:05.999999999Z",
	"2006-01-02 15:04:05Z07:00",
	"2006-01-02 15:04:05",
}

func (t *DBTime) Scan(v any) error {
	switch x := v.(type) {
	case time.Time:
		t.Time = x.UTC()
		return nil
	case string:
		for _, layout := range timeLayouts {
			if ts, err := time.Parse(layout, x); err == nil {
				t.Time = ts.UTC()
				return nil
			}
		}
		return fmt.Errorf("unsupported time string %q", x)
	case []byte:
		return t.Scan(string(x))
	case int64:
		t.Time = time.Unix(x, 0).UTC()
		return nil
	default:
		return fmt.Errorf("unsupported type %T", v)
	}
}

func (t DBTime) Value() (driver.Value, error) {
	return t.Time.UTC(), nil
}

// Time returns time.Time value in UTC.
func (t DBTime) TimeValue() time.Time {
	return t.Time.UTC()
}
