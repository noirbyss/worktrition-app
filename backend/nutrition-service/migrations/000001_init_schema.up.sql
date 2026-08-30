CREATE TABLE plan_templates (
    id            SERIAL PRIMARY KEY,
    user_id       TEXT NOT NULL,
    generation_id TEXT NOT NULL UNIQUE,
    calories      DOUBLE PRECISION NOT NULL,
    protein       DOUBLE PRECISION NOT NULL,
    fat           DOUBLE PRECISION NOT NULL,
    carb          DOUBLE PRECISION NOT NULL,
    water_goal    INTEGER NOT NULL,
    is_active     BOOLEAN NOT NULL DEFAULT true,
    activated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE meal_templates (
    id          SERIAL PRIMARY KEY,
    plan_id     INTEGER NOT NULL REFERENCES plan_templates(id) ON DELETE CASCADE,
    day_of_week SMALLINT NOT NULL,
    calories    DOUBLE PRECISION NOT NULL,
    protein     DOUBLE PRECISION NOT NULL,
    fat         DOUBLE PRECISION NOT NULL,
    carb        DOUBLE PRECISION NOT NULL
);

CREATE TABLE meal_items (
    id                SERIAL PRIMARY KEY,
    meal_template_id  INTEGER NOT NULL REFERENCES meal_templates(id) ON DELETE CASCADE,
    name              TEXT NOT NULL,
    recipe            TEXT NOT NULL,
    calories          DOUBLE PRECISION NOT NULL,
    protein           DOUBLE PRECISION NOT NULL,
    fat               DOUBLE PRECISION NOT NULL,
    carb              DOUBLE PRECISION NOT NULL
);

CREATE TABLE meal_completions (
    id             SERIAL PRIMARY KEY,
    meal_item_id   INTEGER NOT NULL REFERENCES meal_items(id) ON DELETE CASCADE,
    completed_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE water_completions (
    id            SERIAL PRIMARY KEY,
    plan_id       INTEGER NOT NULL REFERENCES plan_templates(id) ON DELETE CASCADE,
    amount_ml     INTEGER NOT NULL,
    completed_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);