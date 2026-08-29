import worktritionLogo from '../../assets/worktrition-logo.png'

function Logo({ compact = false }: { compact?: boolean }) {
  return (
    <img
      alt="Worktrition"
      className={compact ? 'brand-image brand-image--compact' : 'brand-image'}
      src={worktritionLogo}
    />
  )
}

export default Logo
