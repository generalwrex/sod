package core

// Rerpresents an absorption effect, e.g. Power Word: Shield.
type Shield struct {
	Spell *Spell

	// Embed Aura so we can use Ishieldctive/Refresh/etc directly.
	*Aura
}

func (shield *Shield) Apply(sim *Simulation, shieldAmount float64) {
	shield.Aura.Deactivate(sim)
	shield.Aura.Priority = shieldAmount
	shield.Aura.Activate(sim)

	caster := shield.Spell.Unit
	target := shield.Aura.Unit

	threat := 0.0 // TODO
	shield.Spell.SpellMetrics[target.UnitIndex].TotalThreat += threat
	shield.Spell.SpellMetrics[target.UnitIndex].TotalShielding += shieldAmount

	if sim.Log != nil {
		caster.Log(sim, "%s %s hit for %0.3f shielding. (Threat: %0.3f)", target.LogLabel(), shield.Spell.ActionID, shieldAmount, threat)
	}
}

func NewShield(config Shield) *Shield {
	shield := &Shield{}
	*shield = config

	oldOnGain := shield.Aura.OnGain
	oldOnExpire := shield.Aura.OnExpire
	shield.Aura.OnGain = func(aura *Aura, sim *Simulation) {
		if oldOnGain != nil {
			oldOnGain(aura, sim)
		}
	}
	shield.Aura.OnExpire = func(aura *Aura, sim *Simulation) {
		if oldOnExpire != nil {
			oldOnExpire(aura, sim)
		}
	}

	return shield
}
