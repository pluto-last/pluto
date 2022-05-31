package xgorm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
)

type SQLFieldsGetter interface {
	SQLFields() (table string, columnNames []string, values []interface{})
}

func BatchInsert(tx *gorm.DB, step int, a []SQLFieldsGetter) error {
	if len(a) == 0 {
		return nil
	}
	t := a
	for {
		var isEnd bool
		lt := len(t)
		if step >= lt {
			step = lt
			isEnd = true
		}
		b := t[:step]
		t = t[step:]
		s, v := GenInsertSql(b)
		err := ExecSQL(tx, s, v)
		if err != nil {
			return err
		}
		if isEnd {
			break
		}
	}
	return nil
}

func ExecSQL(tx *gorm.DB, s string, v []interface{}) (err error) {
	if len(v) == 0 {
		return nil
	}
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%#v", e))
		}
	}()
	if len(v) == 0 {
		return
	}
	err = tx.Exec(s, v...).Error
	return
}

func GenInsertSql(rows []SQLFieldsGetter) (string, []interface{}) {
	table, s, _ := rows[0].SQLFields()
	sql := fmt.Sprintf("INSERT INTO %s ( %s ) VALUES ",
		table,
		strings.Join(s, ","),
	)
	var blanks []string
	var values []interface{}
	for _, row := range rows {
		_, _, v := row.SQLFields()
		blanks = append(blanks, "("+
			strings.TrimRight(strings.Repeat("?,", len(s)), ",")+
			")")
		values = append(values, v...)
	}
	return sql + strings.Join(blanks, ","), values
}
