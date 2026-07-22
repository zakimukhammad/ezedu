package store

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// NewDB creates a new SQLite database connection with optimized settings.
func NewDB(path string) (*sql.DB, error) {
	dsn := fmt.Sprintf("file:%s?_pragma=journal_mode%%3DWAL&_pragma=synchronous%%3DNORMAL&_pragma=cache_size%%3D-20000&_pragma=busy_timeout%%3D5000&_pragma=foreign_keys%%3DON", path)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// Connection pool settings for SQLite
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(30 * time.Minute)

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}

// Migrate runs all database migrations.
func Migrate(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS accounts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			parent_name TEXT NOT NULL,
			parent_pin TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS children (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			account_id INTEGER NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			birth_year INTEGER NOT NULL,
			age_group TEXT NOT NULL CHECK(age_group IN ('toddlers', 'explorers', 'builders', 'challengers')),
			avatar_id INTEGER DEFAULT 1,
			xp_total INTEGER DEFAULT 0,
			current_level INTEGER DEFAULT 1,
			streak_days INTEGER DEFAULT 0,
			last_active DATE,
			daily_limit_min INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			slug TEXT UNIQUE NOT NULL,
			name TEXT NOT NULL,
			description TEXT DEFAULT '',
			icon TEXT DEFAULT '',
			color TEXT DEFAULT '#6366f1',
			sort_order INTEGER DEFAULT 0
		)`,

		`CREATE TABLE IF NOT EXISTS lessons (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			category_id INTEGER NOT NULL REFERENCES categories(id),
			age_group TEXT NOT NULL CHECK(age_group IN ('toddlers', 'explorers', 'builders', 'challengers')),
			level INTEGER NOT NULL CHECK(level BETWEEN 1 AND 5),
			sort_order INTEGER DEFAULT 0,
			title TEXT NOT NULL,
			description TEXT DEFAULT '',
			content_json TEXT DEFAULT '{}',
			estimated_minutes INTEGER DEFAULT 10,
			xp_reward INTEGER DEFAULT 10,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS activities (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			lesson_id INTEGER NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
			type TEXT NOT NULL CHECK(type IN ('multiple_choice', 'drag_drop', 'fill_blank', 'matching', 'sequencing', 'drawing', 'pixel_art', 'block_code', 'timed')),
			sort_order INTEGER DEFAULT 0,
			question_json TEXT NOT NULL DEFAULT '{}',
			max_score INTEGER DEFAULT 10
		)`,

		`CREATE TABLE IF NOT EXISTS child_progress (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			child_id INTEGER NOT NULL REFERENCES children(id) ON DELETE CASCADE,
			lesson_id INTEGER NOT NULL REFERENCES lessons(id),
			status TEXT NOT NULL DEFAULT 'not_started' CHECK(status IN ('not_started', 'in_progress', 'completed')),
			score INTEGER DEFAULT 0,
			max_possible INTEGER DEFAULT 0,
			attempts INTEGER DEFAULT 0,
			started_at DATETIME,
			completed_at DATETIME,
			time_spent_sec INTEGER DEFAULT 0,
			UNIQUE(child_id, lesson_id)
		)`,

		`CREATE TABLE IF NOT EXISTS activity_results (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			child_id INTEGER NOT NULL REFERENCES children(id) ON DELETE CASCADE,
			activity_id INTEGER NOT NULL REFERENCES activities(id),
			answer_json TEXT DEFAULT '{}',
			is_correct BOOLEAN DEFAULT 0,
			score INTEGER DEFAULT 0,
			attempt_number INTEGER DEFAULT 1,
			answered_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS badges (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			slug TEXT UNIQUE NOT NULL,
			name TEXT NOT NULL,
			description TEXT DEFAULT '',
			icon TEXT DEFAULT 'badge-default',
			category_id INTEGER REFERENCES categories(id),
			criteria_json TEXT DEFAULT '{}'
		)`,

		`CREATE TABLE IF NOT EXISTS child_badges (
			child_id INTEGER NOT NULL REFERENCES children(id) ON DELETE CASCADE,
			badge_id INTEGER NOT NULL REFERENCES badges(id),
			earned_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (child_id, badge_id)
		)`,

		`CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			account_id INTEGER NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			expires_at DATETIME NOT NULL,
			ip_address TEXT DEFAULT ''
		)`,

		`CREATE TABLE IF NOT EXISTS daily_challenges (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			challenge_date DATE NOT NULL,
			age_group TEXT NOT NULL CHECK(age_group IN ('toddlers', 'explorers', 'builders', 'challengers')),
			activity_type TEXT NOT NULL,
			question_json TEXT NOT NULL DEFAULT '{}',
			max_score INTEGER DEFAULT 10,
			UNIQUE(challenge_date, age_group)
		)`,

		`CREATE TABLE IF NOT EXISTS daily_challenge_results (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			child_id INTEGER NOT NULL REFERENCES children(id) ON DELETE CASCADE,
			challenge_id INTEGER NOT NULL REFERENCES daily_challenges(id),
			score INTEGER DEFAULT 0,
			completed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(child_id, challenge_id)
		)`,

		// Indexes
		`CREATE INDEX IF NOT EXISTS idx_children_account ON children(account_id)`,
		`CREATE INDEX IF NOT EXISTS idx_lessons_category_age ON lessons(category_id, age_group, level)`,
		`CREATE INDEX IF NOT EXISTS idx_progress_child ON child_progress(child_id, lesson_id)`,
		`CREATE INDEX IF NOT EXISTS idx_results_child ON activity_results(child_id, activity_id)`,
		`CREATE INDEX IF NOT EXISTS idx_badges_child ON child_badges(child_id)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_expires ON sessions(expires_at)`,
	}

	for i, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			return fmt.Errorf("migration %d failed: %w", i, err)
		}
	}

	// Dynamic schema upgrade for pre-existing SQLite databases missing 'toddlers' in CHECK constraint
	var childrenDDL string
	err := db.QueryRow(`SELECT sql FROM sqlite_master WHERE type='table' AND name='children'`).Scan(&childrenDDL)
	if err == nil && !strings.Contains(childrenDDL, "toddlers") {
		_, _ = db.Exec(`PRAGMA foreign_keys=OFF;`)
		_, _ = db.Exec(`CREATE TABLE children_new (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			account_id INTEGER NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			birth_year INTEGER NOT NULL,
			age_group TEXT NOT NULL CHECK(age_group IN ('toddlers', 'explorers', 'builders', 'challengers')),
			avatar_id INTEGER DEFAULT 1,
			xp_total INTEGER DEFAULT 0,
			current_level INTEGER DEFAULT 1,
			streak_days INTEGER DEFAULT 0,
			last_active DATE,
			daily_limit_min INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`)
		_, _ = db.Exec(`INSERT INTO children_new SELECT * FROM children;`)
		_, _ = db.Exec(`DROP TABLE children;`)
		_, _ = db.Exec(`ALTER TABLE children_new RENAME TO children;`)
		_, _ = db.Exec(`PRAGMA foreign_keys=ON;`)
	}

	var lessonsDDL string
	err = db.QueryRow(`SELECT sql FROM sqlite_master WHERE type='table' AND name='lessons'`).Scan(&lessonsDDL)
	if err == nil && !strings.Contains(lessonsDDL, "toddlers") {
		_, _ = db.Exec(`PRAGMA foreign_keys=OFF;`)
		_, _ = db.Exec(`CREATE TABLE lessons_new (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			category_id INTEGER NOT NULL REFERENCES categories(id),
			age_group TEXT NOT NULL CHECK(age_group IN ('toddlers', 'explorers', 'builders', 'challengers')),
			level INTEGER NOT NULL CHECK(level BETWEEN 1 AND 5),
			sort_order INTEGER DEFAULT 0,
			title TEXT NOT NULL,
			description TEXT DEFAULT '',
			content_json TEXT DEFAULT '{}',
			estimated_minutes INTEGER DEFAULT 10,
			xp_reward INTEGER DEFAULT 10,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`)
		_, _ = db.Exec(`INSERT INTO lessons_new SELECT * FROM lessons;`)
		_, _ = db.Exec(`DROP TABLE lessons;`)
		_, _ = db.Exec(`ALTER TABLE lessons_new RENAME TO lessons;`)
		_, _ = db.Exec(`PRAGMA foreign_keys=ON;`)
	}

	var activitiesDDL string
	err = db.QueryRow(`SELECT sql FROM sqlite_master WHERE type='table' AND name='activities'`).Scan(&activitiesDDL)
	if err == nil && !strings.Contains(activitiesDDL, "pixel_art") {
		_, _ = db.Exec(`PRAGMA foreign_keys=OFF;`)
		_, _ = db.Exec(`CREATE TABLE activities_new (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			lesson_id INTEGER NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
			type TEXT NOT NULL CHECK(type IN ('multiple_choice', 'drag_drop', 'fill_blank', 'matching', 'sequencing', 'drawing', 'pixel_art', 'block_code', 'timed')),
			sort_order INTEGER DEFAULT 0,
			question_json TEXT NOT NULL DEFAULT '{}',
			max_score INTEGER DEFAULT 10
		);`)
		_, _ = db.Exec(`INSERT INTO activities_new SELECT * FROM activities;`)
		_, _ = db.Exec(`DROP TABLE activities;`)
		_, _ = db.Exec(`ALTER TABLE activities_new RENAME TO activities;`)
		_, _ = db.Exec(`PRAGMA foreign_keys=ON;`)
	}

	return nil
}

// SeedDefaults inserts default categories and badges if they don't exist.
func SeedDefaults(db *sql.DB) error {
	// Seed categories
	categories := []struct {
		slug, name, description, icon, color string
		sortOrder                            int
	}{
		{"toddlers", "Mengenal Dunia", "Eksplorasi sensorik, bentuk, warna, dan suara hewan untuk balita", "sparkles", "#f59e0b", 0},
		{"math", "Matematika", "Belajar angka, hitung, dan geometri", "calculator", "#f59e0b", 1},
		{"science", "Sains & Alam", "Jelajahi dunia sains dan alam sekitar", "microscope", "#10b981", 2},
		{"coding", "Koding & Logika", "Belajar dasar pemrograman dan berpikir logis", "code", "#6366f1", 3},
		{"language", "Bahasa & Literasi", "Membaca, menulis, dan kosakata", "book-open", "#ec4899", 4},
		{"logic", "Logika & Teka-teki", "Asah otak dengan puzzle dan teka-teki seru", "puzzle", "#8b5cf6", 5},
		{"art", "Seni & Kreativitas", "Gambar, warna, musik, dan kreasi digital", "palette", "#f97316", 6},
	}

	for _, c := range categories {
		_, err := db.Exec(
			`INSERT OR IGNORE INTO categories (slug, name, description, icon, color, sort_order)
			 VALUES (?, ?, ?, ?, ?, ?)`,
			c.slug, c.name, c.description, c.icon, c.color, c.sortOrder,
		)
		if err != nil {
			return fmt.Errorf("seed category %s: %w", c.slug, err)
		}
	}

	// Seed initial badges (Bahasa Indonesia names - 21 total badges)
	badges := []struct {
		slug, name, description, icon string
		criteriaJSON                  string
	}{
		{"first_lesson", "Langkah Pertama", "Selesaikan pelajaran pertamamu!", "badge-first", `{"type":"lessons_completed","value":1}`},
		{"three_lessons", "Pembelajar Giat", "Selesaikan 3 pelajaran", "badge-three", `{"type":"lessons_completed","value":3}`},
		{"five_lessons", "Jagoan Belajar", "Selesaikan 5 pelajaran", "badge-five", `{"type":"lessons_completed","value":5}`},
		{"ten_lessons", "Murid Rajin", "Selesaikan 10 pelajaran", "badge-ten", `{"type":"lessons_completed","value":10}`},
		{"twenty_lessons", "Master Pelajaran", "Selesaikan 20 pelajaran", "badge-twenty", `{"type":"lessons_completed","value":20}`},
		{"math_explorer", "Penjelajah Angka", "Selesaikan 3 pelajaran Matematika", "badge-math", `{"type":"category_lessons","category":"math","value":3}`},
		{"coding_starter", "Pembuat Kode Pemula", "Selesaikan 3 pelajaran Koding", "badge-code", `{"type":"category_lessons","category":"coding","value":3}`},
		{"science_junior", "Ilmuwan Cilik", "Selesaikan 3 pelajaran Sains", "badge-science", `{"type":"category_lessons","category":"science","value":3}`},
		{"word_wizard", "Penyihir Kata", "Selesaikan 3 pelajaran Bahasa", "badge-language", `{"type":"category_lessons","category":"language","value":3}`},
		{"logic_solver", "Pemecah Teka-Teki", "Selesaikan 3 pelajaran Logika", "badge-logic", `{"type":"category_lessons","category":"logic","value":3}`},
		{"art_creator", "Seniman Kreatif", "Selesaikan 3 pelajaran Seni", "badge-art", `{"type":"category_lessons","category":"art","value":3}`},
		{"toddler_star", "Bintang Cilik", "Selesaikan 3 pelajaran Mengenal Dunia", "badge-toddler", `{"type":"category_lessons","category":"toddlers","value":3}`},
		{"streak_3", "Rajin Belajar", "Belajar 3 hari berturut-turut!", "badge-streak3", `{"type":"streak","value":3}`},
		{"streak_7", "Bintang Mingguan", "Belajar 7 hari berturut-turut!", "badge-streak7", `{"type":"streak","value":7}`},
		{"streak_14", "Pejuang Hebat", "Belajar 14 hari berturut-turut!", "badge-streak14", `{"type":"streak","value":14}`},
		{"streak_30", "Pahlawan Belajar", "Belajar 30 hari berturut-turut!", "badge-streak30", `{"type":"streak","value":30}`},
		{"perfect_quiz", "Nilai Sempurna", "Dapatkan skor sempurna di sebuah kuis!", "badge-perfect", `{"type":"perfect_score","value":1}`},
		{"xp_100", "Pengumpul Bintang", "Capai total 100 XP!", "badge-xp100", `{"type":"xp_total","value":100}`},
		{"xp_500", "Bintang Cemerlang", "Capai total 500 XP!", "badge-xp500", `{"type":"xp_total","value":500}`},
		{"level_2", "Naik Level 2", "Berhasil mencapai Level 2!", "badge-level2", `{"type":"level","value":2}`},
		{"level_5", "Puncak Prestasi", "Berhasil mencapai Level 5!", "badge-level5", `{"type":"level","value":5}`},
	}

	for _, b := range badges {
		_, err := db.Exec(
			`INSERT OR IGNORE INTO badges (slug, name, description, icon, criteria_json)
			 VALUES (?, ?, ?, ?, ?)`,
			b.slug, b.name, b.description, b.icon, b.criteriaJSON,
		)
		if err != nil {
			return fmt.Errorf("seed badge %s: %w", b.slug, err)
		}
	}

	return nil
}
