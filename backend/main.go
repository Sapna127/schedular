package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

type Event struct {
    StartTime int `json:"start_time"`
    EndTime   int `json:"end_time"`
}

type Scheduler struct {
    Events []Event
}

var scheduler = Scheduler{}

func (s *Scheduler) addEvent(event Event) bool {
    for _, e := range s.Events {
        if (event.StartTime < e.EndTime && event.StartTime >= e.StartTime) ||
            (event.EndTime > e.StartTime && event.EndTime <= e.EndTime) {
            return false
        }
    }
    s.Events = append(s.Events, event)
    return true
}

func addEventHandler(c *gin.Context) {
    var newEvent Event
    if err := c.ShouldBindJSON(&newEvent); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if newEvent.StartTime < 0 || newEvent.EndTime > 23 || newEvent.StartTime >= newEvent.EndTime {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event time"})
        return
    }

    if scheduler.addEvent(newEvent) {
        c.JSON(http.StatusOK, gin.H{"success": true})
    } else {
        c.JSON(http.StatusConflict, gin.H{"error": "Event overlaps with an existing event"})
    }
}

func getEventsHandler(c *gin.Context) {
    c.JSON(http.StatusOK, scheduler.Events)
}

func main() {
    r := gin.Default()

    r.Use(func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if c.Request.Method == http.MethodOptions {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }

        c.Next()
    })

    r.POST("/events", addEventHandler)
    r.GET("/events", getEventsHandler)
    r.Run(":8080")
}
