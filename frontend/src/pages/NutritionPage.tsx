import { useState } from 'react'
import { AppFrame } from '../components/app/AppFrame'
import { NotchBar } from '../components/app/PlaceholderUi'

const sidebarProfile = {
  badge: '7',
  meta: '820 / 1000 XP',
  name: 'Артём Ковалёв',
}

const dayOrder = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс'] as const

type DayKey = (typeof dayOrder)[number]

type Macro = {
  prefix?: string
  suffix: string
  value: string
}

type Meal = {
  items: string
  macros: Macro[]
  name: string
  status: string
  tone: 'done' | 'pending'
}

type NutritionPlan = {
  calories: string
  carbs: string
  fats: string
  meals: Meal[]
  protein: string
  title: string
  totalStatus: string
  water: string
}

const nutritionDays: Record<DayKey, NutritionPlan> = {
  Пн: {
    calories: '1590 / 2100 ккал',
    carbs: '168 / 230 г',
    fats: '48 / 65 г',
    meals: [
      {
        items: 'Овсянка + банан',
        macros: [
          { suffix: 'ккал', value: '480' },
          { prefix: 'Б', suffix: 'г', value: '22' },
          { prefix: 'Ж', suffix: 'г', value: '14' },
          { prefix: 'У', suffix: 'г', value: '66' },
        ],
        name: 'Завтрак',
        status: 'выполнено',
        tone: 'done',
      },
      {
        items: 'Индейка + гречка + овощи',
        macros: [
          { suffix: 'ккал', value: '690' },
          { prefix: 'Б', suffix: 'г', value: '47' },
          { prefix: 'Ж', suffix: 'г', value: '18' },
          { prefix: 'У', suffix: 'г', value: '74' },
        ],
        name: 'Обед',
        status: 'выполнено',
        tone: 'done',
      },
      {
        items: 'Творог + ягоды',
        macros: [
          { suffix: 'ккал', value: '420' },
          { prefix: 'Б', suffix: 'г', value: '34' },
          { prefix: 'Ж', suffix: 'г', value: '16' },
          { prefix: 'У', suffix: 'г', value: '28' },
        ],
        name: 'Ужин',
        status: 'не отмечен',
        tone: 'pending',
      },
    ],
    protein: '103 / 140 г',
    title: 'Рацион на понедельник',
    totalStatus: '2 из 3 приёмов пищи',
    water: '1.2 / 2.5 л',
  },
  Вт: {
    calories: '1710 / 2100 ккал',
    carbs: '174 / 230 г',
    fats: '49 / 65 г',
    meals: [
      {
        items: 'Сырники + йогурт',
        macros: [
          { suffix: 'ккал', value: '510' },
          { prefix: 'Б', suffix: 'г', value: '24' },
          { prefix: 'Ж', suffix: 'г', value: '17' },
          { prefix: 'У', suffix: 'г', value: '63' },
        ],
        name: 'Завтрак',
        status: 'выполнено',
        tone: 'done',
      },
      {
        items: 'Говядина + картофель + салат',
        macros: [
          { suffix: 'ккал', value: '760' },
          { prefix: 'Б', suffix: 'г', value: '48' },
          { prefix: 'Ж', suffix: 'г', value: '22' },
          { prefix: 'У', suffix: 'г', value: '82' },
        ],
        name: 'Обед',
        status: 'выполнено',
        tone: 'done',
      },
      {
        items: 'Лосось + брокколи',
        macros: [
          { suffix: 'ккал', value: '440' },
          { prefix: 'Б', suffix: 'г', value: '36' },
          { prefix: 'Ж', suffix: 'г', value: '14' },
          { prefix: 'У', suffix: 'г', value: '29' },
        ],
        name: 'Ужин',
        status: 'не отмечен',
        tone: 'pending',
      },
    ],
    protein: '108 / 140 г',
    title: 'Рацион на вторник',
    totalStatus: '2 из 3 приёмов пищи',
    water: '1.5 / 2.5 л',
  },
  Ср: {
    calories: '1680 / 2100 ккал',
    carbs: '177 / 230 г',
    fats: '51 / 65 г',
    meals: [
      {
        items: 'Овсянка + яйца',
        macros: [
          { suffix: 'ккал', value: '500' },
          { prefix: 'Б', suffix: 'г', value: '25' },
          { prefix: 'Ж', suffix: 'г', value: '15' },
          { prefix: 'У', suffix: 'г', value: '65' },
        ],
        name: 'Завтрак',
        status: 'выполнено',
        tone: 'done',
      },
      {
        items: 'Курица + рис + овощи',
        macros: [
          { suffix: 'ккал', value: '700' },
          { prefix: 'Б', suffix: 'г', value: '50' },
          { prefix: 'Ж', suffix: 'г', value: '20' },
          { prefix: 'У', suffix: 'г', value: '80' },
        ],
        name: 'Обед',
        status: 'выполнено',
        tone: 'done',
      },
      {
        items: 'Рыба + овощи',
        macros: [
          { suffix: 'ккал', value: '480' },
          { prefix: 'Б', suffix: 'г', value: '38' },
          { prefix: 'Ж', suffix: 'г', value: '16' },
          { prefix: 'У', suffix: 'г', value: '32' },
        ],
        name: 'Ужин',
        status: 'не отмечен',
        tone: 'pending',
      },
    ],
    protein: '113 / 140 г',
    title: 'Рацион на среду',
    totalStatus: '2 из 3 приёмов пищи',
    water: '1.4 / 2.5 л',
  },
  Чт: {
    calories: '1620 / 2100 ккал',
    carbs: '170 / 230 г',
    fats: '47 / 65 г',
    meals: [
      {
        items: 'Тосты + омлет',
        macros: [
          { suffix: 'ккал', value: '470' },
          { prefix: 'Б', suffix: 'г', value: '23' },
          { prefix: 'Ж', suffix: 'г', value: '15' },
          { prefix: 'У', suffix: 'г', value: '58' },
        ],
        name: 'Завтрак',
        status: 'выполнено',
        tone: 'done',
      },
      {
        items: 'Тунец + паста + овощи',
        macros: [
          { suffix: 'ккал', value: '710' },
          { prefix: 'Б', suffix: 'г', value: '49' },
          { prefix: 'Ж', suffix: 'г', value: '18' },
          { prefix: 'У', suffix: 'г', value: '84' },
        ],
        name: 'Обед',
        status: 'выполнено',
        tone: 'done',
      },
      {
        items: 'Индейка + салат',
        macros: [
          { suffix: 'ккал', value: '440' },
          { prefix: 'Б', suffix: 'г', value: '35' },
          { prefix: 'Ж', suffix: 'г', value: '14' },
          { prefix: 'У', suffix: 'г', value: '28' },
        ],
        name: 'Ужин',
        status: 'не отмечен',
        tone: 'pending',
      },
    ],
    protein: '107 / 140 г',
    title: 'Рацион на четверг',
    totalStatus: '2 из 3 приёмов пищи',
    water: '1.3 / 2.5 л',
  },
  Пт: {
    calories: '1765 / 2100 ккал',
    carbs: '182 / 230 г',
    fats: '54 / 65 г',
    meals: [
      {
        items: 'Гранола + творог',
        macros: [
          { suffix: 'ккал', value: '520' },
          { prefix: 'Б', suffix: 'г', value: '26' },
          { prefix: 'Ж', suffix: 'г', value: '16' },
          { prefix: 'У', suffix: 'г', value: '67' },
        ],
        name: 'Завтрак',
        status: 'выполнено',
        tone: 'done',
      },
      {
        items: 'Курица терияки + рис',
        macros: [
          { suffix: 'ккал', value: '730' },
          { prefix: 'Б', suffix: 'г', value: '51' },
          { prefix: 'Ж', suffix: 'г', value: '21' },
          { prefix: 'У', suffix: 'г', value: '85' },
        ],
        name: 'Обед',
        status: 'выполнено',
        tone: 'done',
      },
      {
        items: 'Стейк + овощи',
        macros: [
          { suffix: 'ккал', value: '515' },
          { prefix: 'Б', suffix: 'г', value: '40' },
          { prefix: 'Ж', suffix: 'г', value: '17' },
          { prefix: 'У', suffix: 'г', value: '30' },
        ],
        name: 'Ужин',
        status: 'не отмечен',
        tone: 'pending',
      },
    ],
    protein: '117 / 140 г',
    title: 'Рацион на пятницу',
    totalStatus: '2 из 3 приёмов пищи',
    water: '1.6 / 2.5 л',
  },
  Сб: {
    calories: '1490 / 2100 ккал',
    carbs: '153 / 230 г',
    fats: '44 / 65 г',
    meals: [
      {
        items: 'Блины + творог',
        macros: [
          { suffix: 'ккал', value: '490' },
          { prefix: 'Б', suffix: 'г', value: '22' },
          { prefix: 'Ж', suffix: 'г', value: '15' },
          { prefix: 'У', suffix: 'г', value: '62' },
        ],
        name: 'Завтрак',
        status: 'выполнено',
        tone: 'done',
      },
      {
        items: 'Паста + индейка',
        macros: [
          { suffix: 'ккал', value: '640' },
          { prefix: 'Б', suffix: 'г', value: '43' },
          { prefix: 'Ж', suffix: 'г', value: '16' },
          { prefix: 'У', suffix: 'г', value: '74' },
        ],
        name: 'Обед',
        status: 'не отмечен',
        tone: 'pending',
      },
      {
        items: 'Салат с тунцом',
        macros: [
          { suffix: 'ккал', value: '360' },
          { prefix: 'Б', suffix: 'г', value: '31' },
          { prefix: 'Ж', suffix: 'г', value: '13' },
          { prefix: 'У', suffix: 'г', value: '17' },
        ],
        name: 'Ужин',
        status: 'не отмечен',
        tone: 'pending',
      },
    ],
    protein: '96 / 140 г',
    title: 'Рацион на субботу',
    totalStatus: '1 из 3 приёмов пищи',
    water: '1.1 / 2.5 л',
  },
  Вс: {
    calories: '1380 / 2100 ккал',
    carbs: '149 / 230 г',
    fats: '41 / 65 г',
    meals: [
      {
        items: 'Йогурт + мюсли',
        macros: [
          { suffix: 'ккал', value: '430' },
          { prefix: 'Б', suffix: 'г', value: '19' },
          { prefix: 'Ж', suffix: 'г', value: '12' },
          { prefix: 'У', suffix: 'г', value: '54' },
        ],
        name: 'Завтрак',
        status: 'выполнено',
        tone: 'done',
      },
      {
        items: 'Рыба + киноа',
        macros: [
          { suffix: 'ккал', value: '610' },
          { prefix: 'Б', suffix: 'г', value: '41' },
          { prefix: 'Ж', suffix: 'г', value: '16' },
          { prefix: 'У', suffix: 'г', value: '62' },
        ],
        name: 'Обед',
        status: 'не отмечен',
        tone: 'pending',
      },
      {
        items: 'Овощной суп + курица',
        macros: [
          { suffix: 'ккал', value: '340' },
          { prefix: 'Б', suffix: 'г', value: '28' },
          { prefix: 'Ж', suffix: 'г', value: '13' },
          { prefix: 'У', suffix: 'г', value: '19' },
        ],
        name: 'Ужин',
        status: 'не отмечен',
        tone: 'pending',
      },
    ],
    protein: '88 / 140 г',
    title: 'Рацион на воскресенье',
    totalStatus: '1 из 3 приёмов пищи',
    water: '0.9 / 2.5 л',
  },
}

export function NutritionPage() {
  const [activeDay, setActiveDay] = useState<DayKey>('Ср')
  const activePlan = nutritionDays[activeDay]

  return (
    <AppFrame
      description="Рацион по дням недели, БЖУ и трекер воды."
      eyebrow="Экран 03"
      sidebarProfile={sidebarProfile}
      title="Питание"
    >
      <div className="day-tabs">
        {dayOrder.map((day) => (
          <button
            className={day === activeDay ? 'day-tab active' : 'day-tab'}
            key={day}
            onClick={() => setActiveDay(day)}
            type="button"
          >
            {day}
          </button>
        ))}
      </div>

      <section className="grid g-2">
        <div className="card">
          <div className="card-title">
            {activePlan.title} <span className="pill done">{activePlan.totalStatus}</span>
          </div>

          {activePlan.meals.map((meal) => (
            <div className="meal" key={meal.name}>
              <div className="meal-top">
                <span className="meal-name">{meal.name}</span>
                <span className={`pill ${meal.tone}`}>{meal.status}</span>
              </div>
              <p className="meal-items">{meal.items}</p>
              <div className="macro-row">
                {meal.macros.map((macro) => (
                  <span className="macro" key={`${meal.name}-${macro.prefix ?? 'ккал'}-${macro.value}`}>
                    {macro.prefix ? `${macro.prefix} ` : null}
                    <b>{macro.value}</b> {macro.suffix}
                  </span>
                ))}
              </div>
            </div>
          ))}
        </div>

        <div className="card">
          <div className="card-title">Норма на день</div>

          <MacroBlock active={16} label="Калории" value={activePlan.calories} />

          <hr className="rule" />

          <MacroBlock active={16} label="Белки" value={activePlan.protein} variant="balance" />
          <MacroBlock active={16} label="Жиры" marginTop value={activePlan.fats} variant="balance" />
          <MacroBlock active={15} label="Углеводы" marginTop value={activePlan.carbs} variant="balance" />

          <hr className="rule" />

          <MacroBlock active={11} label="Вода" value={activePlan.water} variant="water" />

          <div className="reward">
            <b>+20 XP</b>
            <span>за дневную норму воды, когда счётчик дойдет до конца</span>
          </div>
        </div>
      </section>

      <footer className="foot">ФОРМА — пример интерфейса для хакатона · экран nutrition</footer>
    </AppFrame>
  )
}

function MacroBlock({
  active,
  label,
  marginTop = false,
  value,
  variant,
}: {
  active: number
  label: string
  marginTop?: boolean
  value: string
  variant?: 'balance' | 'water'
}) {
  return (
    <div className={marginTop ? 'metric-block metric-block--spaced' : 'metric-block'}>
      <div className="bar-row">
        <span className="bar-label">{label}</span>
        <span className="bar-value">{value}</span>
      </div>
      <NotchBar active={active} total={20} variant={variant} />
    </div>
  )
}
