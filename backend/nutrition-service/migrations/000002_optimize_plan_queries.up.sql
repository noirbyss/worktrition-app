WITH ranked_active_plans AS (
    SELECT
        id,
        ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY activated_at DESC, id DESC) AS row_num
    FROM plan_templates
    WHERE is_active = true
)
UPDATE plan_templates AS pt
SET is_active = false
FROM ranked_active_plans AS rap
WHERE pt.id = rap.id
  AND rap.row_num > 1;

CREATE INDEX idx_plan_templates_user ON plan_templates (user_id);
CREATE UNIQUE INDEX idx_plan_templates_user_active_unique ON plan_templates (user_id) WHERE is_active;
CREATE INDEX idx_meal_templates_plan_day ON meal_templates (plan_id, day_of_week);
CREATE INDEX idx_meal_items_meal_template ON meal_items (meal_template_id);
CREATE INDEX idx_meal_completions_meal_item_completed_at ON meal_completions (meal_item_id, completed_at);
CREATE INDEX idx_water_completions_plan_completed_at ON water_completions (plan_id, completed_at);
