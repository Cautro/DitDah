package lesson

type LessonEntity struct {
	Id int `json:"id"`
	Order int `json:"order"`
	Title string `json:"title"`
	Description string `json:"description"`
	XPReward int `json:"xpReward"`
	Lang LanguageLesson `json:"language"`
	Type TypeLesson `json:"type"`
}

type LanguageLesson string
const (
	Ru LanguageLesson = "ru"
	En LanguageLesson = "en"
)

type TypeLesson string
const (
	TypeText TypeLesson = "text"
	TypeTap TypeLesson = "tap"
	TypeListen TypeLesson = "listen"
	TypeQuiz TypeLesson = "quiz"
)

type LessonSymbols struct {
	LessonId int `json:"lessonId"`
	SymbolId int `json:"symbolId"`
}

type Symbol struct {
	Symbol rune `json:"symbol"`
	Morse string `json:"morse"`
}