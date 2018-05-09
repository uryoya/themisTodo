package module

import (
	"database/sql"
	"log"
	"../models"
	"time"
	"sync"
)

type TasksModule struct {
	db     *sql.DB
	dbLock sync.Mutex
}

func NewTaskModule(db *sql.DB) *TasksModule {
	return &TasksModule{
		db,
		sync.Mutex{},
	}
}

func (self *TasksModule) GetLastId(projectId int) int {
	self.dbLock.Lock()
	defer self.dbLock.Unlock()
	rows, err := self.db.Query("SELECT `id` FROM `todo_list` WHERE `project` = ? ORDER BY `id` DESC LIMIT 0,1;", projectId)

	if err != nil {
		return 0
	}

	if !rows.Next() {
		return 0
	}

	var lastId int
	if err := rows.Scan(&lastId); err != nil {
		log.Printf("TasksModule.GetLastId Error: %+v\n", err)
		return 0
	}

	return lastId
}

func (self *TasksModule) Add(task *models.Task) *models.Task {
	self.dbLock.Lock()
	defer self.dbLock.Unlock()

	stmt, err := self.db.Prepare("INSERT INTO `todo_list` (`id`, `project`, `name`, `creator`, `status`, `deadline`, `description`, `createDate`) VALUE (?, ?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		log.Printf("TasksModule.Add Error: %+v", err)
		return nil
	}

	defer stmt.Close()

	now := time.Now().UnixNano()
	_, err = stmt.Exec(task.TaskId, task.ProjectId, task.Name, task.Creator, task.Status, task.Deadline, task.Description, now)
	if err != nil {
		log.Printf("TasksModule.Add Error: %+v", err)
		return nil
	}

	task.CreateDate = now

	return task
}

func (self *TasksModule) GetList(projectId int) (error bool, list []models.Task) {
	list = []models.Task{}

	rows, err := self.db.Query("SELECT `id`, `todo_list`.`name`, `creator`, `status`, `deadline`, `description`, `createDate`, `displayName` "+
		"FROM `todo_list` INNER JOIN `users` ON `users`.`uuid` = `creator` WHERE `project` = ? ORDER BY `id` ASC;", projectId)

	if err != nil {
		return true, nil
	}

	for rows.Next() {
		listOne := models.Task{}
		if err := rows.Scan(&listOne.TaskId, &listOne.Name, &listOne.Creator, &listOne.Status, &listOne.Deadline, &listOne.Description, &listOne.CreateDate,  &listOne.CreatorName); err != nil {
			log.Printf("TasksModule.GetList Error: %+v\n", err)
			return true, nil
		}
		list = append(list, listOne)
	}

	return false, list
}