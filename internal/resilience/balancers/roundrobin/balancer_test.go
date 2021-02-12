package roundrobin

import "testing"

func TestBalancer(t *testing.T) {
	t.Parallel()

	t.Run("balance", func(t *testing.T) {
		g := group{size: 4}
		b := New(g)
		for i := 0; i < 10; i++ {
			idx, err := b.Index()
			if err != nil {
				t.Errorf("expected err to be nil %v", err)
			}

			if idx != uint64(i%4) {
				t.Errorf("expected idx %v to be %v", idx, (i % 5))
			}
		}
	})
}

type group struct {
	size int
}

func (g group) Size() int {
	return g.size
}
