package usecase

import (
	"context"
	"database/sql"
	"time"

	pb "github.com/noirbyss/worktrition-app/gen/gamification-service"
	"go.uber.org/zap"
)

type GamificationUseCase struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewGamificationUseCase(db *sql.DB, logger *zap.SugaredLogger) *GamificationUseCase {
	return &GamificationUseCase{db: db, logger: logger}
}

type RewardResult struct {
	Character *pb.Character
	GainedXP  int32
	LeveledUp bool
}

func getNextLevelXP(level int32) int32 {
	if level < 20 {
		return level * 200
	}
	return 4000
}

func (uc *GamificationUseCase) ensureCharacterExists(ctx context.Context, userID string) error {
	query := `INSERT INTO characters (user_id) VALUES ($1) ON CONFLICT (user_id) DO NOTHING`
	_, err := uc.db.ExecContext(ctx, query, userID)
	return err
}

func (uc *GamificationUseCase) GetCharacter(ctx context.Context, userID string) (*pb.Character, error) {
	if err := uc.ensureCharacterExists(ctx, userID); err != nil {
		return nil, err
	}

	query := `
		SELECT user_id, level, current_xp, hp, strength, endurance, discipline, balance, current_streak
		FROM characters WHERE user_id = $1
	`
	c := &pb.Character{}
	err := uc.db.QueryRowContext(ctx, query, userID).Scan(
		&c.UserId, &c.Level, &c.CurrentXp, &c.Hp,
		&c.Strength, &c.Endurance, &c.Discipline, &c.Balance, &c.CurrentStreak,
	)
	if err != nil {
		return nil, err
	}
	c.NextLevelXp = getNextLevelXP(c.Level)
	return c, nil
}

func (uc *GamificationUseCase) ApplyReward(ctx context.Context, userID string, xpGain, strGain, endGain, discGain, balGain int32) (*RewardResult, error) {
	if err := uc.ensureCharacterExists(ctx, userID); err != nil {
		return nil, err
	}

	tx, err := uc.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	querySelect := `
		SELECT level, current_xp, hp, strength, endurance, discipline, balance, current_streak, last_active_date
		FROM characters WHERE user_id = $1 FOR NO KEY UPDATE
	`
	var (
		level, currentXP, strength, endurance, discipline, balance, currentStreak int32
		hp                                                                        float64
		lastActive                                                                sql.NullTime
	)

	err = tx.QueryRowContext(ctx, querySelect, userID).Scan(
		&level, &currentXP, &hp, &strength, &endurance, &discipline, &balance, &currentStreak, &lastActive,
	)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	if !lastActive.Valid {
		currentStreak = 1
		hp = 6.0
	} else {
		last := time.Date(lastActive.Time.Year(), lastActive.Time.Month(), lastActive.Time.Day(), 0, 0, 0, 0, time.UTC)
		daysDiff := int(today.Sub(last).Hours() / 24)

		if daysDiff == 1 {
			currentStreak++
			hp += 0.5
		} else if daysDiff > 1 {
			currentStreak = 1
			hp -= float64(daysDiff-1) * 1.0
			hp += 0.5
		}
	}

	if hp > 6.0 {
		hp = 6.0
	}
	if hp < 0.0 {
		hp = 0.0
	}

	currentXP += xpGain
	leveledUp := false
	for currentXP >= getNextLevelXP(level) {
		currentXP -= getNextLevelXP(level)
		level++
		leveledUp = true
	}

	strength += strGain
	endurance += endGain
	discipline += discGain
	balance += balGain

	queryUpdate := `
		UPDATE characters SET 
			level = $1, current_xp = $2, hp = $3, strength = $4, endurance = $5, 
			discipline = $6, balance = $7, current_streak = $8, last_active_date = $9, updated_at = NOW()
		WHERE user_id = $10
	`
	_, err = tx.ExecContext(ctx, queryUpdate,
		level, currentXP, hp, strength, endurance, discipline, balance, currentStreak, today, userID,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &RewardResult{
		Character: &pb.Character{
			UserId: userID, Level: level, CurrentXp: currentXP, NextLevelXp: getNextLevelXP(level),
			Hp: hp, Strength: strength, Endurance: endurance, Discipline: discipline, Balance: balance, CurrentStreak: currentStreak,
		},
		GainedXP: xpGain, LeveledUp: leveledUp,
	}, nil
}
