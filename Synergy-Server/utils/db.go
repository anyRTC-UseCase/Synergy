package utils

import (
	"database/sql"
	"errors"
)

func RowToMap(rows *sql.Rows, tableData *[]map[string]interface{}) (err error) {
	columns, err := rows.Columns()
	if err != nil {
		return errors.New("get column error")
	}
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		err = rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		*tableData = append(*tableData, entry)
	}
	return nil
}
