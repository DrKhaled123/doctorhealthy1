-- Comprehensive Workout Database Schema
-- Based on VIP JSON workout data structure

-- Exercises table - stores all exercise data
CREATE TABLE IF NOT EXISTS exercises (
    id TEXT PRIMARY KEY,
    name_en TEXT NOT NULL,
    name_ar TEXT,
    category TEXT NOT NULL, -- strength, cardio, flexibility, etc.
    primary_muscles TEXT NOT NULL, -- JSON array
    secondary_muscles TEXT, -- JSON array
    equipment TEXT, -- JSON array (dumbbells, barbell, bodyweight, etc.)
    difficulty TEXT NOT NULL, -- beginner, intermediate, advanced
    instructions_en TEXT NOT NULL,
    instructions_ar TEXT,
    common_mistakes_en TEXT, -- JSON array
    common_mistakes_ar TEXT, -- JSON array
    injury_risks_en TEXT, -- JSON array
    injury_risks_ar TEXT, -- JSON array
    tips_en TEXT, -- JSON array
    tips_ar TEXT, -- JSON array
    alternatives TEXT, -- JSON array of exercise IDs
    progressions TEXT, -- JSON array of exercise IDs
    evidence_links TEXT, -- JSON array of scientific links
    sets_default INTEGER DEFAULT 3,
    reps_default TEXT DEFAULT '8-12',
    rest_default TEXT DEFAULT '60 sec',
    is_gym_exercise BOOLEAN DEFAULT TRUE,
    is_home_exercise BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Workout programs table - stores structured workout programs
CREATE TABLE IF NOT EXISTS workout_programs (
    id TEXT PRIMARY KEY,
    name_en TEXT NOT NULL,
    name_ar TEXT,
    level TEXT NOT NULL, -- beginner, intermediate, advanced
    split_type TEXT NOT NULL, -- full_body, upper_lower, push_pull_legs, etc.
    description_en TEXT NOT NULL,
    description_ar TEXT,
    duration_weeks INTEGER NOT NULL DEFAULT 4,
    sessions_per_week INTEGER NOT NULL DEFAULT 3,
    goals TEXT NOT NULL, -- JSON array (weight_loss, muscle_gain, strength, etc.)
    equipment_required TEXT, -- JSON array
    target_audience TEXT, -- JSON array (time_crunched, home_based, etc.)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Weekly plans table - stores week-by-week workout plans
CREATE TABLE IF NOT EXISTS weekly_plans (
    id TEXT PRIMARY KEY,
    program_id TEXT NOT NULL,
    week_number INTEGER NOT NULL,
    progression_notes_en TEXT,
    progression_notes_ar TEXT,
    intensity_level TEXT, -- light, moderate, high
    volume_increase_percent REAL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (program_id) REFERENCES workout_programs(id) ON DELETE CASCADE,
    UNIQUE(program_id, week_number)
);

-- Daily workouts table - stores individual workout days
CREATE TABLE IF NOT EXISTS daily_workouts (
    id TEXT PRIMARY KEY,
    weekly_plan_id TEXT NOT NULL,
    day_number INTEGER NOT NULL, -- 1-7
    focus_en TEXT NOT NULL, -- Full Body A, Upper Body, Rest, etc.
    focus_ar TEXT,
    duration_minutes INTEGER DEFAULT 45,
    is_rest_day BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (weekly_plan_id) REFERENCES weekly_plans(id) ON DELETE CASCADE,
    UNIQUE(weekly_plan_id, day_number)
);

-- Workout exercises table - links exercises to daily workouts
CREATE TABLE IF NOT EXISTS workout_exercises (
    id TEXT PRIMARY KEY,
    daily_workout_id TEXT NOT NULL,
    exercise_id TEXT NOT NULL,
    exercise_order INTEGER NOT NULL,
    sets INTEGER NOT NULL DEFAULT 3,
    reps TEXT NOT NULL DEFAULT '8-12',
    rest_seconds INTEGER NOT NULL DEFAULT 60,
    weight_kg REAL,
    intensity_notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (daily_workout_id) REFERENCES daily_workouts(id) ON DELETE CASCADE,
    FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE,
    UNIQUE(daily_workout_id, exercise_order)
);

-- User workout progress table - tracks user progress
CREATE TABLE IF NOT EXISTS workout_progress (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    program_id TEXT NOT NULL,
    current_week INTEGER NOT NULL DEFAULT 1,
    completed_sessions INTEGER DEFAULT 0,
    total_sessions INTEGER NOT NULL,
    start_date DATE NOT NULL,
    last_session_date DATE,
    completion_percentage REAL DEFAULT 0,
    status TEXT DEFAULT 'active', -- active, paused, completed, abandoned
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (program_id) REFERENCES workout_programs(id) ON DELETE CASCADE,
    UNIQUE(user_id, program_id)
);

-- Session logs table - tracks individual workout sessions
CREATE TABLE IF NOT EXISTS session_logs (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    daily_workout_id TEXT NOT NULL,
    session_date DATE NOT NULL,
    duration_minutes INTEGER,
    rpe_score INTEGER, -- Rate of Perceived Exertion (1-10)
    notes TEXT,
    completed BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (daily_workout_id) REFERENCES daily_workouts(id) ON DELETE CASCADE
);

-- Exercise logs table - tracks individual exercise performance
CREATE TABLE IF NOT EXISTS exercise_logs (
    id TEXT PRIMARY KEY,
    session_log_id TEXT NOT NULL,
    exercise_id TEXT NOT NULL,
    set_number INTEGER NOT NULL,
    reps_completed INTEGER,
    weight_kg REAL,
    rest_seconds INTEGER,
    rpe_score INTEGER, -- 1-10
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_log_id) REFERENCES session_logs(id) ON DELETE CASCADE,
    FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE
);

-- Nutrition recommendations table - workout-specific nutrition
CREATE TABLE IF NOT EXISTS workout_nutrition (
    id TEXT PRIMARY KEY,
    program_id TEXT NOT NULL,
    meal_timing TEXT NOT NULL, -- pre_workout, post_workout, general
    timing_minutes INTEGER, -- minutes before/after workout
    recommendations_en TEXT NOT NULL,
    recommendations_ar TEXT,
    food_suggestions TEXT, -- JSON array
    hydration_guidelines_en TEXT,
    hydration_guidelines_ar TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (program_id) REFERENCES workout_programs(id) ON DELETE CASCADE
);

-- Supplement protocols table - workout-specific supplements
CREATE TABLE IF NOT EXISTS supplement_protocols (
    id TEXT PRIMARY KEY,
    program_id TEXT NOT NULL,
    supplement_name TEXT NOT NULL,
    dosage TEXT NOT NULL,
    timing TEXT NOT NULL, -- pre_workout, post_workout, daily
    timing_minutes INTEGER, -- minutes before/after workout
    purpose_en TEXT NOT NULL,
    purpose_ar TEXT,
    benefits_en TEXT, -- JSON array
    benefits_ar TEXT, -- JSON array
    contraindications TEXT, -- JSON array
    side_effects TEXT, -- JSON array
    interactions TEXT, -- JSON array
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (program_id) REFERENCES workout_programs(id) ON DELETE CASCADE
);

-- Warmup routines table - specific warmup routines
CREATE TABLE IF NOT EXISTS warmup_routines (
    id TEXT PRIMARY KEY,
    program_id TEXT NOT NULL,
    workout_type TEXT NOT NULL, -- upper_body, lower_body, full_body, cardio
    duration_minutes INTEGER NOT NULL DEFAULT 10,
    exercises TEXT NOT NULL, -- JSON array of warmup exercises
    instructions_en TEXT NOT NULL,
    instructions_ar TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (program_id) REFERENCES workout_programs(id) ON DELETE CASCADE
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_exercises_category ON exercises(category);
CREATE INDEX IF NOT EXISTS idx_exercises_difficulty ON exercises(difficulty);
CREATE INDEX IF NOT EXISTS idx_exercises_primary_muscles ON exercises(primary_muscles);
CREATE INDEX IF NOT EXISTS idx_exercises_equipment ON exercises(equipment);
CREATE INDEX IF NOT EXISTS idx_exercises_gym_home ON exercises(is_gym_exercise, is_home_exercise);

CREATE INDEX IF NOT EXISTS idx_workout_programs_level ON workout_programs(level);
CREATE INDEX IF NOT EXISTS idx_workout_programs_split_type ON workout_programs(split_type);
CREATE INDEX IF NOT EXISTS idx_workout_programs_goals ON workout_programs(goals);

CREATE INDEX IF NOT EXISTS idx_weekly_plans_program_id ON weekly_plans(program_id);
CREATE INDEX IF NOT EXISTS idx_weekly_plans_week_number ON weekly_plans(week_number);

CREATE INDEX IF NOT EXISTS idx_daily_workouts_weekly_plan_id ON daily_workouts(weekly_plan_id);
CREATE INDEX IF NOT EXISTS idx_daily_workouts_day_number ON daily_workouts(day_number);

CREATE INDEX IF NOT EXISTS idx_workout_exercises_daily_workout_id ON workout_exercises(daily_workout_id);
CREATE INDEX IF NOT EXISTS idx_workout_exercises_exercise_id ON workout_exercises(exercise_id);

CREATE INDEX IF NOT EXISTS idx_workout_progress_user_id ON workout_progress(user_id);
CREATE INDEX IF NOT EXISTS idx_workout_progress_program_id ON workout_progress(program_id);
CREATE INDEX IF NOT EXISTS idx_workout_progress_status ON workout_progress(status);

CREATE INDEX IF NOT EXISTS idx_session_logs_user_id ON session_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_session_logs_daily_workout_id ON session_logs(daily_workout_id);
CREATE INDEX IF NOT EXISTS idx_session_logs_session_date ON session_logs(session_date);

CREATE INDEX IF NOT EXISTS idx_exercise_logs_session_log_id ON exercise_logs(session_log_id);
CREATE INDEX IF NOT EXISTS idx_exercise_logs_exercise_id ON exercise_logs(exercise_id);

CREATE INDEX IF NOT EXISTS idx_workout_nutrition_program_id ON workout_nutrition(program_id);
CREATE INDEX IF NOT EXISTS idx_workout_nutrition_meal_timing ON workout_nutrition(meal_timing);

CREATE INDEX IF NOT EXISTS idx_supplement_protocols_program_id ON supplement_protocols(program_id);
CREATE INDEX IF NOT EXISTS idx_supplement_protocols_timing ON supplement_protocols(timing);

CREATE INDEX IF NOT EXISTS idx_warmup_routines_program_id ON warmup_routines(program_id);
CREATE INDEX IF NOT EXISTS idx_warmup_routines_workout_type ON warmup_routines(workout_type);

-- Full-text search views for exercises
CREATE VIEW IF NOT EXISTS exercise_search AS
SELECT 
    id,
    name_en,
    name_ar,
    category,
    primary_muscles,
    secondary_muscles,
    equipment,
    difficulty,
    name_en || ' ' || COALESCE(name_ar, '') || ' ' || category || ' ' || 
    primary_muscles || ' ' || COALESCE(secondary_muscles, '') || ' ' || 
    COALESCE(equipment, '') AS search_text
FROM exercises;

-- View for workout program search
CREATE VIEW IF NOT EXISTS program_search AS
SELECT 
    id,
    name_en,
    name_ar,
    level,
    split_type,
    goals,
    equipment_required,
    name_en || ' ' || COALESCE(name_ar, '') || ' ' || level || ' ' || 
    split_type || ' ' || goals || ' ' || COALESCE(equipment_required, '') AS search_text
FROM workout_programs;