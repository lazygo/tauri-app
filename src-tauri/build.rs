use failure::{format_err, Error};
use std::env;
use std::process::Command;

/// Just useful trait to run a command
trait RunIt {
    fn run_it(&mut self, err: &str) -> Result<(), Error>;
}

impl RunIt for Command {
    fn run_it(&mut self, err: &str) -> Result<(), Error> {
        let output = self.output()?;
        if !output.status.success() {
            let out = String::from_utf8_lossy(&output.stderr);
            eprintln!("{}", out);
            Err(format_err!("{}", err))
        } else {
            Ok(())
        }
    }
}

fn main() -> Result<(), Error> {
    let out_dir = env::var("OUT_DIR").unwrap();
    let lib = "service";

    Command::new("sh")
        .arg("../src-golib/build.sh")
        .args(&[out_dir.clone(), format!("{}", lib)])
        .status()
        .expect("service fail");

    println!("cargo:rustc-link-search=native={}", out_dir);
    println!("cargo:rustc-link-lib=static={}", lib);
    tauri_build::build();

    // Activate this feature to rebuild this dependency everytime
    if cfg!(feature = "refresh") {
        Command::new("touch")
            .args(&["build.rs"])
            .run_it("Can't touch the build file")?;
    }

    Ok(())
}
