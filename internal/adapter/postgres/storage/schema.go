package storage

func (s *Storage) initSchema() error {
    if err := s.initUser(); err != nil {
        return err
    }

    if err := s.initUserFriends(); err != nil {
        return err
    }

    if err := s.initUserSymbolStat(); err != nil {
        return err
    }

    if err := s.initAchievement(); err != nil {
        return err
    }

    if err := s.initUserAchievement(); err != nil {
        return err
    }

    if err := s.initItem(); err != nil {
        return err
    }

    if err := s.initUserItem(); err != nil {
        return err
    }

	if err := s.initLesson(); err != nil {
		return err
	}

	if err := s.initTask(); err != nil {
		return err
	}

	if err := s.initHTTPOnlyStorage(); err != nil {
		return err
	}

    return nil
}

func (s *Storage) initUser() error {
	query := `CREATE TABLE IF NOT EXISTS users (
		id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

		username VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255),

		is_admin BOOLEAN NOT NULL DEFAULT FALSE,

		xp INT NOT NULL DEFAULT 0,
		need_xp INT NOT NULL DEFAULT 0,
		level INT NOT NULL DEFAULT 1,

		lessons_done_en INT NOT NULL DEFAULT 0,
		lessons_done_ru INT NOT NULL DEFAULT 0,

		elo INT NOT NULL DEFAULT 500,
		duels_win INT NOT NULL DEFAULT 0,
		duel_max_score INT NOT NULL DEFAULT 0,

		coins INT NOT NULL DEFAULT 0,

		day_streak INT NOT NULL DEFAULT 0,
		answer_streak INT NOT NULL DEFAULT 0,

		last_login TIMESTAMPTZ NOT NULL DEFAULT NOW(),

		invited_by INT REFERENCES users(id) ON DELETE SET NULL,
		referral_code VARCHAR(50) UNIQUE,

		registered_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

		CHECK (xp >= 0),
		CHECK (need_xp >= 0),
		CHECK (level >= 1),
		CHECK (lessons_done_en >= 0),
		CHECK (lessons_done_ru >= 0),
		CHECK (elo >= 0),
		CHECK (duels_win >= 0),
		CHECK (duel_max_score >= 0),
		CHECK (coins >= 0),
		CHECK (day_streak >= 0),
		CHECK (answer_streak >= 0)
	);`

	if _, err := s.DB.Exec(query); err != nil {
		return err
	}

	if _, err := s.DB.Exec(`CREATE INDEX IF NOT EXISTS idx_users_id ON users(id)`); err != nil {
		return err
	}

	return nil
}

func (s *Storage) initUserSymbolStat() error {
	query := `CREATE TABLE IF NOT EXISTS user_symbol_stats (
		user_id INT NOT NULL
			REFERENCES users(id)
			ON DELETE CASCADE,

		symbol VARCHAR(10) NOT NULL,

		correct INT NOT NULL DEFAULT 0,
		wrong INT NOT NULL DEFAULT 0,
		consecutive_errors INT NOT NULL DEFAULT 0,

		weight DOUBLE PRECISION NOT NULL DEFAULT 1.0,

		last_practiced TIMESTAMPTZ NOT NULL DEFAULT NOW(),

		PRIMARY KEY (user_id, symbol),

		CHECK (correct >= 0),
		CHECK (wrong >= 0),
		CHECK (consecutive_errors >= 0),
		CHECK (weight >= 0)
	);`

	_, err := s.DB.Exec(query)
	return err
}

func (s *Storage) initUserAchievement() error {
	query := `CREATE TABLE IF NOT EXISTS user_achievements (
		user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		achievement_id INT NOT NULL REFERENCES achievements(id) ON DELETE CASCADE,
		unlocked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

		PRIMARY KEY (user_id, achievement_id)
	);`

	_, err := s.DB.Exec(query)
	return err
}

func (s *Storage) initUserItem() error {
	query := `CREATE TABLE IF NOT EXISTS user_items (
		user_id INT NOT NULL
			REFERENCES users(id)
			ON DELETE CASCADE,

		item_id INT NOT NULL
			REFERENCES items(id)
			ON DELETE CASCADE,
		
		amount INT NOT NULL DEFAULT 1,
		
		PRIMARY KEY (user_id, item_id),
		CHECK (amount > 0)
	);`

	_, err := s.DB.Exec(query)
	return err
}

func (s *Storage) initUserFriends() error {
	query := `CREATE TABLE IF NOT EXISTS user_friends (
		user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		friend_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

		PRIMARY KEY (user_id, friend_id),

		CHECK (user_id < friend_id)
	);`

	_, err := s.DB.Exec(query)
	return err
}

func (s *Storage) initAchievement() error {
	query := `CREATE TABLE IF NOT EXISTS achievements (
		id SERIAL PRIMARY KEY,
		code VARCHAR(100) UNIQUE NOT NULL DEFAULT '',
		name VARCHAR(255) NOT NULL DEFAULT '',
		description VARCHAR(500) NOT NULL DEFAULT '',
		xp_reward INT NOT NULL DEFAULT 0
	);`

	_, err := s.DB.Exec(query)
	return err
}

func (s *Storage) initItem() error {
	query := `CREATE TABLE IF NOT EXISTS items(
		id SERIAL PRIMARY KEY,
		title VARCHAR(64) NOT NULL DEFAULT '',
		description VARCHAR(500) NOT NULL DEFAULT ''
	);`

	_, err := s.DB.Exec(query)
	return err
}

func (s *Storage) initLesson() error {
	const query = `CREATE TABLE IF NOT EXISTS lessons (
		id SERIAL PRIMARY KEY,
		order_num INT NOT NULL,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		xp_reward INT NOT NULL DEFAULT 10,
		language VARCHAR(10) NOT NULL
	);`

	_, err := s.DB.Exec(query)
	return err
}

func (s *Storage) initTask() error {
	const query = `CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		lesson_id INT NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
		order_num INT NOT NULL,
		type VARCHAR(50) NOT NULL,
		payload JSONB NOT NULL
	);`
	_, err := s.DB.Exec(query)
	return err
}

func (s *Storage) initHTTPOnlyStorage() error {

	query := `
	CREATE TABLE IF NOT EXISTS refresh_tokens (
		id BIGSERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		token TEXT NOT NULL UNIQUE,
		expires_at TIMESTAMPTZ NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(Id) ON DELETE CASCADE
	);`

	if _, err := s.DB.Exec(query); err != nil {
		return err
	}
	if _, err := s.DB.Exec(`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);`); err != nil {
		return err
	}

	return nil
}
