package game

import (
	"ebijam23/resources"
	"ebijam23/states"
	"image/color"
	"time"
)

type HintGroup struct {
	Items   []string
	OffsetY float64
	Prefix  string
}

type Hints struct {
	hintGroup          map[string]HintGroup
	activeGroups       []string
	activeGroupIndices []int
	active             bool
	ticker             int
	vfxs               resources.VFXList
}

func (h *Hints) Update(ctx states.Context) error {
	if !h.active {
		return nil
	}
	if len(h.activeGroups) == 0 {
		return nil
	}
	h.ticker++
	if h.ticker > 0 {
		h.ticker = -230

		for i := 0; i < len(h.activeGroups); {
			h.activeGroupIndices[i]++
			if h.activeGroupIndices[i] >= len(h.hintGroup[h.activeGroups[i]].Items) {
				h.activeGroupIndices = append(h.activeGroupIndices[:i], h.activeGroupIndices[i+1:]...)
				h.activeGroups = append(h.activeGroups[:i], h.activeGroups[i+1:]...)
				continue
			} else {
				i++
			}
		}
		for i, g := range h.activeGroups {
			hint := h.hintGroup[g].Items[h.activeGroupIndices[i]]
			h.vfxs.Add(&resources.Text{
				X:            320,
				Y:            325 + h.hintGroup[g].OffsetY,
				Scale:        1.0,
				InDuration:   1 * time.Second,
				HoldDuration: 5 * time.Second,
				OutDuration:  1 * time.Second,
				Text:         h.hintGroup[g].Prefix + ctx.L(hint),
				Color:        color.NRGBA{0xff, 0xff, 0x44, 0xff},
				OutlineColor: color.NRGBA{0x00, 0x00, 0x00, 0xff},
				Outline:      true,
			})
		}
	}
	return nil
}

func (h *Hints) AddHintGroup(group string, hintGroup HintGroup) {
	h.hintGroup[group] = hintGroup
}

func (h *Hints) ActivateGroup(group string) {
	h.activeGroups = append(h.activeGroups, group)
	h.activeGroupIndices = append(h.activeGroupIndices, -1)
}

func (h *Hints) DeactivateGroups() {
	h.activeGroups = []string{}
	h.activeGroupIndices = []int{}
}

func (h *Hints) Draw(ctx states.DrawContext) {
	if !h.active {
		return
	}
	h.vfxs.Process(ctx, nil)
}
