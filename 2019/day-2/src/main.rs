use std::fs;
use std::mem;

const FILENAME: &str = "input.txt";

fn parse_input_string(string: &str) -> Vec<u32> {
    return string
                .trim()
                .split(",")
                .map(|x| x.parse())
                .collect::<Result<Vec<u32>, _>>()
                .unwrap();
}

fn execute_intcode(mut vec: Vec<u32>, main: &bool, noun: &u32, verb: &u32) -> u32 {
    // initial set up
    if *main {
        let _  = mem::replace(&mut vec[1], *noun);
        let _ = mem::replace(&mut vec[2], *verb);
    }

    let mut addition: bool = true;

    for i in (0..(vec.len() - 1)).step_by(4) {
        // starts with 0 .. we should use this to access what's in 0, 1, 2, 3 i.e. i, i+1, i+2, i+3

        let pos_0: usize = vec[i] as usize;
        let pos_1: usize = vec[i+1] as usize;
        let pos_2: usize = vec[i+2] as usize;
        let pos_3: usize = vec[i+3] as usize;

        match pos_0 {
            1 => addition = true,
            2 => addition = false,  // so doing product
            99 => return vec[0],
            _ => println!("Oh no oh no oh no!!!")
        }

        let replacement_value: u32;
        if addition {
            replacement_value = vec[pos_1] + vec[pos_2];
        }
        else {
            replacement_value = vec[pos_1] * vec[pos_2];
        }

        let _ = mem::replace(&mut vec[pos_3], replacement_value);
    }
    
    return vec[0];
}

fn part2(vec: Vec<u32>) -> u32 {
    for noun in 0..99 {
        for verb in 0..99 {
            let tmp_vec = vec.clone();
            if execute_intcode(tmp_vec, &true, &(noun as u32), &(verb as u32)) == 19690720 {
                return (100 * noun) + verb;
            }
        }
    }
    panic!("Shouldn't be here! Didn't exit from our search.");
}

fn main() {
    let input_string = fs::read_to_string(FILENAME).expect("Failed to read input");

    // Test case
    let test_case = parse_input_string("1,9,10,3,2,3,11,0,99,30,40,50");
    let test_case_ans = execute_intcode(test_case, &false, &12, &2);
    assert_eq!(test_case_ans, 3500);

    // Real input
    let input_vec: Vec<u32> = parse_input_string(&input_string);
    let part1_ans: u32 = execute_intcode(input_vec, &true, &12, &2);

    let part2_ans: u32 = part2(
        parse_input_string(&input_string)
    );

    println!("---AOC 2019 - Day 2 - Intcode!---\n\tPart 1: {}\n\tPart2: {}", part1_ans, part2_ans);
}
