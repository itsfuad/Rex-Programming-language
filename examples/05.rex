mod main;

import "core:fmt";
import { readFile, writeFile } from "core:io";

let x : i8 = 5;
let a := 2;
let b := 10;

let myName := "John";

a = -10 + 5 - b;

if a == 4 {
    let a := 10;
    b = a + 3;
}

struct Person {
    pub name: str;
    pub age: i8;
}

let p := Person { name: "John", age: 23 };

fn add(a: i8, b: i8) -> i8 {
    ret a + b;
}

/*
for i := 0; i < 10; ++i {
    // print(i);
}
let array := [1, 2, 3, 4, 5];

foreach i, arr in array {
    // i is the index
    // arr is the value on each iteration
}

foreach val in array where val % 2 == 0{
    // val is the values that are even
}
*/



fn mul(a: i8, b: i8) -> i8 {
    ret a * b;
}

fn square(a: i8) -> i8 {
    // return a * a;
    ret mul(a, a);
    //ret 4;
}