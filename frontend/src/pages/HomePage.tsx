import { AppFrame } from '../components/app/AppFrame'
import { NotchBar } from '../components/app/PlaceholderUi'

const sidebarProfile = {
  badge: '7',
  meta: '820 / 1000 XP',
  name: 'Артём Ковалёв',
}

const characterStats = [
  { active: 14, label: 'Сила', value: '35', variant: 'strength' as const },
  { active: 11, label: 'Выносливость', value: '28', variant: 'endurance' as const },
  { active: 17, label: 'Дисциплина', value: '42', variant: 'discipline' as const },
  { active: 12, label: 'Баланс', value: '31', variant: 'balance' as const },
]

const achievements = [
  { icon: '✓', name: 'Первая тренировка', progress: 'получено', status: 'unlocked' as const },
  { icon: '✓', name: 'Неделя без пропуска воды', progress: 'получено', status: 'unlocked' as const },
  { icon: '5/7', name: '7 дней подряд', progress: 'осталось 2 дня', status: 'locked' as const },
  { icon: '7/10', name: '10 дней питания', progress: 'осталось 3 дня', status: 'locked' as const },
  { icon: '14/20', name: '20 тренировок', progress: 'осталось 6', status: 'locked' as const },
]

export function HomePage() {
  return (
    <AppFrame
      description="Тело развивается в реальной жизни — персонаж растёт в приложении."
      eyebrow="Главная · Экран 01"
      sidebarProfile={sidebarProfile}
      title="Персонаж"
    >
      <section className="card frame placeholder-section">
        <div className="hero">
          <div className="ring-wrap">
            <svg viewBox="0 0 118 118">
              <circle cx="59" cy="59" fill="none" r="50" stroke="rgba(255,255,255,.08)" strokeWidth="8" />
              <circle
                cx="59"
                cy="59"
                fill="none"
                r="50"
                stroke="url(#characterLevelGradient)"
                strokeDasharray="314.2"
                strokeDashoffset="56.5"
                strokeLinecap="round"
                strokeWidth="8"
              />
              <defs>
                <linearGradient id="characterLevelGradient" x1="0" x2="1" y1="0" y2="1">
                  <stop offset="0" stopColor="#FFD23F" />
                  <stop offset="1" stopColor="#FF7A1A" />
                </linearGradient>
              </defs>
            </svg>
            <div className="ring-center">
              <div className="ring-lv">УРОВЕНЬ</div>
              <div className="ring-num">7</div>
            </div>
          </div>

          <div className="hero-copy">
            <h2 className="hero-name">Артём Ковалёв</h2>
            <div className="hero-title">Новичок → Воин · 3 тренировки до повышения класса</div>
            <div className="bar-row hero-progress-copy">
              <span className="bar-label">Опыт</span>
              <span className="bar-value">820 / 1000 XP</span>
            </div>
            <NotchBar active={28} large total={34} />
          </div>
        </div>
      </section>

      <section className="grid g-2 placeholder-grid">
        <div className="card">
          <div className="card-title">Характеристики</div>

          {characterStats.map((item) => (
            <div className="stat-row" key={item.label}>
              <span className="stat-label">
                <span className={`stat-dot stat-dot--${item.variant}`} />
                {item.label}
              </span>
              <NotchBar active={item.active} total={20} variant={item.variant} />
              <span className="stat-value">{item.value}</span>
            </div>
          ))}

          <hr className="rule" />

          <div className="legend-list">
            <div className="legend-row">
              <span className="k">Сила</span>
              <span className="v">растёт за силовые тренировки</span>
            </div>
            <div className="legend-row">
              <span className="k">Выносл.</span>
              <span className="v">растёт за кардио и регулярность</span>
            </div>
            <div className="legend-row">
              <span className="k">Дисципл.</span>
              <span className="v">растёт за задачи без пропусков</span>
            </div>
            <div className="legend-row">
              <span className="k">Баланс</span>
              <span className="v">растёт за питание и воду</span>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="card-title">План на сегодня</div>

          <div className="check-item">
            <span className="checkbox" />
            <span className="check-label">Тренировка — кардио, дом</span>
            <span className="pill pending">не выполнена</span>
          </div>
          <div className="check-item">
            <span className="checkbox done" />
            <span className="check-label">План питания — завтрак и обед</span>
            <span className="pill done">выполнено</span>
          </div>
          <div className="check-item">
            <span className="checkbox" />
            <span className="check-label">Вода — 2.5 л</span>
            <span className="pill pending">1.4 / 2.5 л</span>
          </div>

          <hr className="rule" />

          <div className="bar-row">
            <span className="bar-label">Прогресс дня</span>
            <span className="bar-value">3 из 5 задач</span>
          </div>
          <NotchBar active={14} total={24} />

          <div className="qs qs--compact">
            <QuickStat label="Текущий вес" value="88 кг" />
            <QuickStat label="Серия дней подряд" value="5" />
          </div>
        </div>
      </section>

      <section className="card placeholder-section">
        <div className="card-title">Достижения</div>
        <div className="ach-grid">
          {achievements.map((achievement) => (
            <div className={`ach ${achievement.status}`} key={achievement.name}>
              <div className="ach-icon">{achievement.icon}</div>
              <div className="ach-name">{achievement.name}</div>
              <div className="ach-progress">{achievement.progress}</div>
            </div>
          ))}
        </div>
      </section>

      <footer className="foot">ФОРМА — пример интерфейса для хакатона · экран index</footer>
    </AppFrame>
  )
}

function QuickStat({ label, value }: { label: string; value: string }) {
  return (
    <div className="qs-item">
      <div className="qs-num">{value}</div>
      <div className="qs-label">{label}</div>
    </div>
  )
}
