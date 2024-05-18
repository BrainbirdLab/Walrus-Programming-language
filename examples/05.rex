
let x : i8 = 5;
let a : puka = 2;
let b := 10;

a = -10 + 5 * 2;

fn add(a: i8, b: i8) -> i8 {
    ret a + b;
}

fn mul(a: i8, b: i8) -> i8 {
    ret a * b;
}

fn square(a: i8) -> i8 {
    // return a * a;
    ret mul(a, a);
    //ret 4;
}