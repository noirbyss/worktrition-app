type NotchBarVariant = 'default' | 'strength' | 'endurance' | 'discipline' | 'balance' | 'water'

export function NotchBar({
  active,
  large = false,
  total = 20,
  variant = 'default',
}: {
  active: number
  large?: boolean
  total?: number
  variant?: NotchBarVariant
}) {
  const ticks = Array.from({ length: total }, (_, index) => index < active)
  const className = ['notch-bar', variant !== 'default' ? variant : null, large ? 'lg' : null]
    .filter(Boolean)
    .join(' ')

  return (
    <div className={className}>
      {ticks.map((isActive, index) => (
        <span className={isActive ? 'tick on' : 'tick'} key={index} />
      ))}
    </div>
  )
}
