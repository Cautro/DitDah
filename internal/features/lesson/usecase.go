package lesson

import (
	"context"
	user "ditdah/internal/features/user"
	"encoding/json"
	"errors"

	"ditdah/pkg/check"
)

type LessonUseCase struct {
	lessonRepo LessonRepository
	userRepo user.UserRepository
}

func NewLessonUseCase(repo LessonRepository, userRepo user.UserRepository) *LessonUseCase {
	return &LessonUseCase{
		lessonRepo: repo,
		userRepo: userRepo,
	}
}

func (l *LessonUseCase) GetLessonById(lessonId int) (LessonEntity, error) {
	lesson, err := l.lessonRepo.GetLessonById(lessonId)
	if err != nil {
		return LessonEntity{}, err
	}

	return lesson, nil
}

func (l *LessonUseCase) GetAllLessons() ([]LessonEntity, error) {
	return l.lessonRepo.GetAllLessons()
}

func (l *LessonUseCase) AddLesson(lesson LessonEntity, userId int) error {
	user, err := l.userRepo.GetFullUserById(context.Background(), userId)
	if err != nil {
		return err
	}

	if !check.CheckToAdmin(user) {
		return errors.New("User is not an admin")
	}
	
	if err := l.validateLesson(lesson); err != nil {
		return err
	}

	return l.lessonRepo.AddLesson(lesson)
}

func (l *LessonUseCase) DeleteLesson(lessonId int, userId int) error {
	user, err := l.userRepo.GetFullUserById(context.Background(), userId)
	if err != nil {
		return err
	}
	if !check.CheckToAdmin(user) {
		return errors.New("User is not an admin")
	}

	return l.lessonRepo.DeleteLesson(lessonId)
}

func (l *LessonUseCase) validateLesson(lesson LessonEntity) error {
	if lesson.Title == "" {
		return errors.New("Lesson title cannot be empty")
	} else if lesson.Description == "" {
		return errors.New("Lesson description cannot be empty")
	} else if lesson.XPReward < 0 {
		return errors.New("Lesson XP reward cannot be negative")
	} else if lesson.Lang != Ru && lesson.Lang != En {
		return errors.New("Lesson language must be either 'ru' or 'en'")
	} else if len(lesson.Tasks) <= 0 {
		return errors.New("Lesson must have at least one task")
	} else if len(lesson.Tasks) > 20 || len(lesson.Tasks) < 7 {
		return errors.New("Lesson must have between 7 and 20 tasks")
	}

	for _, task := range lesson.Tasks {
		_, err := parseTaskPayload(task)
		if err != nil {
			return errors.New("Invalid payload in task")
		}
	}

	return nil
}

func parseTaskPayload(task TaskEntity) (any, error) {
	switch task.Type {
	case TypeQuiz:
		var p QuizPayload
		err := json.Unmarshal(task.Payload, &p)
		return p, err
	case TypeListen:
		var p ListenPayload
		err := json.Unmarshal(task.Payload, &p)
		return p, err
	case TypeTap:
		var p TapPayload
		err := json.Unmarshal(task.Payload, &p)
		return p, err
	case TypeText:
		var p TextPayload
		err := json.Unmarshal(task.Payload, &p)
		return p, err
	default:
		return nil, errors.New("Unknown task type")
	}
}