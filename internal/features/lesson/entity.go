package lesson

import "encoding/json"

type Language string

const (
	Ru Language = "ru"
	En Language = "en"
)

type TaskType string

const (
	TypeText   TaskType = "text"
	TypeTap    TaskType = "tap"
	TypeListen TaskType = "listen"
	TypeQuiz   TaskType = "quiz"
)

type LessonEntity struct {
	Id          int          `json:"id"`
	Order       int          `json:"order"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	XPReward    int          `json:"xpReward"`
	Lang        Language     `json:"language"`
	Tasks       []TaskEntity `json:"tasks,omitempty"`
}

type TaskEntity struct {
	Id       int             `json:"id"`
	LessonId int             `json:"lessonId"`
	Order    int             `json:"order"`
	Type     TaskType        `json:"type"`
	Payload  json.RawMessage `json:"payload"`
}

// another payloads for different task types

type TextPayload struct {
	Text  string `json:"text"`
	Morse string `json:"morse,omitempty"`
}

type QuizPayload struct {
	Question     string   `json:"question"`
	Options      []string `json:"options"`
	CorrectIndex int      `json:"correctIndex"`
}

type ListenPayload struct {
	MorseText     string `json:"morseText"`
	CorrectAnswer string `json:"correctAnswer"`
}

type TapPayload struct {
	Words 	string   `json:"words"`
	Morse 	string   `json:"morse"`
}