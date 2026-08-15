package user

import (
	"context"
	"database/sql"
	"time"
)

type postgresRepo struct {
	db *sql.DB
}

func New(db *sql.DB) UserRepository {
	return &postgresRepo{
		db: db,
	}
}

func (r *postgresRepo) GetAllUsers(ctx context.Context) ([]*UserEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	const query = `SELECT
			id,
			username,
			xp,
			need_xp,
			level,
			lessons_done_en,
			lessons_done_ru,
			elo,
			duels_win,
			duel_max_score,
			coins,
			day_streak,
			answer_streak,
			last_login,
			invited_by,
			referral_code,
			registered_date
		FROM users;`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*UserEntity
	for rows.Next() {
		u := &UserEntity{}
		err := rows.Scan(
			&u.Id,
			&u.Username,
			&u.XP,
			&u.NeedXP,
			&u.Level,
			&u.LessonsDoneEn,
			&u.LessonsDoneRu,
			&u.Elo,
			&u.DuelsWin,
			&u.DuelMaxScore,
			&u.Coins,
			&u.DayStreak,
			&u.AnswerStreak,
			&u.LastLogin,
			&u.InvitedBy,
			&u.ReferralCode,
			&u.RegisteredDate,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *postgresRepo) GetFullUserById(ctx context.Context, userID int) (*UserEntity, error) {
    user, err := r.getUserById(ctx, userID)
    if err != nil {
        return nil, err
    }

    friends, err := r.getUserFriendsById(ctx, userID)
    if err != nil {
        return nil, err
    }

    achievements, err := r.getUserAchievementsById(ctx, userID)
    if err != nil {
        return nil, err
    }

    stats, err := r.getUserSymbolStatsById(ctx, userID)
    if err != nil {
        return nil, err
    }

    user.Friends = friends
    user.UnlockedAchievements = achievements
    user.SymbolStats = stats

    return user, nil
}

func (r *postgresRepo) GetUserByLogin(ctx context.Context, username string) (*UserEntity, error) {
	output, err := r.getUserByLogin(ctx, username)
	if err != nil {
		return &UserEntity{}, err
	}

	return output, nil
}

func (r *postgresRepo) Register(ctx context.Context, username, passwordHash string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	const query = `
		INSERT INTO users (username, password)
		VALUES ($1, $2)`
	
	_, err := r.db.ExecContext(ctx, query, username, passwordHash)
	return err
}


func (r *postgresRepo) SaveRefreshToken(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`, userID, token, expiresAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *postgresRepo) GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, token, expires_at, created_at
		FROM refresh_tokens
		WHERE token = $1
	`, token)

	var rt RefreshToken

	err := row.Scan(
		&rt.ID,
		&rt.UserID,
		&rt.Token,
		&rt.ExpiresAt,
		&rt.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &rt, nil
}

func (r *postgresRepo) DeleteRefreshToken(ctx context.Context, token string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, `
		DELETE FROM refresh_tokens
		WHERE token = $1
	`, token)
	if err != nil {
		return err
	}

	return nil
}

func (r *postgresRepo) getUserByLogin(ctx context.Context, username string) (*UserEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	const query = `SELECT
			id,
			username,
			password,
			xp,
			need_xp,
			level,
			lessons_done_en,
			lessons_done_ru,
			elo,
			duels_win,
			duel_max_score,
			coins,
			day_streak,
			answer_streak,
			last_login,
			invited_by,
			referral_code,
			registered_date
		FROM users
		WHERE username = $1;`

	u := &UserEntity{}

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&u.Id,
		&u.Username,
		&u.Password,
		&u.XP,
		&u.NeedXP,
		&u.Level,
		&u.LessonsDoneEn,
		&u.LessonsDoneRu,
		&u.Elo,
		&u.DuelsWin,
		&u.DuelMaxScore,
		&u.Coins,
		&u.DayStreak,
		&u.AnswerStreak,
		&u.LastLogin,
		&u.InvitedBy,
		&u.ReferralCode,
		&u.RegisteredDate,
	)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *postgresRepo) getUserById(ctx context.Context, id int) (*UserEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	const query = `SELECT
			id,
			username,
			xp,
			need_xp,
			level,
			lessons_done_en,
			lessons_done_ru,
			elo,
			duels_win,
			duel_max_score,
			coins,
			day_streak,
			answer_streak,
			last_login,
			invited_by,
			referral_code,
			registered_date
		FROM users
		WHERE id = $1;`

	u := &UserEntity{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.Id,
		&u.Username,
		&u.XP,
		&u.NeedXP,
		&u.Level,
		&u.LessonsDoneEn,
		&u.LessonsDoneRu,
		&u.Elo,
		&u.DuelsWin,
		&u.DuelMaxScore,
		&u.Coins,
		&u.DayStreak,
		&u.AnswerStreak,
		&u.LastLogin,
		&u.InvitedBy,
		&u.ReferralCode,
		&u.RegisteredDate,
	)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *postgresRepo) getUserFriendsById(ctx context.Context, id int) ([]int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	const query = `SELECT
		CASE
			WHEN user_id = $1 THEN friend_id
			ELSE user_id
		END AS friend_id
	FROM user_friends
	WHERE user_id = $1
	OR friend_id = $1;`
	
	rows, err := r.db.QueryContext(ctx, query, id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    friends := make([]int, 0)

    for rows.Next() {
        var friendID int

        if err := rows.Scan(&friendID); err != nil {
            return nil, err
        }

        friends = append(friends, friendID)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return friends, nil
}

func (r *postgresRepo) getUserAchievementsById(ctx context.Context, id int) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	const query = `
        SELECT a.code
        FROM user_achievements ua
        JOIN achievements a ON a.id = ua.achievement_id
        WHERE ua.user_id = $1;
    `

    rows, err := r.db.QueryContext(ctx, query, id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    achievements := make([]string, 0)

    for rows.Next() {
        var code string

        if err := rows.Scan(&code); err != nil {
            return nil, err
        }

        achievements = append(achievements, code)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return achievements, nil
}

func (r *postgresRepo) getUserSymbolStatsById(ctx context.Context, id int) ([]UserSymbolStatEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	const query = `SELECT
			user_id,
			symbol,
			correct,
			wrong,
			consecutive_errors,
			weight,
			last_practiced
		FROM user_symbol_stats
		WHERE user_id = $1;`
	
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []UserSymbolStatEntity

	for rows.Next() {
		var s UserSymbolStatEntity
		err := rows.Scan(
			&s.UserId,
			&s.Symbol,
			&s.Correct,
			&s.Wrong,
			&s.ConsecutiveErrors,
			&s.Weight,
			&s.LastPracticed,
		)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}
