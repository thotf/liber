package parser

type listObject struct {
	list	array
	cmap    map[Value]Value
	islist  bool
}

func Newlist () listObject {
	return listObject{
		list:   NewAarray(),
		cmap:   make(map[Value]Value),
		islist: true,
	}
}

func (l *listObject) add (index Value,value Value) {
	if index.T == VT_STRING && l.islist {   //当list为数组且index不为int时讲数据全部复制到map
		l.listDataToCmap()
		l.islist = false
	}

	//将当前list当作数组使用
	if l.islist {
		l.list.append(value)
		return
	}

	l.cmap[index] = value
}

func (l *listObject) listDataToCmap () {
	for index,v := range l.list.buf {
		a := new(Value)
		a.T = VT_INT
		a.V = index
		l.cmap[*a] = v
	}
	l.list = array{}
}

func (l *listObject) Set (index Value,value Value) {
	if l.islist && index.T == VT_INT {
		l.list.set(index.V.(int),value)
		return
	}
	l.cmap[index] = value
}

func (l *listObject) del (name Value) {
	if !l.islist {
		delete(l.cmap,name)
	}
}

func (l *listObject) Get (index Value) Value {
	if l.islist {
		return  l.list.find(index.V.(int))
	}

	return l.cmap[index]
}

func NewListValue () Value {
	return Value{
		T: VT_lIST,
		V:Newlist(),
	}
}
