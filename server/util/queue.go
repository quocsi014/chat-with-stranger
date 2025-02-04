package util

type Queue struct{
	items []interface{}
}

func NewQueue() *Queue{
	return &Queue{
		items: make([]interface{}, 0),
	}
}

func (q *Queue)EnQueue(item interface{}){
	q.items = append(q.items, item)
}

func (q *Queue)DeQueue() interface{}{
	if len(q.items) == 0{
		return nil
	}
	item := q.items[0]
	q.items = q.items[1:]

	return item
}


func (q *Queue)Size() int{
	return len(q.items)
}