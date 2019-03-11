package db

import (
    "database/sql"
    "fmt"
)

func Start() {
    if err := AuthInit(); err != nil {
        return
    }
}

func insert(dbObj *sql.DB, dbName string, tableName string, cols string, values string, args ...interface{}) (int64, error) {
    if dbObj == nil {
        return 0, fmt.Errorf("db wasn't initialized")
    }

    result, err := dbObj.Exec("INSERT INTO "+dbName+"."+tableName+" ("+cols+") VALUES ("+values+")", args...)
    if err != nil {
        return 0, err
    }

    return result.LastInsertId()
}

func findRowBy(dbObj *sql.DB, dbName string, tableName string, cols string, where string, args ...interface{}) (*sql.Row, error) {
    if dbObj == nil {
        return nil, fmt.Errorf("db wasn't initialized")
    }

    if where == "" {
        where = "1"
    }
    return dbObj.QueryRow("SELECT "+cols+" FROM "+dbName+"."+tableName+" WHERE "+where, args...), nil
}

func findRowsBy(dbObj *sql.DB, dbName string, tableName string, cols string, where string, args ...interface{}) (*sql.Rows, error) {
    if dbObj == nil {
        return nil, fmt.Errorf("db wasn't initialized")
    }

    if where == "" {
        where = "1"
    }
    return dbObj.Query("SELECT "+cols+" FROM "+dbName+"."+tableName+" WHERE "+where, args...)
}

func updateBy(dbObj *sql.DB, dbName string, tableName string, set string, where string, args ...interface{}) (int64, error) {
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

func removeBy(dbObj *sql.DB, dbName string, tableName string, where string, args ...interface{}) (int64, error) {
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
