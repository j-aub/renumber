#[macro_use] extern crate rocket;

use std::io::Write;
use std::process::{Command, Stdio};

use rocket::http::Status;
use rocket::form::{Form, Contextual, FromForm};

use rocket_dyn_templates::{context, Template};

#[derive(Debug, FromForm)]
#[allow(dead_code)]
struct Renumber<'v> {
    list: &'v str,
}

#[get("/")]
fn index() -> Template {
    Template::render("index", context! {list: ""})
}

#[post("/", data = "<form>")]
fn submit<'r>(form: Form<Contextual<'r, Renumber<'r>>>) -> (Status, Template) {
    let child = Command::new("gawk")
        .arg("{if(!/[0-9]+\\./) $0=\"0. \" $0; line++; sub(/^[0-9]+/,line)}1")
        .stdin(Stdio::piped())
        .stdout(Stdio::piped())
        .spawn();

    let list = match form.context.field_value("list") {
        Some(list) => list.as_bytes(),
        None => b"",
    };

    match child {
        Ok(mut child) => {
            let child_stdin = child.stdin.as_mut().unwrap();
            let _  = child_stdin.write_all(&list);


            let output = child.wait_with_output().unwrap().stdout;

            match std::str::from_utf8(&output) {
                Ok(renumbered_list) => (Status{code: 200}, Template::render("index", context! {list: renumbered_list})),
                Err(_) => (Status{code: 500}, Template::render("index", &form.context)),
            }
        }

        Err(_) => (Status{code: 500}, Template::render("index", &form.context)),
    }
}

#[launch]
fn rocket() -> _ {
    rocket::build()
        .mount("/", routes![index, submit])
        .attach(Template::fairing())
}
