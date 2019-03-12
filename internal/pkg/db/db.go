package db

import (
    "database/sql"
    "fmt"
)

var dbObj *sql.DB

func Open() (err error) {
    if dbObj != nil {
        return fmt.Errorf("db already initialized")
    }
    dbObj, err = sql.Open("mysql", username+":"+password+"@/")
    return
}

func Close() error {
    return dbObj.Close()
}

func Query(query string, args ...interface{}) (*sql.Rows, error) {
    return dbObj.Query(query, args...)
}

func Exec(query string, args ...interface{}) (sql.Result, error) {
    return dbObj.Exec(query, args...)
}

func isExists(dbName string, tableName string, where string, args ...interface{}) (id int64, err error) {
    row, err := findRowBy(dbName, tableName, "id", where, args...)
    if err != nil {
        return
    }
    err = row.Scan(&id)
    if err != nil {
        if err == sql.ErrNoRows {
            return 0, nil
        }
        return
    }
    return id, nil
}

func insert(dbName string, tableName string, cols string, values string, args ...interface{}) (int64, error) {
    if dbObj == nil {
        return 0, fmt.Errorf("db wasn't initialized")
    }

    result, err := dbObj.Exec("INSERT INTO "+dbName+"."+tableName+" ("+cols+") VALUES ("+values+")", args...)
    if err != nil {
        return 0, err
    }

    return result.LastInsertId()
}

func findRowBy(dbName string, tableName string, cols string, where string, args ...interface{}) (*sql.Row, error) {
    if dbObj == nil {
        return nil, fmt.Errorf("db wasn't initialized")
    }

    if where == "" {
        where = "1"
    }
    return dbObj.QueryRow("SELECT "+cols+" FROM "+dbName+"."+tableName+" WHERE "+where, args...), nil
}

func findRowsBy(dbName string, tableName string, cols string, where string, args ...interface{}) (*sql.Rows, error) {
    if dbObj == nil {
        return nil, fmt.Errorf("db wasn't initialized")
    }

    if where == "" {
        where = "1"
    }
    return dbObj.Query("SELECT "+cols+" FROM "+dbName+"."+tableName+" WHERE "+where, args...)
}

func joinRowBy(dbName string, tableName string, cols string, where string, args ...interface{}) (*sql.Row, error) {
    if dbObj == nil {
        return nil, fmt.Errorf("db wasn't initialized")
    }

    if where == "" {
        where = "1"
    }
    return dbObj.QueryRow("SELECT "+cols+" FROM "+dbName+"."+tableName+" WHERE "+where, args...), nil
}

func joinRowsBy(dbName string, tableName string, cols string, where string, args ...interface{}) (*sql.Rows, error) {
    if dbObj == nil {
        return nil, fmt.Errorf("db wasn't initialized")
    }

    if where == "" {
        where = "1"
    }
    return dbObj.Query("SELECT "+cols+" FROM "+dbName+"."+tableName+" WHERE "+where, args...)
}

func updateBy(dbName string, tableName string, set string, where string, args ...interface{}) (int64, error) {
    if dbObj == nil {
        return 0, fmt.Errorf("db wasn't initialized")
    }

    if where == "" {
        where = "1"
    }
    result, err := dbObj.Exec("UPDATE "+dbName+"."+tableName+" SET "+set+" WHERE "+where, args...)
    if err != nil {
        return 0, err
    }

    return result.RowsAffected()
}

func removeBy(dbName string, tableName string, where string, args ...interface{}) (int64, error) {
    if dbObj == nil {
        return 0, fmt.Errorf("db wasn't initialized")
    }

    if where == "" {
        where = "1"
    }
    result, err := dbObj.Exec("DELETE FROM "+dbName+"."+tableName+" WHERE "+where, args...)
    if err != nil {
        return 0, err
    }

    return result.RowsAffected()
}
