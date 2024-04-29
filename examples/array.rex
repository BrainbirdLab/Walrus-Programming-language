
let arr : []i8 = [1, 2, 3, 4, 5];

struct Array {

    length: i8, // public
    _elements: []i8, // private

    LENGTH = 5, // constant

    $_copies: i8, // static, private

    pub length: i8;

    priv static readonly _arr: []i8;
    $_ARR: []i8;
    priv _capacity: i8;

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