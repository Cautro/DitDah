package lesson

import "database/sql"

type repositoryRepo struct {
	db *sql.DB
}

func New(db *sql.DB) *repositoryRepo {
	return &repositoryRepo{
		db: db,
	}
}

func (r *repositoryRepo) GetLessonById(lessonId int) (LessonEntity, error) {
	var lesson LessonEntity
	err := r.db.QueryRow("SELECT id, order_num, title, description, xp_reward, language FROM lessons WHERE id = $1", lessonId).
		Scan(&lesson.Id, &lesson.Order, &lesson.Title, &lesson.Description, &lesson.XPReward, &lesson.Lang)
	if err != nil {
		return LessonEntity{}, err
	}
	rows, err := r.db.Query("SELECT id, lesson_id, order_num, type, payload FROM tasks WHERE lesson_id = $1", lessonId)
	if err != nil {
		return LessonEntity{}, err
	}
	defer rows.Close()

	lesson.Tasks = []TaskEntity{}

	for rows.Next() {
		t := TaskEntity{}
		var payloadBytes []byte
		if err := rows.Scan(&t.Id, &t.LessonId, &t.Order, &t.Type, &payloadBytes); err != nil {
			return LessonEntity{}, err
		}
		t.Payload = payloadBytes
		lesson.Tasks = append(lesson.Tasks, t)
	}
	return lesson, nil
}

func (r *repositoryRepo) GetAllLessons() ([]LessonEntity, error) {
	rows, err := r.db.Query("SELECT id, order_num, title, description, xp_reward, language FROM lessons")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lessons := []LessonEntity{}
	lessonMap := make(map[int]*LessonEntity) 

	for rows.Next() {
		lesson := LessonEntity{}
		if err := rows.Scan(&lesson.Id, &lesson.Order, &lesson.Title, &lesson.Description, &lesson.XPReward, &lesson.Lang); err != nil {
			return nil, err
		}
		lesson.Tasks = []TaskEntity{} 
		lessons = append(lessons, lesson)
		lessonMap[lesson.Id] = &lessons[len(lessons)-1]
	}

	taskRows, err := r.db.Query("SELECT id, lesson_id, order_num, type, payload FROM tasks")
	if err != nil {
		return nil, err
	}
	defer taskRows.Close()

	for taskRows.Next() {
		t := TaskEntity{}
		var payloadBytes []byte
		if err := taskRows.Scan(&t.Id, &t.LessonId, &t.Order, &t.Type, &payloadBytes); err != nil {
			return nil, err
		}
		t.Payload = payloadBytes 
		
		if lesson, exists := lessonMap[t.LessonId]; exists {
			lesson.Tasks = append(lesson.Tasks, t)
		}
	}

	return lessons, nil
}

func (r *repositoryRepo) AddLesson(lesson LessonEntity) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var lessonID int
	err = tx.QueryRow(
		"INSERT INTO lessons (order_num, title, description, xp_reward, language) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		lesson.Order, lesson.Title, lesson.Description, lesson.XPReward, lesson.Lang,
	).Scan(&lessonID)
	
	if err != nil {
		return err
	}

	for _, task := range lesson.Tasks {
		_, err = tx.Exec(
			"INSERT INTO tasks (lesson_id, order_num, type, payload) VALUES ($1, $2, $3, $4)",
			lessonID, task.Order, task.Type, task.Payload,
		)
		if err != nil {
			return err 
		}
	}

	return tx.Commit()
}

func (r *repositoryRepo) DeleteLesson(lessonId int) error {
	_, err := r.db.Exec("DELETE FROM lessons WHERE id = $1", lessonId)
	return err
}