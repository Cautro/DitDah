package lesson

type LessonRepository interface {
	GetAllLessons() ([]LessonEntity, error)
	GetLessonById(lessonId int) (LessonEntity, error)
	AddLesson(lesson LessonEntity) error
	DeleteLesson(lessonId int) error
}