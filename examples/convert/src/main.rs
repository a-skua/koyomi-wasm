wit_bindgen::generate!({
    world: "app",
    path: "wit",
    generate_all,
});

use a_skua::koyomi::convert::{Era, Month, WarekiDate, WesternDate};

fn main() {
    let args: Vec<String> = std::env::args().collect();

    if args.len() < 3 {
        print_usage();
        std::process::exit(1);
    }

    match args[1].as_str() {
        "-s" => convert_to_wareki(&args[2]),
        "-w" => convert_to_seireki(&args[2]),
        _ => {
            print_usage();
            std::process::exit(1);
        }
    }
}

fn print_usage() {
    eprintln!("Usage:");
    eprintln!("  convert -s 2019-05-01");
    eprintln!("  convert -w 令和元年5月1日");
    eprintln!("  convert -w 令和1年5月1日");
}

fn convert_to_wareki(input: &str) {
    let parts: Vec<&str> = input.split('-').collect();
    if parts.len() != 3 {
        eprintln!("Error: invalid date format, expected YYYY-MM-DD");
        std::process::exit(1);
    }

    let year: i32 = parts[0].parse().unwrap_or_else(|_| {
        eprintln!("Error: invalid year");
        std::process::exit(1);
    });
    let month: u32 = parts[1].parse().unwrap_or_else(|_| {
        eprintln!("Error: invalid month");
        std::process::exit(1);
    });
    let day: u8 = parts[2].parse().unwrap_or_else(|_| {
        eprintln!("Error: invalid day");
        std::process::exit(1);
    });

    let month = parse_month_number(month);
    let date = WesternDate::new(year, month, day);

    match date.to_wareki() {
        Ok(wareki) => println!("{}", wareki.to_string()),
        Err(e) => {
            eprintln!("Error: {e}");
            std::process::exit(1);
        }
    }
}

fn convert_to_seireki(input: &str) {
    let (era, year, month, day) = parse_wareki(input).unwrap_or_else(|| {
        eprintln!("Error: invalid wareki format, expected e.g. 令和元年5月1日 or 令和1年5月1日");
        std::process::exit(1);
    });

    let date = WarekiDate::new(era, year, month, day);

    match date.to_seireki() {
        Ok(seireki) => println!("{}", seireki.to_string()),
        Err(e) => {
            eprintln!("Error: {e}");
            std::process::exit(1);
        }
    }
}

fn parse_month_number(n: u32) -> Month {
    match n {
        1 => Month::January,
        2 => Month::February,
        3 => Month::March,
        4 => Month::April,
        5 => Month::May,
        6 => Month::June,
        7 => Month::July,
        8 => Month::August,
        9 => Month::September,
        10 => Month::October,
        11 => Month::November,
        12 => Month::December,
        _ => {
            eprintln!("Error: invalid month {n}");
            std::process::exit(1);
        }
    }
}

fn parse_wareki(input: &str) -> Option<(Era, i32, Month, u8)> {
    let era_table: &[(&str, Era)] = &[
        ("明治", Era::Meiji),
        ("大正", Era::Taisho),
        ("昭和", Era::Showa),
        ("平成", Era::Heisei),
        ("令和", Era::Reiwa),
    ];

    let (era, rest) = era_table
        .iter()
        .find_map(|(name, era)| input.strip_prefix(name).map(|rest| (*era, rest)))?;

    // "元年5月1日" or "1年5月1日"
    let (year_str, rest) = rest.split_once('年')?;
    let year: i32 = if year_str == "元" {
        1
    } else {
        year_str.parse().ok()?
    };

    let (month_str, rest) = rest.split_once('月')?;
    let month_num: u32 = month_str.parse().ok()?;

    let day_str = rest.strip_suffix('日')?;
    let day: u8 = day_str.parse().ok()?;

    let month = match month_num {
        1 => Month::January,
        2 => Month::February,
        3 => Month::March,
        4 => Month::April,
        5 => Month::May,
        6 => Month::June,
        7 => Month::July,
        8 => Month::August,
        9 => Month::September,
        10 => Month::October,
        11 => Month::November,
        12 => Month::December,
        _ => return None,
    };

    Some((era, year, month, day))
}
