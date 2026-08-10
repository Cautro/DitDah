package user

import "time"

type UserEntity struct {
    Id           int        `json:"id"`
    Username     string     `json:"username"`
    Password     string     `json:"-"`
    XP           int        `json:"xp"`
    NeedXP       int        `json:"needXp"`
    Level        int        `json:"level"`

    LessonsDoneEn int `json:"lessonsDoneEn"`
    LessonsDoneRu int `json:"lessonsDoneRu"`

    Elo          int `json:"elo"`
    DuelsWin     int `json:"duelsWin"`
    DuelMaxScore int `json:"duelMaxScore"`
    Coins        int `json:"coins"`

    DayStreak    int `json:"dayStreak"`
    AnswerStreak int `json:"answerStreak"`

    LastLogin     *time.Time `json:"lastLogin"`
    InvitedBy     *int       `json:"invitedBy"`
    ReferralCode  *string    `json:"referralCode"`
    RegisteredDate time.Time `json:"registeredDate"`

    Friends              []int                `json:"friends"`
    UnlockedAchievements []string             `json:"unlockedAchievements"`
    SymbolStats          []UserSymbolStatEntity `json:"symbolStats"`
}

type UserSymbolStatEntity struct {
	UserId            int       `json:"user_id"`
	Symbol            string    `json:"symbol"`
	Correct           int       `json:"correct"`
	Wrong             int       `json:"wrong"`
	
	ConsecutiveErrors int       `json:"consecutive_errors"` // ошибки подряд, сбрасывается при правильном ответе
	Weight            float64   `json:"weight"` 
	LastPracticed     time.Time `json:"last_practiced"`
}

type UserFriendEntity struct {
	UserId     int   `json:"userId"`
	FriendId   int   `json:"friendId"`
}

type UserAchievementEntity struct {
	UserId        int       `json:"userId"`
	AchievementId int       `json:"achievementId"`
	UnlockedAt    time.Time `json:"unlockedAt"`
}

type UserItemEntity struct {
	UserId    int    `json:"userId"`
	ItemId    int    `json:"itemId"`
	Amount    int    `json:"amount"`
}

type UserRegisterDTO struct {
    Username     string     `json:"username"`
    Password     string     `json:"password"`
}

type RefreshToken struct {
	ID        int
	UserID    int
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}