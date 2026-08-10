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

    return nil
}

func (s *Storage) initUser() error {
	query := `CREATE TABLE users (
		id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

		username VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255),

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

		last_login TIMESTAMPTZ,

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

	if _, err := s.DB.Exec(`CREATE INDEX idx_users_id ON users(id)`); err != nil {
		return err
	}

	return nil
}

func (s *Storage) initUserSymbolStat() error {
	query := `CREATE TABLE user_symbol_stats (
		user_id INT NOT NULL
			REFERENCES users(id)
			ON DELETE CASCADE,

		symbol VARCHAR(10) NOT NULL,

		correct INT NOT NULL DEFAULT 0,
		wrong INT NOT NULL DEFAULT 0,
		consecutive_errors INT NOT NULL DEFAULT 0,

		weight DOUBLE PRECISION NOT NULL DEFAULT 1.0,

		last_practiced TIMESTAMPTZ,

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
	query := `CREATE TABLE user_achievements (
		user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		achievement_id INT NOT NULL REFERENCES achievements(id) ON DELETE CASCADE,
		unlocked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

		PRIMARY KEY (user_id, achievement_id)
	);`

	_, err := s.DB.Exec(query)
	return err
}

func (s *Storage) initUserItem() error {
	query := `CREATE TABLE user_items (
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
	query := `CREATE TABLE user_friends (
		user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		friend_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

		PRIMARY KEY (user_id, friend_id),

		CHECK (user_id < friend_id)
	);`

	_, err := s.DB.Exec(query)
	return err
}

func (s *Storage) initAchievement() error {
	query := `CREATE TABLE achievements (
		id SERIAL PRIMARY KEY,
		code VARCHAR(100) UNIQUE NOT NULL,
		name VARCHAR(255) NOT NULL,
		description VARCHAR(500) NOT NULL,
		xp_reward INT NOT NULL DEFAULT 0
	);`

	_, err := s.DB.Exec(query)
	return err
}

func (s *Storage) initItem() error {
	query := `CREATE TABLE items(
		id SERIAL PRIMARY KEY,
		title VARCHAR(64) NOT NULL,
		description VARCHAR(500) NOT NULL
	);`

	_, err := s.DB.Exec(query)
	return err
}