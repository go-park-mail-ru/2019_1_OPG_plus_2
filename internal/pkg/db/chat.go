package db

import (
    "database/sql"
    "time"

    "2019_1_OPG_plus_2/internal/pkg/models"
)

func CreateMessage(data models.ChatMessage) (models.ChatMessage, error) {
    userId := int64(0)
    if data.Username != "" {
        row, err := QueryRow("SELECT a.id, c.avatar FROM "+AuthDbName+"."+AuthUsersTable+" AS a "+
            "JOIN "+CoreDbName+"."+CoreUsersTable+" AS c ON a.id = c.id WHERE a.username = ?", data.Username)
        if err != nil {
            return data, err
        }
        err = row.Scan(&userId, &data.Avatar)
        if err != nil && err != sql.ErrNoRows {
            return data, err
        }
    }

    typeId := int64(0)
    row, err := findRowBy(ChatDbName, ChatTypesTable, "id", "type = ?", data.Type)
    if err != nil {
        return data, err
    }
    err = row.Scan(&typeId)
    if err == sql.ErrNoRows {
        return data, models.NotFound
    }

    data.Id, err = insert(ChatDbName, ChatMessagesTable, "created, user_id, type_id, content", "?, ?, ?, ?",
        time.Time(data.Date), userId, typeId, data.Content)
    return data, err
}

func GetMessages(limit, offset int64) (chatMessages []models.ChatMessage, count uint64, err error) {
    row, err := QueryRow("SELECT COUNT(id) FROM " + ChatDbName + "." + ChatMessagesTable)
    if err != nil {
        return
    }
    err = row.Scan(&count)

    rows, err := Query(
        "SELECT cm.id, CASE WHEN a.username IS NULL THEN '' ELSE a.username END, "+
            "CASE WHEN c.avatar IS NULL THEN '' ELSE c.avatar END, cm.created, ct.type, cm.content "+
            "FROM "+ChatDbName+"."+ChatMessagesTable+" AS cm "+
            "JOIN "+ChatDbName+"."+ChatTypesTable+" AS ct ON cm.type_id = ct.id "+
            "LEFT JOIN "+AuthDbName+"."+AuthUsersTable+" AS a ON cm.user_id = a.id "+
            "LEFT JOIN "+CoreDbName+"."+CoreUsersTable+" AS c ON cm.user_id = c.id "+
            "ORDER BY cm.created DESC, cm.id DESC LIMIT ? OFFSET ?", limit, offset)
    if err != nil {
        return
    }

    defer rows.Close()
    for rows.Next() {
        chatMessage := models.ChatMessage{}
        var strDate string
        err = rows.Scan(&chatMessage.Id, &chatMessage.Username, &chatMessage.Avatar,
            &strDate, &chatMessage.Type, &chatMessage.Content)
        if err != nil {
            return
        }

        var timeDate time.Time
        timeDate, err = time.Parse("2006-01-02 15:04:05", strDate)
        if err != nil {
            return
        }
        chatMessage.Date = models.JSONTime(timeDate)

        chatMessages = append(chatMessages, chatMessage)
    }
    return
}
