use std::iter::zip;
use std::fs;

struct BoatInfo {
    times: Vec<u64>,
    speeds: Vec<u64>,
}

fn parse_line(raw_input_line: &str, part2: &bool) -> Vec<u64> {
    if *part2 {
        let x =  raw_input_line
            .trim()
            .split_ascii_whitespace()
            .skip(1)  // similar to below this skips the "X:" bit at the start
            .collect::<Vec<&str>>()  // this then collects it all into a vector
            .join("")  // this joins together like python "".join(list[str])
            .parse();  // may have been able to do unwrap here to stop it eing a Result type thing and make it the int, then wouldn't need an expect?

        return vec![x.expect("Should be a u64?")];    
    }
    else {
        return raw_input_line 
            .trim()
            .split_ascii_whitespace()
            .skip(1)  // "Time:" or "Distance:"
            .map(|x| x.parse())  // I think parse() parses the string as number because its number characters?
            .collect::<Result<Vec<u64>, _>>()
            .expect("Couldn't parse?")
    };
}

fn parse_input(example: &bool, part2: &bool) -> BoatInfo {
    let filename;
    if *example {
        filename = "example.txt";
    }
    else {
        filename = "input.txt";
    }

    let contents = fs::read_to_string(&filename).expect("Should have been able to read input.");


    let (times, distances) = contents
                                .trim()
                                .split_once("\n")
                                .map(|(t, d)| (parse_line(t, part2), parse_line(d, part2)))
                                .expect("not enough lines??");
    
    let boat_info = BoatInfo {
        times: times,
        speeds: distances,
    };

    return boat_info;
}

fn calculate_if_win(button_press: &u64, total_time: &u64, distance_record: &u64) -> bool {
    let speed = *button_press;
    let time_remaining = *total_time - button_press;
    let distance_travelled = speed * time_remaining;
    return distance_travelled > *distance_record
}

fn calculate_number_of_ways_to_win(total_time: &u64, distance_record: &u64) -> u64 {
    let mut count: u64 = 0;
    for i in 0..*total_time+1 {
        if calculate_if_win(&i, total_time, distance_record) {        
            count = count + 1;
        }
    }

    return count
}

fn day6(time_arr: &Vec<u64>, distance_arr: &Vec<u64>) -> u64 {
    let time_dist_iter = zip(time_arr, distance_arr);
    let res: u64 = time_dist_iter.into_iter().map(|(t, d)| calculate_number_of_ways_to_win(&t, &d)).product();

    return res
}

fn main() {
    let boat_info_example_p1: BoatInfo = parse_input(&true, &false);
    let boat_info_real_p1: BoatInfo = parse_input(&false, &false);
    let boat_info_example_p2: BoatInfo = parse_input(&true, &true);
    let boat_info_real_p2: BoatInfo = parse_input(&false, &true);

    println!("--------------------\n     AOC DAY 6\n--------------------\n");

    let example_ans_p1 = day6(&boat_info_example_p1.times, &boat_info_example_p1.speeds);
    let real_ans_p1 = day6(&boat_info_real_p1.times, &boat_info_real_p1.speeds);

    assert_eq!(example_ans_p1, 288);
    assert_eq!(real_ans_p1, 74698);

    let example_ans_p2 = day6(&boat_info_example_p2.times, &boat_info_example_p2.speeds);
    let real_ans_p2 = day6(&boat_info_real_p2.times, &boat_info_real_p2.speeds);

    assert_eq!(example_ans_p2, 71503);
    assert_eq!(real_ans_p2, 27563421);

    println!("PART 1");
    println!("\tExample: {}", example_ans_p1);
    println!("\tReal: {}", real_ans_p1);

    println!("PART 2");
    println!("\tExample: {}", example_ans_p2);
    println!("\tReal: {}", real_ans_p2);

    println!("\n\n:)");
}
