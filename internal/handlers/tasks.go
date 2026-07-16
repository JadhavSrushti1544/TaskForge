package handlers

import (
    "net/http"
    "strconv"
    "task-queue/models"
    "task-queue/repositories"
    "github.com/gin-gonic/gin"
)

type TaskHandler struct {
    taskRepo *repositories.TaskRepository
}

func NewTaskHandler(taskRepo *repositories.TaskRepository) *TaskHandler {
    return &TaskHandler{taskRepo: taskRepo}
}

// CreateTask POST /tasks
func (h *TaskHandler) CreateTask(c *gin.Context) {
    userID := c.GetInt("user_id")
    
    var req models.CreateTaskRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    task, err := h.taskRepo.CreateTask(userID, req.Title, req.Description, req.Priority)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
        return
    }
    
    c.JSON(http.StatusCreated, task)
}

// GetTasks GET /tasks?limit=10&offset=0
func (h *TaskHandler) GetTasks(c *gin.Context) {
    userID := c.GetInt("user_id")
    
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
    
    if limit > 100 {
        limit = 100
    }
    
    tasks, err := h.taskRepo.GetTasksByUser(userID, limit, offset)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"tasks": tasks, "total": len(tasks)})
}

// GetTask GET /tasks/:id
func (h *TaskHandler) GetTask(c *gin.Context) {
    userID := c.GetInt("user_id")
    taskID, _ := strconv.Atoi(c.Param("id"))
    
    task, err := h.taskRepo.GetTaskByID(taskID, userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
        return
    }
    
    c.JSON(http.StatusOK, task)
}

// UpdateTask PATCH /tasks/:id
func (h *TaskHandler) UpdateTask(c *gin.Context) {
    userID := c.GetInt("user_id")
    taskID, _ := strconv.Atoi(c.Param("id"))
    
    var req models.UpdateTaskRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    updates := make(map[string]interface{})
    if req.Title != nil {
        updates["title"] = *req.Title
    }
    if req.Status != nil {
        updates["status"] = *req.Status
    }
    if req.Priority != nil {
        updates["priority"] = *req.Priority
    }
    
    if err := h.taskRepo.UpdateTask(taskID, userID, updates); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Task updated"})
}

// DeleteTask DELETE /tasks/:id
func (h *TaskHandler) DeleteTask(c *gin.Context) {
    userID := c.GetInt("user_id")
    taskID, _ := strconv.Atoi(c.Param("id"))
    
    if err := h.taskRepo.DeleteTask(taskID, userID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}