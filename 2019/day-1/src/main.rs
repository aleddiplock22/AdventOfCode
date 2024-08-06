use std::fs;

const FILENAME: &str = "input.txt";

fn parse_input(filename: &str) -> Vec<i128> {
    let contents = fs::read_to_string(filename).expect("Input should be read as string.");

    let parsed_input = contents
                        .lines()
                        .map(|x| x.trim().parse())
                        .collect::<Result<Vec<i128>, _>>()
                        .expect("Couldn't parse?");
    
    return parsed_input;
}

fn part1(input: &Vec<i128>) -> i128 {
    let answer: i128 = input
                        .iter()
                        .map(|&x| (x / 3) - 2)  // Doing division in integer automatically does what we want (divise then floor func)
                        .collect::<Vec<i128>>()
                        .iter()
                        .sum();

    return answer;
}

fn accumalate_fuel(mass: &i128, total: &i128) -> i128 {
    let mass = (mass / 3) - 2;
    if mass <= 0 {
        return *total;
    }
    let total = total + mass;
    
    return accumalate_fuel(&mass, &total);
}

fn part2(input: &Vec<i128>) -> i128 {
    return input
            .iter()
            .map(|&x| accumalate_fuel(&x, &0))
            .collect::<Vec<i128>>()
            .iter()
            .sum();
}

fn main() {
    let input = parse_input(FILENAME);

    let part1_result = part1(&input);
    let part2_result = part2(&input);

    println!("---AOC Day 1---\n\tPart1: {}\n\tPart2: {}", part1_result, part2_result);
}
