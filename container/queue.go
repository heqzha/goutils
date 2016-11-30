package container

type Queue []interface{}

func (q *Queue) Clear() {
	*q = []interface{}{}
}

func (q *Queue) Push(n interface{}) {
	*q = append(*q, n)
}

func (q *Queue) Pop() interface{} {
	if len(*q) > 0 {
		n := (*q)[0]
		*q = (*q)[1:]
		return n
	}
	return nil
}

func (q *Queue) Len() int {
	return len(*q)
}
