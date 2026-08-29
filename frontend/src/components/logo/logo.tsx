import heroImg from '../../assets/hero.png'

function Logo({ compact = false }: { compact?: boolean }) {
  return (
    <div className={compact ? 'brand-lockup brand-lockup--compact' : 'brand-lockup'}>
      <div className="brand-mark">
        <img alt="" aria-hidden="true" src={heroImg} />
      </div>
      <div className="brand-copy">
        <span className="brand-name">WORKTRITION</span>
        <span className="brand-tag">eat clean, train hard, recover smart</span>
      </div>
    </div>
  )
}

export default Logo
