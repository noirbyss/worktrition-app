import { AppFrame } from '../components/app/AppFrame'
import { AppLink } from '../components/navigation/AppLink'

export function WorkoutsPage() {
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
      description="Временная защищенная страница раздела тренировок. Доступ сюда открывается только после сохранения анкеты."
      title="Тренировки"
    >
      <section className="panel-grid">
        <article className="panel">
          <h2 className="panel-title">Доступ работает</h2>
          <p className="panel-copy">
            Если пользователь бросил анкету и вернулся на сайт позже, guard все равно удержит его
            на onboarding-странице до успешной отправки профиля.
          </p>
        </article>

        <article className="panel">
          <h2 className="panel-title">Можно развивать дальше</h2>
          <p className="panel-copy">
            Здесь потом удобно разместить план тренировок, прогресс и карточки упражнений.
          </p>
        </article>
      </section>
    </AppFrame>
  )
}
