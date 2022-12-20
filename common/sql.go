package common

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

//
// Transaction
// @Description: 事务
func Transaction(db *sqlx.DB, task func(tx *sqlx.Tx) error) error {
	tx, err := db.Beginx() // 开启事务
	if err != nil {
		fmt.Printf("begin trans failed, err:%v\n", err)
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}()

	return task(tx)
}
