import { AppFrame } from '../components/app/AppFrame'
import { NotchBar } from '../components/app/PlaceholderUi'

const sidebarProfile = {
  badge: '7',
  meta: '820 / 1000 XP',
  name: 'Артём Ковалёв',
}

const weekPlan = [
  { day: 'Пн', title: 'Силовая' },
  { day: 'Вт', title: 'Отдых', tone: 'rest' as const },
  { day: 'Ср', title: 'Кардио', tone: 'today' as const },
  { day: 'Чт', title: 'Отдых', tone: 'rest' as const },
  { day: 'Пт', title: 'Силовая' },
  { day: 'Сб', title: 'Кардио' },
  { day: 'Вс', title: 'Отдых', tone: 'rest' as const },
]

export function WorkoutsPage() {
  return (
    <AppFrame
      description="Недельный план и список упражнений на сегодня."
      eyebrow="Экран 02"
      sidebarProfile={sidebarProfile}
      title="Тренировки"
    >
      <section className="card placeholder-section">
        <div className="card-title">План на неделю</div>
        <div className="week">
          {weekPlan.map((item) => (
            <div className={item.tone ? `week-day ${item.tone}` : 'week-day'} key={item.day}>
              <div className="d">{item.day}</div>
              <div className="t">{item.title}</div>
            </div>
          ))}
        </div>
      </section>

      <section className="grid g-2">
        <div className="card frame">
          <div className="card-title">
            Сегодня · Кардио <span className="pill pending">не выполнена</span>
          </div>

          <div className="check-item">
            <span className="checkbox done" />
            <span className="check-label">Выпады</span>
            <span className="check-meta">3×15</span>
          </div>
          <div className="check-item">
            <span className="checkbox done" />
            <span className="check-label">Упражнения на пресс</span>
            <span className="check-meta">3×20</span>
          </div>
          <div className="check-item">
            <span className="checkbox" />
            <span className="check-label">Кардио</span>
            <span className="check-meta">20 мин</span>
          </div>

          <button className="btn btn-primary placeholder-button" type="button">
            Завершить тренировку
          </button>

          <div className="reward">
            <b>+80 XP</b>
            <span>Выносливость +2 · Дисциплина +1</span>
          </div>
        </div>

        <div className="card">
          <div className="card-title">Награды за тренировки</div>

          <div className="legend-list">
            <div className="legend-row">
              <span className="k">Силовая</span>
              <span className="v">+100 XP, Сила +2, Дисциплина +1</span>
            </div>
            <div className="legend-row">
              <span className="k">Кардио</span>
              <span className="v">+80 XP, Выносливость +2</span>
            </div>
            <div className="legend-row">
              <span className="k">Серия 7д</span>
              <span className="v">дополнительный бонус к опыту</span>
            </div>
          </div>

          <hr className="rule" />

          <div className="bar-row">
            <span className="bar-label">Выполнено за неделю</span>
            <span className="bar-value">4 из 7</span>
          </div>
          <NotchBar active={11} total={20} />
        </div>
      </section>

      <footer className="foot">ФОРМА — пример интерфейса для хакатона · экран workouts</footer>
    </AppFrame>
  )
}
