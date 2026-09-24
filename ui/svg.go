//go:build !wasm

package ui

import (
	"webtyp.com/svg"
	"webtyp.com/svg/sprite"
)

// Icons es un tipo que aloja IconID/IconSvg para su descubrimiento y uso en config.
type Icons struct{}

func (m *Icons) IconID() string {
	return ID
}

func (m *Icons) IconSvg() *sprite.Sprite {
	return sprite.NewSprite(
		sprite.Define(svg.Icon(ID), "0 0 24 24",
			sprite.Path("M12 12a4 4 0 1 0-4-4 4 4 0 0 0 4 4zm0 2c-3.3 0-8 1.7-8 5v1h16v-1c0-3.3-4.7-5-8-5z"),
		),
	)
}
