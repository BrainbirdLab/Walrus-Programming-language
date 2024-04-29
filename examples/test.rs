

// system methods

// os methods

os.time(); // returns the current time in seconds since the epoch

os.sleep(5); // sleeps for 5 seconds

// multi threading syntax proposed

fn func() {
    println("Hello from a thread!");
}


//struct in rust
struct Point {

    pub let x: i8;
    pub const y: i8;

    // method
    pub fn distance(self, other: Point) -> i8 {
        return (self.x - other.x).abs() + (self.y - other.y).abs();
    }

    // no return method
    pub fn print(self) {
        println("x: {}, y: {}", self.x, self.y);
    }
}



fn y() -> i8 {
    // error: no return statement
}


fn a() -> i8 {
    return 1;
}

const thread = threads::new(); // creates a new thread object

/* thread = {
    id: thread_id,
    runnable: fn,

    assign: fn(runnable: fn()) {
        // assign the function to the thread
        runnable = runnable;
        return this;
    }

    sleep_for: fn(ms: i64) {
        // sleep for the specified number of seconds
    }

    run: fn() {
        // run the function
        runnable();
    }
}
*/

thread.assign(func).run(); // runs the function in a new thread


fn main() {
    let a : i8 = 1;  // 8 bit int
    const c : i64 = 4389235677832; // 64 bit int
    
    let x; // error: type must be specified or assigned to a value

    let b := 2; // type inference to smallest possible type (i8)

    // array syntax
    let arr : [i8] = [1, 2, 3, 4, 5]; // array of i8s // static size
    let arr2 : [i8; 5] = [1, 2, 3, 4, 5]; // array of 5 i8s // static size
    let arr3 := [1, 2, 3, 4, 5]; // type inference to array of i8s // static size

    //dynamic
    let arr4 : [i8]; // array of i8  dynamic size
    let arr5 : []; // array of any type // dynamic size
}