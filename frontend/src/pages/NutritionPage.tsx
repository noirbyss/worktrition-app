import { AppFrame } from '../components/app/AppFrame'
import { AppLink } from '../components/navigation/AppLink'

export function NutritionPage() {
  return (
    <AppFrame
      actions={
        <div className="button-row">
          <AppLink className="btn btn--ghost btn--small" href="/app">
            Аккаунт
          </AppLink>
          <AppLink className="btn btn--ghost btn--small" href="/profile">
            Профиль
          </AppLink>
        </div>
      }
      description="Временная защищенная страница раздела питания. Она открывается только после завершения анкеты."
      title="Питание"
    >
      <section className="panel-grid">
        <article className="panel">
          <h2 className="panel-title">Раздел подключен</h2>
          <p className="panel-copy">
            Маршрут уже защищен: пользователь без заполненной анкеты сюда не попадет даже по
            прямой ссылке.
          </p>
        </article>

        <article className="panel">
          <h2 className="panel-title">Следующий шаг</h2>
          <p className="panel-copy">
            Сюда позже можно добавить план питания, рекомендации и дневник приемов пищи.
          </p>
        </article>
      </section>
    </AppFrame>
  )
}
