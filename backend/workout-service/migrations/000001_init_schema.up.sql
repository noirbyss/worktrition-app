CREATE TABLE plan_templates (
    id            SERIAL PRIMARY KEY,
    user_id       TEXT NOT NULL,
    generation_id TEXT NOT NULL UNIQUE,
    is_active     BOOLEAN NOT NULL DEFAULT true,
    activated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE training_templates (
    id          SERIAL PRIMARY KEY,
    plan_id     INTEGER NOT NULL REFERENCES plan_templates(id) ON DELETE CASCADE,
    day_of_week SMALLINT NOT NULL,
    type        TEXT NOT NULL,
    UNIQUE (plan_id, day_of_week)
);

CREATE TABLE exercises (
    id                   SERIAL PRIMARY KEY,
    training_template_id INTEGER NOT NULL REFERENCES training_templates(id) ON DELETE CASCADE,
    position             SMALLINT NOT NULL,
    name                 TEXT NOT NULL
);

CREATE TABLE training_completions (
    id                   SERIAL PRIMARY KEY,
    training_template_id INTEGER NOT NULL REFERENCES training_templates(id) ON DELETE CASCADE,
    duration_seconds     INTEGER NOT NULL,
    completed_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_plan_templates_user_active ON plan_templates (user_id) WHERE is_active;
CREATE INDEX idx_training_templates_plan ON training_templates (plan_id);
CREATE INDEX idx_exercises_training ON exercises (training_template_id);
CREATE INDEX idx_training_completions_template ON training_completions (training_template_id);
