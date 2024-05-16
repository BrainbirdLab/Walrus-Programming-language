
let a := i8 = 5;
let a := 2;
let b := 10;

let c := a + b;

let d := (a + b) * 2 - -c;



let arr : []i8 = [1, 2, 3, 4, 5]

let arr2 : [][]i8;

struct Array {

    pub length: i8;

    priv _arr: []i8;

    priv _capacity: i8;

    pub static readonly count: i8;

    pub push(elem: i8); // returns nothing
    pub pop() -> i8;
}

/*
impl Array::push(elem: i8) {
    // code here
}

impl Array::pop() -> i8 {
    // code here
    ret -1;
}
*/






/*
impl Array::new(capacity: i8) -> Array {
    return Array {
        _arr: []i8,
        _length: 0,
        _capacity: capacity,
    };
}

impl Array::push(&mut self) {
    if self._length < self._capacity {
        self._arr[self._length] = self._length;
        self._length += 1;
    }
}

impl Array::pop(&mut self) -> i8 {
    if self._length > 0 {
        self._length -= 1;
        return self._arr[self._length];
    }
    return -1;
}
*/