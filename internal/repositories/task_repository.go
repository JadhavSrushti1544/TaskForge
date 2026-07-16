package repositories

import (
    "database/sql"
    "github.com/JadhavSrushti1544/TaskForge/internal/models"
)

type TaskRepository struct {
    db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
    return &TaskRepository{db: db}
}

// CreateTask inserts a new task
func (r *TaskRepository) CreateTask(userID int, title, description, priority string) (*models.Task, error) {
    var task models.Task
    err := r.db.QueryRow(
        `INSERT INTO tasks (user_id, title, description, priority, status) 
         VALUES ($1, $2, $3, $4, 'pending') 
         RETURNING id, user_id, title, description, status, priority, created_at, updated_at`,
        userID, title, description, priority,
    ).Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.CreatedAt, &task.UpdatedAt)
    
    return &task, err
}

// GetTasksByUser retrieves all tasks for a user with pagination
func (r *TaskRepository) GetTasksByUser(userID int, limit, offset int) ([]models.Task, error) {
    rows, err := r.db.Query(
        `SELECT id, user_id, title, description, status, priority, created_at, updated_at 
         FROM tasks 
         WHERE user_id = $1 
         ORDER BY created_at DESC 
         LIMIT $2 OFFSET $3`,
        userID, limit, offset,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var tasks []models.Task
    for rows.Next() {
        var task models.Task
        rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.CreatedAt, &task.UpdatedAt)
        tasks = append(tasks, task)
    }
    return tasks, nil
}

// GetTaskByID retrieves a single task
func (r *TaskRepository) GetTaskByID(taskID, userID int) (*models.Task, error) {
    var task models.Task
    err := r.db.QueryRow(
        `SELECT id, user_id, title, description, status, priority, created_at, updated_at 
         FROM tasks 
         WHERE id = $1 AND user_id = $2`,
        taskID, userID,
    ).Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.CreatedAt, &task.UpdatedAt)
    
    return &task, err
}

// UpdateTask patches a task
func (r *TaskRepository) UpdateTask(taskID, userID int, updates map[string]interface{}) error {
    query := "UPDATE tasks SET updated_at = NOW()"
    args := []interface{}{}
    argIndex := 1
    
    if title, ok := updates["title"].(string); ok {
        query += ", title = $" + string(rune('0'+argIndex))
        args = append(args, title)
        argIndex++
    }
    if status, ok := updates["status"].(string); ok {
        query += ", status = $" + string(rune('0'+argIndex))
        args = append(args, status)
        argIndex++
    }
    if priority, ok := updates["priority"].(string); ok {
        query += ", priority = $" + string(rune('0'+argIndex))
        args = append(args, priority)
        argIndex++
    }
    
    query += " WHERE id = $" + string(rune('0'+argIndex)) + " AND user_id = $" + string(rune('0'+argIndex+1))
    args = append(args, taskID, userID)
    
    _, err := r.db.Exec(query, args...)
    return err
}

// DeleteTask removes a task
func (r *TaskRepository) DeleteTask(taskID, userID int) error {
    _, err := r.db.Exec("DELETE FROM tasks WHERE id = $1 AND user_id = $2", taskID, userID)
    return err
}
