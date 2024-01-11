package menu

import (
	"fmt"
	"image/color"
	"net"
	"time"

	"github.com/ketMix/retromancer/states"
	"github.com/ketMix/retromancer/states/game"

	rnet "github.com/ketMix/retromancer/net"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/ketMix/retromancer/resources"
)

type Lobby struct {
	clickSound      *resources.Sound
	items           []resources.MenuItem
	multiplayerItem *resources.ButtonItem
	joinItem        *resources.ButtonItem
	cancelItem      *resources.ButtonItem
	hostItem        *resources.ButtonItem
	backItem        *resources.TextItem
	lobbyItem       *resources.InputItem
	playerEntries   []*PlayerEntry
	overlay         game.Overlay
	shouldStart     bool
	difficulty      states.Difficulty
	net             rnet.ServerClient
}

func init() {
	rnet.RegisterMessage(HatMessage{})
	rnet.RegisterMessage(StartMessage{})
}

func (s *Lobby) Init(ctx states.Context) error {
	s.net.Init()

	s.overlay.Init(ctx)
	//
	s.clickSound = ctx.R.GetAs("sounds", "click", (*resources.Sound)(nil)).(*resources.Sound)

	s.playerEntries = append(s.playerEntries, &PlayerEntry{
		player: game.NewLocalPlayer(),
	})

	for _, e := range s.playerEntries {
		e.Init(s, ctx)
	}

	s.multiplayerItem = &resources.ButtonItem{
		Text: ctx.L.Get("Multiplayer"),
		X:    500,
		Y:    20,
		Callback: func() bool {
			s.clickSound.Play(1.0)
			for _, item := range s.items {
				if item == s.multiplayerItem {
					s.multiplayerItem.SetHidden(true)
					s.lobbyItem.SetHidden(false)
					s.joinItem.SetHidden(false)
					s.hostItem.SetHidden(false)
					s.playerEntries = append(s.playerEntries, &PlayerEntry{})
					s.playerEntries[len(s.playerEntries)-1].Init(s, ctx)
					break
				}
			}
			return false
		},
	}

	s.lobbyItem = &resources.InputItem{
		X:           350,
		Y:           20,
		Width:       150,
		Placeholder: ctx.L.Get("Address"),
		Callback: func() bool {
			return false
		},
	}
	s.lobbyItem.SetHidden(true)

	s.joinItem = &resources.ButtonItem{
		Text: ctx.L.Get("Join"),
		X:    450 + 50,
		Y:    20,
		Callback: func() bool {
			s.clickSound.Play(1.0)
			// TODO: Create network server:
			//   - Check if lobby is an address or ip
			//      - If so, begin directly hosting and wait for a client to connect.
			//			- If not, connect to magnet's matchmaker with the lobby as the advertisement and begin waiting for a client to connect.
			s.JoinHost(s.lobbyItem.Text)
			return true
		},
	}
	s.joinItem.SetHidden(true)

	s.hostItem = &resources.ButtonItem{
		Text: ctx.L.Get("Host"),
		X:    450,
		Y:    20,
		Callback: func() bool {
			s.clickSound.Play(1.0)
			// TODO: Create network client:
			//   - Check if lobby is an address or ip
			//      - If so, directly connect to it
			//			- If not, connect to magnet's matchmaker and use the lobby as the target name. Wait for response, and if an ip:port, directly connect to it using the same socket.
			s.StartHost(s.lobbyItem.Text)
			return true
		},
	}
	s.hostItem.SetHidden(true)

	s.cancelItem = &resources.ButtonItem{
		Text: ctx.L.Get("Cancel"),
		X:    450,
		Y:    20,
		Callback: func() bool {
			s.clickSound.Play(1.0)
			s.net.Close()
			return true
		},
	}
	s.cancelItem.SetHidden(true)

	s.backItem = &resources.TextItem{
		Text: ctx.L.Get("Back"),
		X:    30,
		Y:    335,
		Callback: func() bool {
			s.clickSound.Play(1.0)
			ctx.StateMachine.PopState(nil)
			return false
		},
	}
	s.items = append(s.items, s.backItem, s.multiplayerItem, s.lobbyItem, s.joinItem, s.hostItem, s.cancelItem)

	return nil
}

func (s *Lobby) Finalize(ctx states.Context) error {
	return nil
}

func (s *Lobby) Enter(ctx states.Context, v interface{}) error {
	return nil
}

func (s *Lobby) Update(ctx states.Context) error {
	// Handle net stuff.
	select {
	case ev := <-s.net.EventChan:
		switch e := ev.(type) {
		case rnet.EventHosting:
			fmt.Println("now hosting....")
		case rnet.EventJoining:
			fmt.Println("now joining....")
		case rnet.EventJoined:
			fmt.Println("now joined....")
		case rnet.EventConnect:
			s.AddNetPlayer(ctx, e.Peer)
		case rnet.EventDisconnect:
			s.RemoveNetPlayer()
		case rnet.EventClosed:
			s.CancelNetworking()
		case rnet.EventMessage:
			switch msg := e.Message.(type) {
			case HatMessage:
				if e := s.GetNetPlayer(e.Peer); e != nil {
					e.hatIndex = msg.Hat
					e.SyncHat(ctx)
				}
			case StartMessage:
				if !s.net.Hosting {
					s.shouldStart = true
				}
			}
		}
	default:
	}

	s.overlay.Update(ctx)

	s.lobbyItem.Update()

	// Check for controller button hit to activate player 2.
	for i, gamepadID := range resources.GetFunctionalGamepads() {
		m := resources.GetBestGamemap(gamepadID)
		if resources.GetButton(m, gamepadID, resources.ButtonStart) {
			if len(s.playerEntries) == 1 {
				s.playerEntries = append(s.playerEntries, &PlayerEntry{})
				s.playerEntries[len(s.playerEntries)-1].Init(s, ctx)
			}
			pl := game.NewLocalPlayer()
			s.playerEntries[1].player = pl
			s.playerEntries[1].controllerIndex = i
			s.playerEntries[1].useController = true
			s.playerEntries[1].SyncController(ctx)
			pl.GamepadID = gamepadID
			pl.GamepadMap = m
			// TODO: Stop network stuff and hide host/join.
			s.hostItem.SetHidden(true)
			s.joinItem.SetHidden(true)
			s.lobbyItem.SetHidden(true)
		}
	}

	x := -(len(s.playerEntries) - 1) * 150 / 2
	for i, e := range s.playerEntries {
		e.Update(ctx, float64(x+i*150))
	}

	x, y := ebiten.CursorPosition()
	for _, m := range s.items {
		m.CheckState(float64(x), float64(y))
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
		for _, m := range s.items {
			if m.Hovered() {
				if m.Activate() {
					return nil
				}
			}
		}
	}

	if s.shouldStart {
		if s.net.Hosting {
			for _, p := range s.net.Peers() {
				p.Send(StartMessage{})
			}
		}

		var seed int64
		// Seed rand with the host's id.
		if s.net.Running {
			if s.net.Hosting {
				seed = int64(s.net.ID())
			} else {
				seed = int64(s.net.Peers()[0].ID())
			}
		} else {
			seed = time.Now().UnixNano()
		}

		players := make([]game.Player, len(s.playerEntries))
		for i, e := range s.playerEntries {
			players[i] = e.player
		}
		// FIXME: Need to agree w/ players to start (or assume host has full control).
		ctx.StateMachine.PopState(nil)
		ctx.StateMachine.PushState(&game.World{
			StartingMap: "start",
			ShowHints:   true,
			Players:     players,
			Net:         s.net,
			Seed:        seed,
			Difficulty:  &s.difficulty,
		})
	}

	return nil
}

func (s *Lobby) Draw(ctx states.DrawContext) {
	ctx.Text.SetColor(color.White)
	for _, e := range s.playerEntries {
		e.Draw(ctx)
	}
	for _, m := range s.items {
		m.Draw(ctx)
	}
	s.overlay.Draw(ctx)
}

// Networking stuff

func (s *Lobby) CancelNetworking() {
	s.hostItem.SetHidden(false)
	s.joinItem.SetHidden(false)
	s.cancelItem.SetHidden(true)
}

func (s *Lobby) StartHost(address string) error {
	// Use the matchmaker if the address is not an ip:port.
	if _, err := net.ResolveUDPAddr("udp", address); err != nil {
		s.net.UseMatchmaker = true
	}
	s.net.Hosting = true

	if s.net.UseMatchmaker {
		if err := s.net.Open(""); err != nil {
			fmt.Println(err)
			return err
		}
		// TODO: Connect to matchmaker and begin advertising.
		fmt.Println("connecting to matchmaker...")
	} else {
		if err := s.net.Open(address); err != nil {
			fmt.Println(err)
			return err
		}
	}
	fmt.Println("opened....")

	s.hostItem.SetHidden(true)
	s.joinItem.SetHidden(true)
	s.cancelItem.SetHidden(false)

	return nil
}

func (s *Lobby) JoinHost(address string) error {
	if _, err := net.ResolveUDPAddr("udp", address); err != nil {
		s.net.UseMatchmaker = true
	}
	s.net.Hosting = false

	if err := s.net.Open(""); err != nil {
		fmt.Println(err)
		return err
	}

	if s.net.UseMatchmaker {
		// TODO: Connect to matchmaker and wait for response.
		fmt.Println("connecting to matchmaker...")
	} else {
		fmt.Println("connecting directly...")
		if err := s.net.ConnectTo(address); err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("connected...")
	}
	fmt.Println("opened...")

	s.hostItem.SetHidden(true)
	s.joinItem.SetHidden(true)
	s.cancelItem.SetHidden(false)

	return nil
}

func (s *Lobby) AddNetPlayer(ctx states.Context, peer *rnet.Peer) {
	if len(s.playerEntries) == 1 {
		return
	}
	pl := game.NewRemotePlayer(peer)
	s.playerEntries[1].player = pl
	s.playerEntries[1].hatLeft.SetHidden(true)
	s.playerEntries[1].hatRight.SetHidden(true)
	s.playerEntries[1].controllerLeft.SetHidden(true)
	s.playerEntries[1].controllerRight.SetHidden(true)
	s.playerEntries[1].controllerItem.Sprite = resources.NewSprite(ctx.R.Get("images", "network").(*ebiten.Image))
	s.playerEntries[1].controllerItem.Sprite.Centered = true
	s.playerEntries[1].waitingText.SetHidden(true)
	s.playerEntries[1].startText.SetHidden(true)

	// Send our hat to them.
	peer.Send(HatMessage{
		Hat: s.playerEntries[0].hatIndex,
	})
}

func (s *Lobby) RemoveNetPlayer() {
	if len(s.playerEntries) == 1 {
		return
	}

	s.playerEntries[1].player = nil
	s.playerEntries[1].hatLeft.SetHidden(false)
	s.playerEntries[1].hatRight.SetHidden(false)
	s.playerEntries[1].controllerLeft.SetHidden(false)
	s.playerEntries[1].controllerRight.SetHidden(false)
	s.playerEntries[1].waitingText.SetHidden(false)
	s.playerEntries[1].startText.SetHidden(false)

	/*for i, e := range s.playerEntries {
		if pl, ok := e.player.(*game.RemotePlayer); ok {
			if pl.Peer() == peer {
				s.playerEntries = append(s.playerEntries[:i], s.playerEntries[i+1:]...)
				return
			}
		}
	}*/
}

func (s *Lobby) GetNetPlayer(peer *rnet.Peer) *PlayerEntry {
	for _, e := range s.playerEntries {
		if pl, ok := e.player.(*game.RemotePlayer); ok {
			if pl.Peer() == peer {
				return e
			}
		}
	}
	return nil
}

func (s *Lobby) SendToNetPlayers(msg rnet.Message) {
	for _, e := range s.playerEntries {
		if pl, ok := e.player.(*game.RemotePlayer); ok {
			pl.Peer().Send(msg)
		}
	}
}
