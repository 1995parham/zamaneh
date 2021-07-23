use tokio::time;

enum State {
    Working,
    Fun,
}

#[tokio::main]
async fn main() {
    let mut sec_one_interval = time::interval(time::Duration::from_secs(1));
    let mut working_interval = time::interval(time::Duration::from_secs(60));
    let mut fun_interval = time::interval(time::Duration::from_secs(60));

    let mut state = State::Fun;

    let mut secs = 0;

    loop {
        match state {
            State::Working => {
                tokio::select! {
                    _ = sec_one_interval.tick() => {
                        secs += 1;

                        print!("{}", termion::cursor::Save);
                        println!(
                            "{:0>2}:{:0>2}:{:0>2}",
                            (secs / 60) / 60,
                            (secs / 60) % 60,
                            secs % 60
                        );
                        print!("{}", termion::cursor::Restore);
                    },
                    _ = working_interval.tick() => {
                        secs = 0;
                        state = State::Fun;
                        println!("time to fun");
                        fun_interval = time::interval(time::Duration::from_secs(60));
                        fun_interval.tick().await;
                    },
                }
            }
            State::Fun => {
                tokio::select! {
                    _ = sec_one_interval.tick() => {
                        secs += 1;

                        print!("{}", termion::cursor::Save);
                        println!(
                            "{:0>2}:{:0>2}:{:0>2}",
                            (secs / 60) / 60,
                            (secs / 60) % 60,
                            secs % 60
                        );
                        print!("{}", termion::cursor::Restore);
                    },
                    _ = fun_interval.tick() => {
                        secs = 0;
                        state = State::Working;
                        println!("time to work");
                        working_interval = time::interval(time::Duration::from_secs(60));
                        working_interval.tick().await;
                    },
                }
            }
        };
    }
}
