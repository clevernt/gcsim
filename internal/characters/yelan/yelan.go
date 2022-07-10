package yelan

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

const (
	breakthroughStatus = "yelan_breakthrough"
	c6Status           = "yelan_c6"
	burstStatus        = "yelanburst"
)

func init() {
	core.RegisterCharFunc(keys.Yelan, NewChar)
}

type char struct {
	*tmpl.Character
	c2icd        int
	burstDiceICD int
	burstTickSrc int
	c6count      int
	c4count      int //keep track of number of enemies tagged
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Hydro
	c.EnergyMax = 70
	c.Weapon.Class = weapon.WeaponClassBow
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5
	c.CharZone = character.ZoneLiyue

	c.c2icd = 0
	c.c6count = 0

	barb, ok := p.Params["barb"]
	if ok && barb > 0 {
		c.SetTag(breakthroughStatus, 1)
	}

	if c.Base.Cons >= 1 {
		c.SetNumCharges(action.ActionSkill, 2)
	}

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a1()
	c.burstStateHook()
	return nil
}