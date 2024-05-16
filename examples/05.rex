
let a := i8 = 5;
let a := 2;
let b := 10;

fn add(a: i8, b: i8) -> i8 {
    return a + b;
}

fn mul(a: i8, b: i8) -> i8 {
    return a * b;
}

fn square(a: i8) -> i8 {
    // return a * a;
    return mul(a, a);
}