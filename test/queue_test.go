package test

import (
	"testing"

	"github.com/heqzha/goutils/container"
)

func TestMtxGroupQueue(t *testing.T) {
	q := container.MtxGroupQueue{}
	q.Init()
	for _, grp := range []string{"g1", "g2", "g3"} {
		for i := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0} {
			q.Push(grp, i)
		}

		for q.Len(grp) != 0 {
			t.Log(q.Pop(grp))
		}
	}

	for _, grp := range []string{"g1", "g2", "g3"} {
		for i := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0} {
			q.Push(grp, i)
		}
		t.Log("group", grp)
		t.Log("len", q.Len(grp))
		q.Clear(grp)
		t.Log("len, after clear", q.Len(grp))
	}

	for _, grp := range []string{"g1", "g2", "g3"} {
		for i := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0} {
			q.Push(grp, i)
		}
	}
	t.Log("groups", q.Groups())
	t.Log("number of groups", q.GroupsLen())
	q.ClearAll()
	t.Log("number of groups, after clear", q.GroupsLen())

}
