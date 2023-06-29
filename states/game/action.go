package game

type ActorActions struct {
	Actor   Actor
	Actions []Action
}

type BulletActions struct {
	Bullet  *Bullet
	Actions []Action
}

type Action interface {
}

type ActionMove struct {
	X, Y float64
}

type ActionReflect struct {
	X, Y float64
}

type ActionDeflect struct {
	X, Y      float64
	Direction float64
}

type ActionShield struct {
}

type ActionSpawnBullets struct {
	Bullets []*Bullet
}

type ActionFindNearestActor struct {
	Actor Actor
}

type ActionSpawnParticle struct {
	Img   string
	X, Y  float64
	Angle float64
	Speed float64
	Life  int
}

type ActionSpawnEnemy struct {
	Name string
	ID   string
	X, Y float64
}
