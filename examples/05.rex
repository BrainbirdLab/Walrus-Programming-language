mod main;

import "core:fmt";
import { readFile, writeFile } from "core:io";

let x : i8 = 5;
let a := 2;
let b := 10;

let myName := "John";

a = -10 + 5 - b;

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