-- Supplement protocols table
CREATE TABLE IF NOT EXISTS supplement_protocols (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    workout_program_id TEXT,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    purpose TEXT NOT NULL,
    safety_notes TEXT, -- JSON array
    contraindications TEXT, -- JSON array
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Supplement timings table
CREATE TABLE IF NOT EXISTS supplement_timings (
    id TEXT PRIMARY KEY,
    protocol_id TEXT NOT NULL,
    timing_type TEXT NOT NULL, -- pre_workout, during_workout, post_workout, daily
    supplement_name TEXT NOT NULL,
    dosage TEXT NOT NULL,
    unit TEXT NOT NULL,
    timing_minutes INTEGER NOT NULL,
    instructions TEXT,
    benefits TEXT, -- JSON array
    side_effects TEXT, -- JSON array
    max_daily_dose TEXT,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (protocol_id) REFERENCES supplement_protocols(id) ON DELETE CASCADE
);

-- Drug interactions table
CREATE TABLE IF NOT EXISTS supplement_drug_interactions (
    id TEXT PRIMARY KEY,
    protocol_id TEXT NOT NULL,
    drug_name TEXT NOT NULL,
    interaction TEXT NOT NULL,
    severity TEXT NOT NULL,
    management TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (protocol_id) REFERENCES supplement_protocols(id) ON DELETE CASCADE
);

-- Supplement safety data table
CREATE TABLE IF NOT EXISTS supplement_safety_data (
    id TEXT PRIMARY KEY,
    supplement_name TEXT UNIQUE NOT NULL,
    max_daily_dose TEXT,
    contraindications TEXT, -- JSON array
    side_effects TEXT, -- JSON array
    drug_interactions TEXT, -- JSON array
    timing_restrictions TEXT,
    special_populations TEXT, -- JSON array
    references TEXT, -- JSON array
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_supplement_protocols_user_id ON supplement_protocols(user_id);
CREATE INDEX IF NOT EXISTS idx_supplement_protocols_category ON supplement_protocols(category);
CREATE INDEX IF NOT EXISTS idx_supplement_timings_protocol_id ON supplement_timings(protocol_id);
CREATE INDEX IF NOT EXISTS idx_supplement_timings_type ON supplement_timings(timing_type);
CREATE INDEX IF NOT EXISTS idx_supplement_interactions_protocol_id ON supplement_drug_interactions(protocol_id);
CREATE INDEX IF NOT EXISTS idx_supplement_safety_name ON supplement_safety_data(supplement_name);

-- Insert basic supplement safety data
INSERT OR REPLACE INTO supplement_safety_data (
    id, supplement_name, max_daily_dose, contraindications, side_effects, 
    drug_interactions, timing_restrictions, special_populations, references, 
    created_at, updated_at
) VALUES 
(
    'caffeine_001',
    'caffeine',
    '400mg',
    '["pregnancy", "heart_conditions", "anxiety_disorders", "insomnia"]',
    '["jitters", "insomnia", "increased_heart_rate", "anxiety", "digestive_upset"]',
    '["blood_thinners", "stimulant_medications", "beta_blockers"]',
    'Avoid 6 hours before bedtime',
    '["pregnant_women", "children", "elderly", "cardiac_patients"]',
    '["https://pubmed.ncbi.nlm.nih.gov/34633073/"]',
    datetime('now'),
    datetime('now')
),
(
    'creatine_001',
    'creatine',
    '10g',
    '["kidney_disease", "liver_disease", "dehydration"]',
    '["water_retention", "digestive_issues", "muscle_cramps"]',
    '["nephrotoxic_drugs", "diuretics"]',
    'Take with plenty of water',
    '["kidney_patients", "liver_patients", "children"]',
    '["https://pubmed.ncbi.nlm.nih.gov/34447738/"]',
    datetime('now'),
    datetime('now')
),
(
    'whey_protein_001',
    'whey_protein',
    '50g per serving',
    '["milk_allergy", "lactose_intolerance"]',
    '["digestive_upset", "bloating", "nausea"]',
    '["none_significant"]',
    'Best within 30 minutes post-workout',
    '["lactose_intolerant", "milk_allergic"]',
    '["https://pubmed.ncbi.nlm.nih.gov/36685500/"]',
    datetime('now'),
    datetime('now')
),
(
    'beta_alanine_001',
    'beta_alanine',
    '5g',
    '["pregnancy", "breastfeeding"]',
    '["tingling_sensation", "flushing"]',
    '["none_significant"]',
    'Split doses to reduce tingling',
    '["pregnant_women", "breastfeeding_women"]',
    '["https://pubmed.ncbi.nlm.nih.gov/34633073/"]',
    datetime('now'),
    datetime('now')
),
(
    'omega3_001',
    'omega3',
    '3g',
    '["fish_allergy", "bleeding_disorders"]',
    '["fishy_aftertaste", "digestive_upset", "bleeding"]',
    '["blood_thinners", "antiplatelet_drugs"]',
    'Take with meals to reduce aftertaste',
    '["bleeding_disorder_patients", "surgery_patients"]',
    '["https://pubmed.ncbi.nlm.nih.gov/36834422/"]',
    datetime('now'),
    datetime('now')
);