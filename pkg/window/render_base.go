package window

type RenderBase struct {
	z  int
	id int
}

func NewRenderBase(z int) RenderBase {
	return RenderBase{
		z:  z,
		id: GetNextId(),
	}
}

func (r *RenderBase) Z() int {
	return r.z
}

func (r *RenderBase) Id() int {
	return r.id
}
