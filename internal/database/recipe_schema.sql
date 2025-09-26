-- Recipe database schema
CREATE TABLE IF NOT EXISTS recipes (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    cuisine TEXT NOT NULL,
    category TEXT NOT NULL,
    ingredients TEXT NOT NULL, -- JSON array
    instructions TEXT NOT NULL, -- JSON array
    nutrition TEXT NOT NULL, -- JSON object
    macros TEXT NOT NULL, -- JSON object
    prep_time INTEGER NOT NULL,
    cook_time INTEGER NOT NULL,
    servings INTEGER NOT NULL,
    difficulty TEXT NOT NULL,
    alternatives TEXT, -- JSON array
    skills_tips TEXT, -- JSON array
    enhancement_tips TEXT, -- JSON array
    health_benefits TEXT, -- JSON array
    tags TEXT, -- JSON array
    image_url TEXT,
    video_url TEXT,
    rating REAL DEFAULT 0,
    review_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_recipes_cuisine ON recipes(cuisine);
CREATE INDEX IF NOT EXISTS idx_recipes_category ON recipes(category);
CREATE INDEX IF NOT EXISTS idx_recipes_difficulty ON recipes(difficulty);
CREATE INDEX IF NOT EXISTS idx_recipes_created_at ON recipes(created_at);
CREATE INDEX IF NOT EXISTS idx_recipes_rating ON recipes(rating);

-- Recipe ratings table
CREATE TABLE IF NOT EXISTS recipe_ratings (
    id TEXT PRIMARY KEY,
    recipe_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    review TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
    UNIQUE(recipe_id, user_id)
);

-- Recipe collections table
CREATE TABLE IF NOT EXISTS recipe_collections (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    user_id TEXT NOT NULL,
    recipe_ids TEXT NOT NULL, -- JSON array
    is_public BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Meal plans table
CREATE TABLE IF NOT EXISTS meal_plans (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    meals TEXT NOT NULL, -- JSON array
    nutrition_goals TEXT NOT NULL, -- JSON object
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Shopping lists table
CREATE TABLE IF NOT EXISTS shopping_lists (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    items TEXT NOT NULL, -- JSON array
    recipe_ids TEXT, -- JSON array
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Recipe search view for full-text search
CREATE VIEW IF NOT EXISTS recipe_search AS
SELECT 
    id,
    name,
    cuisine,
    category,
    ingredients,
    health_benefits,
    name || ' ' || ingredients || ' ' || health_benefits AS search_text
FROM recipes;