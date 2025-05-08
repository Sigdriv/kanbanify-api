package handler

import (
	"fmt"
	"kanbanify-api/db"
	"kanbanify-api/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func getIssues(c *gin.Context) {
	conn, err := db.Connect()
	if err != nil {
		logrus.Error("Error connecting to database << ", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to connect to database"})
		return
	}

	defer conn.Close()

	rows, err := conn.Query(c, "SELECT kanban_id, title, description, status, variant, created_at, updated_at FROM issues")
	if err != nil {
		logrus.Error("Error querying database << ", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to fetch issues"})
		return
	}
	defer rows.Close()

	var issues []model.Issue
	for rows.Next() {
		var issue model.Issue
		err := rows.Scan(&issue.KanbanID, &issue.Title, &issue.Description, &issue.Status, &issue.Variant, &issue.CreatedAt, &issue.UpdatedAt)
		if err != nil {
			logrus.Error("Failed to scan row << ", err)
			c.IndentedJSON(500, gin.H{"error": "Failed to scan row"})
			return
		}
		issues = append(issues, issue)
	}

	err = rows.Err()
	if err != nil {
		logrus.Error("Error iterating over rows << ", err)
		c.IndentedJSON(500, gin.H{"error": "Error iterating over rows"})
		return
	}

	if len(issues) == 0 {
		issues = []model.Issue{}
	}

	logrus.Info("Fetched issues from database")
	c.IndentedJSON(200, issues)
}

func createIssue(c *gin.Context) {
	var newIssue model.Issue

	err := c.BindJSON(&newIssue)
	if err != nil {
		logrus.Error("Error binding JSON << ", err)
		c.IndentedJSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	newIssue.Variant, err = classifyIssue(newIssue)
	if err != nil {
		logrus.Error("Error classifying issue << ", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to classify issue"})
		return
	}

	conn, err := db.Connect()
	if err != nil {
		logrus.Error("Error connecting to database << ", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer conn.Close()

	var id int
	err = conn.QueryRow(c,
		"INSERT INTO issues (title, kanban_id, description, status, variant) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		newIssue.Title, "", newIssue.Description, newIssue.Status, newIssue.Variant).Scan(&id)

	if err != nil {
		logrus.Error("Error inserting issue into database << ", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to create issue"})
		return
	}

	newIssue.KanbanID = fmt.Sprintf("KAN-%d", id)
	_, err = conn.Exec(c, "UPDATE issues SET kanban_id = $1 WHERE id = $2", newIssue.KanbanID, id)
	if err != nil {
		logrus.Error("Error updating issue ID << ", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to update issue ID"})
		return
	}

	logrus.Info("Created new issue in database with ID << ", id)
	c.IndentedJSON(201, newIssue)

}

func updateIssue(c *gin.Context) {
	var updatedIssue model.Issue

	err := c.BindJSON(&updatedIssue)
	if err != nil {
		logrus.Error("Error binding JSON << ", err)
		c.IndentedJSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	conn, err := db.Connect()
	if err != nil {
		logrus.Error("Error connecting to database << ", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to connect to database"})
		return
	}

	defer conn.Close()

	var issueID int
	err = conn.QueryRow(c, "SELECT id FROM issues WHERE kanban_id = $1", updatedIssue.KanbanID).Scan(&issueID)
	if err != nil {
		logrus.Error("Error fetching issue ID << ", err)
		c.IndentedJSON(500, gin.H{"error": "Issue not found"})
		return
	}

	if issueID == 0 {
		logrus.Error("Issue not found << ", updatedIssue.KanbanID)
		c.IndentedJSON(500, gin.H{"error": "Issue not found"})
		return
	}

	_, err = conn.Exec(c, "UPDATE issues SET title = $1, description = $2, status = $3, variant = $4 WHERE kanban_id = $5",
		updatedIssue.Title, updatedIssue.Description, updatedIssue.Status, updatedIssue.Variant, updatedIssue.KanbanID)

	if err != nil {
		logrus.Error("Error updating issue in database << ", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to update issue"})
		return
	}

	logrus.Info("Updated issue in database with ID << ", updatedIssue.KanbanID)
	c.IndentedJSON(200, gin.H{"message": "Issue updated successfully"})
}

func deleteIssue(c *gin.Context) {
	kanbanID := c.Param("id")

	conn, err := db.Connect()
	if err != nil {
		logrus.Error("Error connecting to database << ", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to connect to database"})
		return
	}

	defer conn.Close()

	_, err = conn.Exec(c, "DELETE FROM issues WHERE kanban_id = $1", kanbanID)
	if err != nil {
		logrus.Error("Error deleting issue from database << ", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to delete issue"})
		return
	}

	logrus.Info("Deleted issue from database with ID << ", kanbanID)
	c.IndentedJSON(200, gin.H{"message": "Issue deleted successfully"})
}
