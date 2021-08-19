use std::io::{stdin, stdout};
use std::process::exit;
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

#[tokio::main]
async fn main() {
    let mut args = std::env::args();
    let topic = args.nth(1).unwrap_or_else(|| "-".to_string());

    let mut sec_one_interval = time::interval(time::Duration::from_secs(1));

    let mut secs = 0;

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
                        println!(
                        "{}you are working for {:0>2}:{:0>2}:{:0>2} on {}",
                        cursor::Show,
                        (secs / 60) / 60,
                        (secs / 60) % 60,
                        secs % 60,
                        topic
                        );
                        exit(0);
                    }
                }
            },
            _ = sec_one_interval.tick() => {
                if !is_pause {

                    secs += 1;

                    print!("{}", termion::cursor::Save);
                    println!(
                    "{:0>2}:{:0>2}:{:0>2}",
                    (secs / 60) / 60,
                    (secs / 60) % 60,
                    secs % 60
                    );
                    print!("{}", termion::cursor::Restore);
                }
            }
        }
    }
}
