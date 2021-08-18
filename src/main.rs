use tokio::time;

#[tokio::main]
async fn main() {
    let mut args = std::env::args();
    let topic = args.nth(1).unwrap_or_else(|| "-".to_string());

    let mut sec_one_interval = time::interval(time::Duration::from_secs(1));

    let mut secs = 0;

    println!(
        r#"
Spending time with you is so precious,
I love every minute that we are together.
        "#
    );
    println!("you are working on {}", topic);
    loop {
        sec_one_interval.tick().await;
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
