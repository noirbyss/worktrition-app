import { AppFrame } from '../components/app/AppFrame'
import { NotchBar } from '../components/app/PlaceholderUi'

const sidebarProfile = {
  badge: '7',
  meta: '820 / 1000 XP',
  name: 'Артём Ковалёв',
}

export function StatsPage() {
  return (
    <AppFrame
      description="Реальный прогресс: вес, тренировки, питание."
      eyebrow="Экран 04"
      sidebarProfile={sidebarProfile}
      title="Статистика"
    >
      <section className="grid g-2 placeholder-grid">
        <div className="card">
          <div className="card-title">Вес — динамика</div>
          <svg className="weight-chart" height="160" viewBox="0 0 460 160" width="100%">
            <defs>
              <linearGradient id="weightChartGradient" x1="0" x2="1" y1="0" y2="0">
                <stop offset="0" stopColor="#FFD23F" />
                <stop offset="1" stopColor="#FF7A1A" />
              </linearGradient>
            </defs>
            <polyline
              fill="none"
              points="10,30 110,55 210,78 310,100 410,118"
              stroke="url(#weightChartGradient)"
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="3"
            />
            <g fill="#FF7A1A">
              <circle cx="10" cy="30" r="4" />
              <circle cx="110" cy="55" r="4" />
              <circle cx="210" cy="78" r="4" />
              <circle cx="310" cy="100" r="4" />
              <circle cx="410" cy="118" r="4" />
            </g>
            <g fill="#9195A8" fontFamily="IBM Plex Mono, monospace" fontSize="11">
              <text x="0" y="18">
                90
              </text>
              <text x="100" y="43">
                88
              </text>
              <text x="200" y="66">
                86
              </text>
              <text x="300" y="88">
                85
              </text>
              <text x="386" y="106">
                84 кг
              </text>
            </g>
            <line stroke="rgba(255,255,255,.1)" strokeWidth="1" x1="0" x2="460" y1="145" y2="145" />
            <g fill="#565A70" fontFamily="IBM Plex Mono, monospace" fontSize="10">
              <text x="4" y="158">
                нед. 1
              </text>
              <text x="192" y="158">
                нед. 3
              </text>
              <text x="392" y="158">
                нед. 5
              </text>
            </g>
          </svg>
        </div>

        <div className="card">
          <div className="card-title">Физические показатели</div>
          <SimpleStatRow label="Текущий вес" value="88 кг" />
          <SimpleStatRow label="Начальный вес" value="90 кг" />
          <SimpleStatRow label="Целевой вес" value="80 кг" />
          <SimpleStatRow label="Осталось" value="8 кг" />
          <SimpleStatRow label="ИМТ" value="27.2" />
        </div>
      </section>

      <section className="grid g-2">
        <div className="card">
          <div className="card-title">Тренировки</div>
          <ProgressStatRow active={15} label="Выполнено" value="14" />
          <ProgressStatRow active={15} label="Процент" value="82%" />
          <ProgressStatRow active={3} label="Пропущено" value="3" />
          <ProgressStatRow active={9} label="Серия дней" value="5" />
          <ProgressStatRow active={11} label="Общее время" value="6ч 40м" />
        </div>

        <div className="card">
          <div className="card-title">Питание</div>
          <ProgressStatRow active={16} label="Соблюдение КБЖУ" value="88%" variant="balance" />
          <ProgressStatRow active={16} label="Выполнение плана" value="91%" variant="balance" />
          <ProgressStatRow active={14} label="Норма воды" value="76%" variant="water" />
        </div>
      </section>

      <footer className="foot">ФОРМА — пример интерфейса для хакатона · экран stats</footer>
    </AppFrame>
  )
}

function ProgressStatRow({
  active,
  label,
  value,
  variant,
}: {
  active: number
  label: string
  value: string
  variant?: 'balance' | 'water'
}) {
  return (
    <div className="stat-row">
      <span className="stat-label">{label}</span>
      <NotchBar active={active} total={18} variant={variant} />
      <span className="stat-value">{value}</span>
    </div>
  )
}

function SimpleStatRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="stat-row stat-row--simple">
      <span className="stat-label">{label}</span>
      <span />
      <span className="stat-value">{value}</span>
    </div>
  )
}
