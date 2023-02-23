package raiden

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int

const chargeHitmark = 22

func init() {
	// charge -> x
	chargeFrames = frames.InitAbilSlice(37) //n1, skill, burst all at 37
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionSwap] = 36
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	if c.StatusIsActive(BurstKey) {
		return c.swordCharge(p)
	}

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Charge Attack",
		AttackTag:          attacks.AttackTagExtra,
		ICDTag:             combat.ICDTagExtraAttack,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeSlash,
		Element:            attributes.Physical,
		Durability:         25,
		HitlagHaltFrames:   0.02 * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
		Mult:               charge[c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: -0.1}, 2.8, 4.8),
		chargeHitmark,
		chargeHitmark,
	)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmark,
		State:           action.ChargeAttackState,
	}
}

var swordCAFrames []int
var swordCAHitmarks = []int{24, 31}

func init() {
	// charge (burst) -> x
	swordCAFrames = frames.InitAbilSlice(56)
	swordCAFrames[action.ActionDash] = swordCAHitmarks[len(swordCAHitmarks)-1]
	swordCAFrames[action.ActionJump] = swordCAHitmarks[len(swordCAHitmarks)-1]
}

func (c *char) swordCharge(p map[string]int) action.ActionInfo {
	for i, mult := range chargeSword {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               "Musou Isshin (Charge Attack)",
			AttackTag:          attacks.AttackTagElementalBurst,
			ICDTag:             combat.ICDTagNormalAttack,
			ICDGroup:           combat.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeSlash,
			Element:            attributes.Electro,
			Durability:         25,
			Mult:               mult[c.TalentLvlBurst()] + resolveBonus[c.TalentLvlBurst()]*c.stacksConsumed,
			HitlagHaltFrames:   0.02 * 60, //all raiden normals have 0.02s hitlag
			HitlagFactor:       0.01,
			CanBeDefenseHalted: true,
		}
		// Sword hits are dynamic - group snapshots with damage proc
		if i == 0 { // Only the last hit has hitlag
			ai.HitlagHaltFrames = 0
			ai.CanBeDefenseHalted = false
		}
		if c.Base.Cons >= 2 {
			ai.IgnoreDefPercent = .6
		}
		c.QueueCharTask(func() {
			c.Core.QueueAttack(
				ai,
				combat.NewBoxHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: -0.1}, 7.5, 8),
				0,
				0,
				c.burstRestorefunc,
				c.c6,
			)
		}, swordCAHitmarks[i])
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(swordCAFrames),
		AnimationLength: swordCAFrames[action.InvalidAction],
		CanQueueAfter:   swordCAHitmarks[len(swordCAHitmarks)-1],
		State:           action.ChargeAttackState,
	}
}
