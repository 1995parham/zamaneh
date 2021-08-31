use clap::{App, Arg};
use std::io::{stdin, stdout};
use std::process::exit;
use std::time::Duration;
use termion::cursor;
use termion::event::Key;
use termion::input::TermRead;
use termion::raw::IntoRawMode;
use tokio::sync::mpsc::channel;
use tokio::{select, time};

enum Event {
    Pause,
    Quit,
}

const TOPIC_ARG_NAME: &str = "title";

#[tokio::main]
async fn main() {
    let args = App::new("zamaneh")
        .about("Manage your working periods with ease")
        .arg(
            Arg::with_name(TOPIC_ARG_NAME)
                .takes_value(true)
                .help("working period title")
                .default_value("-"),
        )
        .get_matches();
    let topic = args.value_of(TOPIC_ARG_NAME).unwrap();

    let mut sec_one_interval = time::interval(time::Duration::from_secs(1));

    let mut secs = Duration::new(0, 0);

    let (tx, mut rx) = channel::<Event>(1);

    let stdin = stdin();
    let stdout = stdout().into_raw_mode().unwrap();
    stdout.suspend_raw_mode().unwrap();

    let mut is_pause = false;

    tokio::spawn(async move {
        for c in stdin.keys() {
            match c.unwrap() {
                Key::Char('q') => {
                    let _ = tx.send(Event::Quit).await;
                }
                Key::Char(' ') => {
                    let _ = tx.send(Event::Pause).await;
                }
                _ => {}
            }
        }
    });

    println!(
        r#"{}
Spending time with you is so precious,
I love every minute that we are together.
"#,
        cursor::Hide
    );
    println!("you are working on {}", topic);
    stdout.activate_raw_mode().unwrap();

    loop {
        select! {
            Some(ev) = rx.recv() => {
                match ev {
                    Event::Pause => {
                        is_pause = !is_pause;
                    },
                    Event::Quit => {
                        stdout.suspend_raw_mode().unwrap();
                        println!("{}you are working for {} on on {}", cursor::Show, format_duration(secs), topic);
                        exit(0);
                    }
                }
            },
            _ = sec_one_interval.tick() => {
                print!("{}", termion::cursor::Save);

                if !is_pause {
                    secs += Duration::from_secs(1);
                    print!("{}", termion::color::Fg(termion::color::LightRed));
                } else {
                    print!("{}", termion::color::Fg(termion::color::LightBlue));
                }

                println!("{}", format_duration(secs));
                print!("{}", termion::color::Fg(termion::color::Reset));
                print!("{}", termion::cursor::Restore);
            }
        }
    }
}

fn format_duration(d: Duration) -> String {
    let secs = d.as_secs();
    return format!(
        "{:0>2}:{:0>2}:{:0>2}",
        (secs / 60) / 60,
        (secs / 60) % 60,
        secs % 60
    );
}
