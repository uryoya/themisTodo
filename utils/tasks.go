package utils

import (
	"../models"
	"log"
)

func TasksConvert(tasks []models.Task) []models.Task {
	for key, value := range tasks {
		tasks[key] = *TaskConvert(&value)
	}

	return tasks
}

func TaskConvert(task *models.Task) *models.Task {
	var e bool
	e, task.DeadlineMD = GetDateMD(task.Deadline)
	if e {
		log.Printf("Utils.TaskConvert")
	}

	task.LimitDate = DiffDay(task.Deadline)

	return task
}
