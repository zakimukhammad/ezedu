package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ezedu/backend/internal/model"
)

// LeaderboardStore manages weekly leaderboard entries and rankings.
type LeaderboardStore struct {
	db *sql.DB
}

func NewLeaderboardStore(db *sql.DB) *LeaderboardStore {
	return &LeaderboardStore{db: db}
}

// CurrentWeekStart returns the date string (YYYY-MM-DD) for Monday of the current week.
func CurrentWeekStart() string {
	now := time.Now()
	// Calculate offset to Monday (1 = Monday, ..., 7 = Sunday)
	weekday := int(now.Weekday())
	if weekday == 0 { // Sunday
		weekday = 7
	}
	monday := now.AddDate(0, 0, -(weekday - 1))
	return monday.Format("2006-01-02")
}

var animalAdjectives = []string{
	"Singa Pemberani", "Elang Cerdas", "Lumba-lumba Cepat", "Kancil Cerdik",
	"Harimau Tangkas", "Serigala Gigih", "Panda Ceria", "Gajah Hebat",
	"Rusa Lincah", "Tupai Rajin", "Kucing Pintar", "Zebra Kreatif",
}

// GenerateAnonymousName creates an anonymous handle based on child ID.
func GenerateAnonymousName(childID int64) string {
	idx := int((childID * 7 + 3) % int64(len(animalAdjectives)))
	return animalAdjectives[idx]
}

// GetWeeklyLeaderboard retrieves top entries for a given week and marks the requesting child.
func (s *LeaderboardStore) GetWeeklyLeaderboard(targetChildID int64, weekStart string, limit int) ([]model.LeaderboardEntry, error) {
	if limit <= 0 {
		limit = 20
	}

	rows, err := s.db.Query(
		`SELECT child_id, weekly_xp, display_name, avatar_id
		 FROM leaderboard_entries
		 WHERE week_start = ?
		 ORDER BY weekly_xp DESC, child_id ASC
		 LIMIT ?`,
		weekStart, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("get weekly leaderboard: %w", err)
	}
	defer rows.Close()

	var entries []model.LeaderboardEntry
	rank := 1
	foundMe := false

	for rows.Next() {
		var childID int64
		var e model.LeaderboardEntry
		if err := rows.Scan(&childID, &e.WeeklyXP, &e.DisplayName, &e.AvatarID); err != nil {
			return nil, fmt.Errorf("scan leaderboard entry: %w", err)
		}
		e.Rank = rank
		if childID == targetChildID {
			e.IsMe = true
			foundMe = true
		}
		entries = append(entries, e)
		rank++
	}

	// If requesting child opted in but not in top limit, append their rank at the end if present
	if targetChildID > 0 && !foundMe {
		myRank, myEntry, err := s.GetChildRankAndEntry(targetChildID, weekStart)
		if err == nil && myEntry != nil {
			myEntry.Rank = myRank
			myEntry.IsMe = true
			entries = append(entries, *myEntry)
		}
	}

	return entries, nil
}

// GetChildRankAndEntry returns a child's rank and entry for a given week.
func (s *LeaderboardStore) GetChildRankAndEntry(childID int64, weekStart string) (int, *model.LeaderboardEntry, error) {
	var myXP int
	var displayName string
	var avatarID int

	err := s.db.QueryRow(
		`SELECT weekly_xp, display_name, avatar_id FROM leaderboard_entries WHERE child_id = ? AND week_start = ?`,
		childID, weekStart,
	).Scan(&myXP, &displayName, &avatarID)

	if err == sql.ErrNoRows {
		return 0, nil, nil
	}
	if err != nil {
		return 0, nil, err
	}

	var higherCount int
	err = s.db.QueryRow(
		`SELECT COUNT(*) FROM leaderboard_entries WHERE week_start = ? AND (weekly_xp > ? OR (weekly_xp = ? AND child_id < ?))`,
		weekStart, myXP, myXP, childID,
	).Scan(&higherCount)
	if err != nil {
		return 0, nil, err
	}

	rank := higherCount + 1
	return rank, &model.LeaderboardEntry{
		Rank:        rank,
		DisplayName: displayName,
		AvatarID:    avatarID,
		WeeklyXP:    myXP,
		IsMe:        true,
	}, nil
}

// UpsertWeeklyXP adds XP gained to a child's weekly leaderboard entry.
func (s *LeaderboardStore) UpsertWeeklyXP(childID int64, weekStart string, xpGained int, avatarID int) error {
	displayName := GenerateAnonymousName(childID)

	_, err := s.db.Exec(
		`INSERT INTO leaderboard_entries (child_id, week_start, weekly_xp, display_name, avatar_id)
		 VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT(child_id, week_start) DO UPDATE SET
		   weekly_xp = weekly_xp + excluded.weekly_xp,
		   avatar_id = excluded.avatar_id`,
		childID, weekStart, xpGained, displayName, avatarID,
	)
	return err
}

// RemoveChildEntry removes a child from the current week's leaderboard (e.g. on opt-out).
func (s *LeaderboardStore) RemoveChildEntry(childID int64, weekStart string) error {
	_, err := s.db.Exec(`DELETE FROM leaderboard_entries WHERE child_id = ? AND week_start = ?`, childID, weekStart)
	return err
}
