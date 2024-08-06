const PUZZLE_INPUT: &str = "136760-595730";

fn is_legit_password(password: &u32) -> bool {
    let password_vec = password.to_string().chars().collect::<Vec<char>>();
    let mut adjacent_match: bool = false;
    for idx in 0..password_vec.len()-1 {
        if password_vec[idx].to_digit(10) > password_vec[idx+1].to_digit(10) {
            return false;
        }
        if password_vec[idx].to_digit(10) == password_vec[idx+1].to_digit(10) {
            adjacent_match = true;
        }
    }
    return adjacent_match
}

fn is_legit_password_part2(password: &u32) -> bool {
    let password_vec = password.to_string().chars().map(|x| x.to_digit(10).unwrap()).collect::<Vec<u32>>();
    let mut adjacent_match: bool = false;
    let largest_idx = password_vec.len()-1;
    for idx in 0..largest_idx {
        if password_vec[idx] > password_vec[idx+1] {
            return false;
        }
        if password_vec[idx] == password_vec[idx+1] {
            let statement: bool;
            if idx == 0 {
                statement = password_vec[idx+1] != password_vec[idx+2]
            }
            else if idx == largest_idx-1 {  // no idea why I couldn't do a match statement like this but oh well. Maybe variables dont work in them?
                statement = password_vec[idx-1] != password_vec[idx]
            }
            else {
                statement = (password_vec[idx-1] != password_vec[idx]) & (password_vec[idx+1] != password_vec[idx+2])
            }
            if statement {
                adjacent_match = true;
            }
        }
    }
    return adjacent_match
}

fn main() {
    let parts = PUZZLE_INPUT.split("-").map(|x| x.parse()).collect::<Result<Vec<u32>, _>>().unwrap();
    let lower = parts[0];
    let upper = parts[1];
    
    // num digits = 6
    // value within range lower to upper
    // two adjacent digits are the same
    // going from left to right, monotonically increasing

    // how many different passwords possible?
    let mut count_part_1: u32 = 0;
    let mut count_part_2: u32 = 0;
    for i in lower..upper+1 {
        if is_legit_password(&i) {
            count_part_1 = count_part_1 + 1;
        }
        if is_legit_password_part2(&i) {
            count_part_2 = count_part_2 + 1;
        }
    }

    // part 2 - the two adjacent digits can't be part of a larger group of repeaters

    println!("---AOC 2019 - DAY 4 - p4$$w0rd checker\n\tPart 1: {}\n\tPart 2: {}", count_part_1, count_part_2);
}
